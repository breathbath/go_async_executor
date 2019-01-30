package dto

type OutputMessage interface {
	Serialize() (string, error)
}

type StringOutputMessage string

func (uom StringOutputMessage) Serialize() (string, error) {
	return string(uom), nil
}
