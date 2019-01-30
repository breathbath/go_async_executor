package builder

import (
	"async_executor/examples/func_caller"
	"async_executor/executor"
	"fmt"
	"time"
)

func GetFunctions() []executor.AsyncFunctionExecutor {
	timeOutputFunc := func_caller.NewStringFunctionExecutor("time_executor", func(input string) error {
		fmt.Println(time.Now(), input)
		return nil
	})

	return []executor.AsyncFunctionExecutor{
		timeOutputFunc,
	}
}
