package executor

import (
	"async_executor/dto"
	"async_executor/logger"
	"async_executor/output"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ExecutionFacade struct {
	termChan             chan bool
	execPayloadsProvider ExecutionPayloadsProvider
	settings             ExecutionSettings
	errorWriter          output.Outputter
	outputWriter         output.Outputter
	delayWriter          output.Outputter
	funcExecutorRegistry *AsyncFuncExecutorRegistry
}

func NewExecutionFacade(
	termChan chan bool,
	funcExecutorRegistry *AsyncFuncExecutorRegistry,
	executionPayloadsProvider ExecutionPayloadsProvider,
	settings ExecutionSettings,
	errorWriter output.Outputter,
	outputWriter output.Outputter,
	delayWriter output.Outputter,
) *ExecutionFacade {
	return &ExecutionFacade{
		termChan:             termChan,
		execPayloadsProvider: executionPayloadsProvider,
		settings:             settings,
		errorWriter:          errorWriter,
		outputWriter:         outputWriter,
		delayWriter:          delayWriter,
		funcExecutorRegistry: funcExecutorRegistry,
	}
}

func (ef *ExecutionFacade) ExecuteAsyncFunctions() error {
	termSignalChannel := make(chan os.Signal, 2)
	signal.Notify(termSignalChannel, os.Interrupt, syscall.SIGTERM)

	asyncInputMessageChannel, err := ef.execPayloadsProvider.GetAsyncCalls()
	if err != nil {
		return err
	}

	logger.Log("Will start processing async inputs with %d processors\n", ef.settings.ExecutionProcessorsCount)

	consumerTermChannels := []chan bool{}
	for y := 0; y < ef.settings.ExecutionProcessorsCount; y++ {
		boolTermChan := make(chan bool)
		consumerTermChannels = append(consumerTermChannels, boolTermChan)
		go ef.runInBackground(asyncInputMessageChannel, boolTermChan)
	}

	select {
	case <-ef.termChan:
		ef.sendTermToAllConsumers(consumerTermChannels)
	case <-termSignalChannel:
		ef.sendTermToAllConsumers(consumerTermChannels)
		return nil
	}

	return nil
}

func (ef *ExecutionFacade) runInBackground(asyncInputMessageChannel <-chan dto.RawInput, termChan chan bool) {
	for {
		select {
		case <-termChan:
			return
		case msg := <-asyncInputMessageChannel:
			if len(msg.GetPayload()) == 0 {
				continue
			}
			body := string(msg.GetPayload())
			logger.Log("Received message: %s", body)

			funcInput, err := dto.NewAsyncFuncInput(msg.GetPayload())
			if err != nil {
				logger.LogError(err, "Incompatible message format")

				if ef.settings.PublishBadlyFormattedMessagesToErrorChannel {
					err = ef.errorWriter.OutputMessage(dto.StringOutputMessage(body))
					if err != nil {
						logger.LogError(err, "")
					}
				}

				ef.markAsDone(msg, "")

				continue
			}
			ef.execMsg(msg, funcInput)
		}
	}
}

func (ef *ExecutionFacade) execMsg(origMsg dto.RawInput, funcInput dto.AsyncFuncInput) {
	funcInput.CallsCount += 1

	defer func() {
		ef.markAsDone(origMsg, funcInput.MessageId)
	}()

	err := ef.validateInputWindow(funcInput)
	if err != nil {
		logger.LogError(err, "")
		return
	}

	processor, ok := ef.funcExecutorRegistry.GetFuncExecutor(funcInput.FunctionName)

	if !ok {
		logger.LogErrorFromMsg(
			"Unknown anync function name %s for message %s",
			funcInput.FunctionName,
			funcInput.MessageId,
		)

		funcInput.FailedAttemptsCount += 1
		funcInput.LastError = errors.New("Unknown rpc function name")

		err = ef.errorWriter.OutputMessage(funcInput)
		if err != nil {
			logger.LogErrorFromMsg(err.Error(), "Publishing failure for message %s", funcInput.MessageId)
		}

		return
	}

	outputMsg, processorResult := processor.Process(funcInput.Payload)
	switch processorResult.(type) {
	case output.NonOutputtableSuccess:
		return
	case output.OutputtableSuccess:
		ef.publishResult(outputMsg, ef.outputWriter, funcInput.MessageId)
	case output.NonRecoverableError:
		funcInput.FailedAttemptsCount += 1
		funcInput.LastError = processorResult.GetError()
		ef.publishResult(funcInput, ef.errorWriter, funcInput.MessageId)
	case output.RecoverableError:
		funcInput.FailedAttemptsCount += 1
		funcInput.LastError = processorResult.GetError()
		if funcInput.FailedAttemptsCount > ef.settings.FailedMessagesRepeatAttemptsCount {
			ef.publishResult(funcInput, ef.errorWriter, funcInput.MessageId)
		} else {
			ef.publishResult(funcInput, ef.delayWriter, funcInput.MessageId)
		}
	}
}

func (ef *ExecutionFacade) publishResult(outputMsg dto.OutputMessage, outputter output.Outputter, msgId string) {
	err := outputter.OutputMessage(outputMsg)
	if err != nil {
		logger.LogErrorFromMsg(err.Error(), "Publishing failure for message %s", msgId)
	}
}

func (ef *ExecutionFacade) markAsDone(origMsg dto.RawInput, msgId string) {
	err := ef.execPayloadsProvider.MarkAsDone(origMsg)
	if err != nil {
		logger.LogError(err, fmt.Sprintf("Marking input as done error for message %s", msgId))
	}
}

func (ef *ExecutionFacade) sendTermToAllConsumers(consumerTermChannels []chan bool) {
	for _, consumerTermChannel := range consumerTermChannels {
		consumerTermChannel <- true
		close(consumerTermChannel)
	}
}

func (ef *ExecutionFacade) validateInputWindow(msg dto.AsyncFuncInput) error {
	if msg.ValidWindow <= 0 {
		return nil
	}

	msgValidTillTime := time.Time(msg.TimeStamp).Add(time.Second * time.Duration(msg.ValidWindow))
	if time.Now().UTC().After(msgValidTillTime) {
		return fmt.Errorf("Message valid time %v is elapsed", msgValidTillTime)
	}

	return nil
}
