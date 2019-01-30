package builder

import (
	"async_executor/amqp"
	"async_executor/executor"
)

func BuildAmqpCallerFacade(
	funcExecutors []executor.AsyncFunctionExecutor,
	connectionString string,
	connectionAttemptsCount int,
) (facade *executor.AsyncFuncRegistrationFacade, err error) {
	executorsRegistry := BuildExecutorsRegistry(funcExecutors)

	conn, err := BuildAmqpConnection(connectionString, connectionAttemptsCount)
	if err != nil {
		return
	}

	writingSettings := amqp.AmqpWritingSettings{}

	amqpFuncRegistrator, err := amqp.NewAmqpAsyncFuncRegistrator(AsyncInputsQueueName, writingSettings, conn.Conn)
	if err != nil {
		return
	}

	facade = executor.NewAsyncFuncRegistrationFacade(executorsRegistry, amqpFuncRegistrator)

	return
}
