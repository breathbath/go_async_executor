package executor

import (
	"go_async_executor/dto"
)

//AsyncFuncRegistrator allows polymorfic implementations for enqueuing async function payloads
type AsyncFuncRegistrator interface {
	RegisterAsyncExecution(msg dto.AsyncFuncInput) error
}