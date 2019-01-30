package output

import "async_executor/dto"

type Outputter interface {
	OutputMessage(msg dto.OutputMessage) error
}
