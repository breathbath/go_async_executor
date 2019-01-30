package amqp

import (
	"go_async_executor/dto"
	"go_async_executor/logger"
	"fmt"
	"github.com/streadway/amqp"
)
//AmqpAsyncFuncRegistrator publishes async func payloads to the processing queue
type AmqpAsyncFuncRegistrator struct {
	inputQueueName    string
	inputExchangeName string
	conn              *amqp.Connection
	channel           *amqp.Channel
	writingSettings   AmqpWritingSettings
}

func NewAmqpAsyncFuncRegistrator(
	sourceId string,
	writingSettings AmqpWritingSettings,
	conn *amqp.Connection,
) (*AmqpAsyncFuncRegistrator, error) {
	inputQueueName := fmt.Sprintf("queue_%s", sourceId)
	inputExchangeName := fmt.Sprintf("ex_%s", sourceId)

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		inputExchangeName, // name
		"direct",          // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)

	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		inputQueueName, // name
		true,           // durable
		false,          // delete when usused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,            // queue name
		"",                // routing key
		inputExchangeName, // exchange
		false,
		nil,
	)

	return &AmqpAsyncFuncRegistrator{
		inputQueueName:    inputQueueName,
		inputExchangeName: inputExchangeName,
		conn:              conn,
		channel:           ch,
		writingSettings:   writingSettings,
	}, nil
}

func (ais *AmqpAsyncFuncRegistrator) RegisterAsyncExecution(msg dto.AsyncFuncInput) error {
	msgString, err := msg.Serialize()
	if err != nil {
		return err
	}

	args := ais.writingSettings.Args

	if args == nil {
		args = &amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: 2,
		}
	}

	args.Body = []byte(msgString)

	err = ais.channel.Publish(
		ais.inputExchangeName,
		"",
		ais.writingSettings.Mandatory,
		ais.writingSettings.Immediate,
		*args,
	)

	if err == nil {
		logger.Log(
			"Published message %s to exchange %s",
			msgString,
			ais.inputExchangeName,
		)
	}

	return err
}
