package executor

import (
	"go_async_executor/dto"
	"go_async_executor/output"
)

//AsyncFunctionExecutor implementations for concrete async functions logic
type AsyncFunctionExecutor interface {
	Process(input string) (dto.OutputMessage, output.ExecutionOutput)
	GetFuncName() string
}
