package builder

import (
	"async_executor/amqp"
	"async_executor/executor"
	"async_executor/output"
	"time"
)

type ExecutionBuildSettings struct {
	ConnectionString                     string
	ConnectionAttemptsCount              int
	ProcessorsCount                      int
	FailedMessagesRepeatAttemptsCount    int
	FailedMessagesRepeatDelay            time.Duration
	OutputBadlyFormattedMessagesToErrors bool
	OutputResultToExistingExchange       string
}

func BuildExecutorFacade(
	funcExecutors []executor.AsyncFunctionExecutor,
	buildSettings ExecutionBuildSettings,
) (facade *executor.ExecutionFacade, err error) {
	conn, err := BuildAmqpConnection(buildSettings.ConnectionString, buildSettings.ConnectionAttemptsCount)
	if err != nil {
		return
	}

	amqpPayloadsProvider, err := amqp.NewAmqpAsyncCallsProvider(
		AsyncInputsQueueName,
		&amqp.InputReadingSettings{},
		conn.Conn,
	)
	if err != nil {
		return
	}

	executionSettings := executor.ExecutionSettings{
		buildSettings.ProcessorsCount,
		buildSettings.OutputBadlyFormattedMessagesToErrors,
		buildSettings.FailedMessagesRepeatAttemptsCount,
	}

	errOutputter, err := amqp.NewAmqpOutputter(AsyncErrorsQueueName, amqp.AmqpWritingSettings{}, conn.Conn)
	if err != nil {
		return
	}

	var delayedOutputter output.Outputter = output.NullOutputter{}
	if buildSettings.FailedMessagesRepeatAttemptsCount > 0 && buildSettings.FailedMessagesRepeatDelay > 0 {
		delayedOutputter, err = amqp.NewDelayedAmqpOutputter(
			AsyncInputsQueueName,
			amqp.AmqpWritingSettings{LifeTime: buildSettings.FailedMessagesRepeatDelay},
			conn.Conn,
		)
		if err != nil {
			return
		}
	}

	var resultOutputter output.Outputter = output.NullOutputter{}
	if buildSettings.OutputResultToExistingExchange != "" {
		resultOutputter, err = amqp.NewExistingAmqpOutputter(
			conn.Conn,
			buildSettings.OutputResultToExistingExchange,
			amqp.AmqpWritingSettings{},
		)
		if err != nil {
			return
		}
	}

	executorsRegistry := BuildExecutorsRegistry(funcExecutors)

	termChan := make(chan bool)
	go func() {
		for _ = range conn.ConnClosedNotificationChannel {
			termChan <- true
		}
	}()

	asyncExecFacade := executor.NewExecutionFacade(
		termChan,
		executorsRegistry,
		amqpPayloadsProvider,
		executionSettings,
		errOutputter,
		resultOutputter,
		delayedOutputter,
	)

	return asyncExecFacade, nil
}
