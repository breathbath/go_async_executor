package amqp

import "github.com/streadway/amqp"

//InputReadingSettings customisation for amqp consumer logic
type InputReadingSettings struct {
	ConsumerName string
	AutoAck      bool
	Exclusive    bool
	NoLocal      bool
	NoWait       bool
	Args         amqp.Table
}
