package builder

import "github.com/breathbath/go_async_executor/amqp/connection"

func BuildAmqpConnection(connectionString string, connectionAttemptsCount int) (conn *connection.AmqpConnection, err error) {
	amqpConnProvider := connection.NewAmqpConnectionProvider(
		connectionString,
		connectionAttemptsCount,
	)
	conn, err = amqpConnProvider.GetConnection()

	return
}
