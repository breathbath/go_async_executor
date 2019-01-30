package output

import "async_executor/dto"

//NullOutputter implements Outputter to discard outputs of async func calls
type NullOutputter struct {

}

func (no NullOutputter) OutputMessage(msg dto.OutputMessage) error {
	return nil
}
