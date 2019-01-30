package builder

import (
	"async_executor/examples/func_caller"
	"async_executor/executor"
	"async_executor/funcAdapters"
	"errors"
	"fmt"
	"time"
)

func GetFunctions() []executor.AsyncFunctionExecutor {
	timeOutputFunc := func_caller.NewStringFunctionExecutor("time_executor", func(input string) error {
		fmt.Println(time.Now(), input)
		return nil
	})

	failingFunc := funcAdapters.NewRecoverableNonReturningFunc(
		"fail_me",
		func(input string) error {
			return errors.New("Unknown failure")
		},
	)

	return []executor.AsyncFunctionExecutor{
		timeOutputFunc,
		failingFunc,
	}
}
