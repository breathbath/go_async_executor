package executor

import (
	"github.com/breathbath/go_async_executor/dto"
)

//ExecutionPayloadsProvider implementations for dequeueing raw queue messages
type ExecutionPayloadsProvider interface {
	MarkAsDone(msg dto.RawInput) error
	GetAsyncCalls() (<-chan dto.RawInput, error)
}
