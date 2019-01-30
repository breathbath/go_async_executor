package output

import "async_executor/dto"

type NullOutputter struct {

}

func (no NullOutputter) OutputMessage(msg dto.OutputMessage) error {
	return nil
}
