package builder

import (
	"async_executor/amqp"
	"async_executor/executor"
	"async_executor/output"
	"time"
)

type ExecutionBuildSettings struct {
	AmqpConnectionString                 string
	ConnectionAttemptsCount              int           //will try to connect to RabbitMq this amount of times before failing
	ProcessorsCount                      int           //how many routines to start for processing async func payloads
	FailedMessagesRepeatAttemptsCount    int           //after this amount of times func payload will be discarded
	FailedMessagesRepeatDelay            time.Duration //on failure a func execution will be repeated after delay time
	OutputBadlyFormattedMessagesToErrors bool          //shall badly formatted payloads be published to the errors queue
	OutputResultToExistingExchange       string        //returned result will be published to this exchange for further
	ErrorMessagesLifeTime                time.Duration //for how long queue should keep error messages
}

func BuildExecutorFacade(
	funcExecutors []executor.AsyncFunctionExecutor,
	buildSettings ExecutionBuildSettings,
) (facade *executor.ExecutionFacade, err error) {
	conn, err := BuildAmqpConnection(buildSettings.AmqpConnectionString, buildSettings.ConnectionAttemptsCount)
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

	errOutputter, err := amqp.NewAmqpOutputter(
		AsyncErrorsQueueName,
		amqp.AmqpWritingSettings{LifeTime: buildSettings.ErrorMessagesLifeTime},
		conn.Conn,
	)
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
