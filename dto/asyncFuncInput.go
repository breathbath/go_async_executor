package dto

import (
	"encoding/json"
	"time"
)

type AsyncFuncInput struct {
	FunctionName        string    `json:"function"`
	Payload             string    `json:"payload"`
	CallsCount          int64     `json:"calls_count"`
	TimeStamp           time.Time `json:"timestamp"`
	ValidWindow         int64     `json:"valid_window"`
	MessageId           string    `json:"id"`
	FailedAttemptsCount int     `json:"failed_attempts_count"`
	LastError           error     `json:"last_error"`
}

func NewAsyncFuncInput(rawMsg []byte) (AsyncFuncInput, error) {
	var msg AsyncFuncInput
	err := json.Unmarshal(rawMsg, &msg)
	return msg, err
}

func (afi AsyncFuncInput) Serialize() (string, error) {
	jsonRaw, err := json.Marshal(afi)
	return string(jsonRaw), err
}
