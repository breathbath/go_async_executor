package executor

import (
	"github.com/breathbath/go_async_executor/dto"
	"github.com/breathbath/go_async_executor/output"
)

//AsyncFunctionExecutor implementations for concrete async functions logic
type AsyncFunctionExecutor interface {
	Process(input string) (dto.OutputMessage, output.ExecutionOutput)
	GetFuncName() string
}
