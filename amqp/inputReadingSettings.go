package amqp

import "github.com/streadway/amqp"

type InputReadingSettings struct {
	ConsumerName string
	AutoAck      bool
	Exclusive    bool
	NoLocal      bool
	NoWait       bool
	Args         amqp.Table
}
