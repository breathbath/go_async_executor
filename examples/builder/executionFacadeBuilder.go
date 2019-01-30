package builder

import (
	"go_async_executor/amqp/builder"
	"go_async_executor/executor"
	"time"
)

func BuildAsyncFuncExecutionFacade() (facade *executor.ExecutionFacade, err error) {
	settings := builder.ExecutionBuildSettings{
		AmqpConnectionString:                 "amqp://guest:guest@localhost:5672/",
		ConnectionAttemptsCount:              10,
		ProcessorsCount:                      5,
		FailedMessagesRepeatAttemptsCount:    2,
		FailedMessagesRepeatDelay:            time.Second * 30,
		OutputBadlyFormattedMessagesToErrors: false,
		OutputResultToExistingExchange:       "",
		ErrorMessagesLifeTime:                time.Minute * 10,
	}

	return builder.BuildExecutorFacade(GetFunctions(), settings)
}
