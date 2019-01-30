package dto

//OutputMessage async func returned result
type OutputMessage interface {
	Serialize() (string, error)
}

//StringOutputMessage adapts simple string variables to the OutputMessage interface
type StringOutputMessage string

func (uom StringOutputMessage) Serialize() (string, error) {
	return string(uom), nil
}
