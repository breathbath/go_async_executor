package executor

import (
	"async_executor/dto"
	"async_executor/output"
)

//AsyncFunctionExecutor implementations for concrete async functions logic
type AsyncFunctionExecutor interface {
	Process(input string) (dto.OutputMessage, output.ExecutionOutput)
	GetFuncName() string
}
