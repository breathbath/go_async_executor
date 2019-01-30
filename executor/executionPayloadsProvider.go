package executor

import (
	"async_executor/dto"
)

//ExecutionPayloadsProvider implementations for dequeueing raw queue messages
type ExecutionPayloadsProvider interface {
	MarkAsDone(msg dto.RawInput) error
	GetAsyncCalls() (<-chan dto.RawInput, error)
}
