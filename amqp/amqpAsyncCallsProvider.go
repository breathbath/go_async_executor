package amqp

import (
	"async_executor/dto"
	"async_executor/logger"
	"fmt"
	"github.com/streadway/amqp"
	"strconv"
)

type AmqpAsyncCallsProvider struct {
	readingSettings   *InputReadingSettings
	inputQueueName    string
	inputExchangeName string
	conn              *amqp.Connection
	channel           *amqp.Channel
}

func NewAmqpAsyncCallsProvider(
	sourceId string,
	readingSettings *InputReadingSettings,
	conn *amqp.Connection,
) (*AmqpAsyncCallsProvider, error) {
	inputQueueName := fmt.Sprintf("queue_%s", sourceId)
	inputExchangeName := fmt.Sprintf("ex_%s", sourceId)

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.Qos(10, 0, false)
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

	return &AmqpAsyncCallsProvider{
		readingSettings:   readingSettings,
		inputQueueName:    inputQueueName,
		inputExchangeName: inputExchangeName,
		conn:              conn,
		channel:           ch,
	}, nil
}

func (ais *AmqpAsyncCallsProvider) MarkAsDone(msg dto.RawInput) error {
	idInt, err := strconv.ParseUint(msg.GetId(), 10, 64)
	if err != nil {
		return err
	}

	return ais.channel.Ack(idInt, false)
}

func (ais *AmqpAsyncCallsProvider) GetAsyncCalls() (<-chan dto.RawInput, error) {
	logger.Log("Will fetch amqp messages from channel '%s'", ais.inputQueueName)
	amqpCh, err := ais.channel.Consume(
		ais.inputQueueName,
		ais.readingSettings.ConsumerName,
		ais.readingSettings.AutoAck,
		ais.readingSettings.Exclusive,
		ais.readingSettings.NoLocal,
		ais.readingSettings.NoWait,
		ais.readingSettings.Args,
	)

	if err != nil {
		return nil, err
	}

	rawMsgChannel := make(chan dto.RawInput, 10000)

	go func() {
		for amqpMsg := range amqpCh {
			rawMsgChannel <- AmqpMsgAsRawInputMsg{amqpMsg:amqpMsg}
		}
	}()

	return rawMsgChannel, nil
}
