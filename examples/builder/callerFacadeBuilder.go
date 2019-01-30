package builder

import (
	"async_executor/amqp/builder"
	"async_executor/executor"
)

func BuildAsyncFuncRegistrationFacade() (facade *executor.AsyncFuncRegistrationFacade, err error) {
	return builder.BuildAmqpCallerFacade(
		GetFunctions(),
		"amqp://guest:guest@localhost:5672/",
		10,
	)
}
