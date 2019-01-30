package output

import "github.com/breathbath/go_async_executor/dto"

type Outputter interface {
	OutputMessage(msg dto.OutputMessage) error
}
