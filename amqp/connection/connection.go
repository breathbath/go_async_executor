package connection

import "github.com/streadway/amqp"

type AmqpConnection struct {
	Conn *amqp.Connection
	ConnClosedNotificationChannel chan *amqp.Error
}
