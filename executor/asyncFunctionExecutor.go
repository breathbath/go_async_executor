package executor

import (
	"async_executor/dto"
	"async_executor/output"
)

type AsyncFunctionExecutor interface {
	Process(input string) (dto.OutputMessage, output.ExecutionOutput)
	GetFuncName() string
}
