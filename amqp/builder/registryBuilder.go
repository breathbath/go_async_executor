package builder

import (
	"async_executor/executor"
)

func BuildExecutorsRegistry(funcExecutors []executor.AsyncFunctionExecutor) *executor.AsyncFuncExecutorRegistry {
	funcExecutorsRegistry := executor.NewAsyncFuncExecutorRegistry()
	for _, funcExecutor := range funcExecutors {
		funcExecutorsRegistry.AddAsyncFunction(funcExecutor)
	}

	return funcExecutorsRegistry
}
