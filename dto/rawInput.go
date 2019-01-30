package dto

type RawInput interface {
	GetPayload() []byte
	GetId() string
}
