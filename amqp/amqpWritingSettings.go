package amqp

import (
	"github.com/streadway/amqp"
	"time"
)

//AmqpWritingSettings customisation for amqp publishing logic
type AmqpWritingSettings struct {
	LifeTime             time.Duration
	Args                 *amqp.Publishing
	Mandatory, Immediate bool
}
