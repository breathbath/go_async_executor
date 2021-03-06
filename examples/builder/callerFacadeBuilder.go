package builder

import (
	"github.com/breathbath/go_async_executor/amqp/builder"
	"github.com/breathbath/go_async_executor/executor"
)

func BuildAsyncFuncRegistrationFacade() (facade *executor.AsyncFuncRegistrationFacade, err error) {
	return builder.BuildAmqpCallerFacade(
		GetFunctions(),
		"amqp://guest:guest@localhost:5672/",
		10,
	)
}
