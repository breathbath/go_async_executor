package dto

//RawInput payloads we get directly from queues
type RawInput interface {
	GetPayload() []byte
	GetId() string
}
