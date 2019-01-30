package executor

import (
	"async_executor/dto"
)

type AsyncFuncRegistrator interface {
	RegisterAsyncExecution(msg dto.AsyncFuncInput) error
}