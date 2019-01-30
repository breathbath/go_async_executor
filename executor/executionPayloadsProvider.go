package executor

import (
	"async_executor/dto"
)

type ExecutionPayloadsProvider interface {
	MarkAsDone(msg dto.RawInput) error
	GetAsyncCalls() (<-chan dto.RawInput, error)
}
