package output

import "go_async_executor/dto"

type Outputter interface {
	OutputMessage(msg dto.OutputMessage) error
}
