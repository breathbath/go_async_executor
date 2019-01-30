package amqp

import (
	"github.com/streadway/amqp"
	"time"
)

type AmqpWritingSettings struct {
	RoutingKey           string
	LifeTime             time.Duration
	Args                 *amqp.Publishing
	Mandatory, Immediate bool
}
