package executor

import (
	"go_async_executor/dto"
	"go_async_executor/logger"
	"go_async_executor/output"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//ExecutionFacade main entry point to execute async functions payloads
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
			"Unknown anync function name: %s for message %s",
			funcInput.FunctionName,
			funcInput.MessageId,
		)

		funcInput.FailedAttemptsCount += 1
		funcInput.LastError = "Unknown async function name: " + funcInput.FunctionName

		err = ef.errorWriter.OutputMessage(funcInput)
		if err != nil {
			logger.LogErrorFromMsg(err.Error(), "Publishing failure for message %s", funcInput.MessageId)
		}

		return
	}

	outputMsg, processorResult := processor.Process(funcInput.Payload)
	switch processorResult.(type) {
	case output.NonOutputtableSuccess:
		logger.Log("Execution success for message %s", funcInput.MessageId)
		return
	case output.OutputtableSuccess:
		logger.Log("Execution success of message %s, will send the result further", funcInput.MessageId)
		ef.publishResult(outputMsg, ef.outputWriter, funcInput.MessageId)
	case output.NonRecoverableError:
		logger.LogError(processorResult.GetError(), "Execution of message %s failed, will not process message further", funcInput.MessageId)
		funcInput.FailedAttemptsCount += 1
		funcInput.LastError = processorResult.GetError().Error()
		ef.publishResult(funcInput, ef.errorWriter, funcInput.MessageId)
	case output.RecoverableError:
		funcInput.FailedAttemptsCount += 1
		funcInput.LastError = processorResult.GetError().Error()
		if funcInput.FailedAttemptsCount > ef.settings.FailedMessagesRepeatAttemptsCount {
			logger.LogError(
				processorResult.GetError(),
				"Execution of message %s failed and errors limit %d is elapsed, will not process message further",
				ef.settings.FailedMessagesRepeatAttemptsCount,
				funcInput.MessageId,
			)
			ef.publishResult(funcInput, ef.errorWriter, funcInput.MessageId)
		} else {
			logger.LogError(processorResult.GetError(), "Execution of message %s failed, will try to repeat execution later", funcInput.MessageId)
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
		logger.LogError(err, "Marking input as done error for message %s", msgId)
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
		return fmt.Errorf("Message's %s lifetime %v is elapsed, will ignore it", msg.MessageId, msgValidTillTime)
	}

	return nil
}
