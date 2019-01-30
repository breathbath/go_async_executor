package connection

import (
	"github.com/breathbath/go_utils/utils/connections"
	"github.com/streadway/amqp"
)

type AmqpConnectionProvider struct {
	connectionString string
	connectionAttemptsCount int
}

func NewAmqpConnectionProvider(connectionString string, connectionAttemptsCount int) *AmqpConnectionProvider {
	return &AmqpConnectionProvider{connectionString: connectionString, connectionAttemptsCount: connectionAttemptsCount}
}

func (abc *AmqpConnectionProvider) GetConnection() (
	conn *AmqpConnection,
	err error,
) {
	amqpConnRes, err := connections.WaitForConnection(
		abc.connectionAttemptsCount,
		"RabbitMq",
		func() (interface{}, error) {
			return amqp.Dial(abc.connectionString)
		},
		nil,
	)

	if err != nil {
		return
	}

	amqpConn := amqpConnRes.(*amqp.Connection)
	connClosedNotificationChannel := amqpConn.NotifyClose(make(chan *amqp.Error))
	conn = &AmqpConnection{Conn: amqpConn, ConnClosedNotificationChannel: connClosedNotificationChannel}

	return
}
