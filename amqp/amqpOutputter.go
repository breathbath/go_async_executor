package amqp

import (
	"async_executor/dto"
	"async_executor/logger"
	"fmt"
	"github.com/streadway/amqp"
)

//AmqpOutputter different publishers for error, delayed execution and processable output of async function calls
type AmqpOutputter struct {
	settings     AmqpWritingSettings
	ch           *amqp.Channel
	exchangeName string
}

//NewDelayedAmqpOutputter creates publisher for delayed payloads which should be repeated after the delay time
func NewDelayedAmqpOutputter(id string, settings AmqpWritingSettings, conn *amqp.Connection) (*AmqpOutputter, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	nonDelayedExchangeName := fmt.Sprintf("ex_%s", id)

	queueNameDelayed := fmt.Sprintf("queue_%s_delayed", id)
	q, err := ch.QueueDeclare(
		queueNameDelayed,                                                      // name
		true,                                                           // durable
		false,                                                          // delete when usused
		false,                                                          // exclusive
		false,                                                          // no-wait
		map[string]interface{}{"x-dead-letter-exchange": nonDelayedExchangeName}, // arguments
	)
	if err != nil {
		return nil, err
	}

	exchangeNameDelayed := fmt.Sprintf("ex_%s_delayed", id)

	err = ch.ExchangeDeclare(
		exchangeNameDelayed,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(q.Name, "", exchangeNameDelayed, false, nil)
	if err != nil {
		return nil, err
	}

	return &AmqpOutputter{settings: settings, ch: ch, exchangeName: exchangeNameDelayed}, nil
}

//NewAmqpOutputter publisher to send failed payloads to the error queue
func NewAmqpOutputter(id string, settings AmqpWritingSettings, conn *amqp.Connection) (*AmqpOutputter, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queueName := fmt.Sprintf("queue_%s", id)
	exchangeName := fmt.Sprintf("ex_%s", id)

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &AmqpOutputter{settings: settings, ch: ch, exchangeName: exchangeName}, nil
}

//NewExistingAmqpOutputter publisher to send async function outputs for further processing
func NewExistingAmqpOutputter(conn *amqp.Connection, existingExchangeName string, settings AmqpWritingSettings) (*AmqpOutputter, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		existingExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	return &AmqpOutputter{settings: settings, ch: ch, exchangeName: existingExchangeName}, err
}

//OutputMessage means publish it to the corresponding exchange
func (aot *AmqpOutputter) OutputMessage(msg dto.OutputMessage) error {
	msgString, err := msg.Serialize()
	if err != nil {
		return err
	}

	args := aot.settings.Args

	if args == nil {
		args = &amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: 2,
		}
	}

	expiration := ""
	if aot.settings.LifeTime > 0 {
		expiration = fmt.Sprintf("%.0f", aot.settings.LifeTime.Seconds()*1000)
	}
	args.Expiration = expiration

	args.Body = []byte(msgString)

	err = aot.ch.Publish(
		aot.exchangeName,
		"",
		aot.settings.Mandatory,
		aot.settings.Immediate,
		*args,
	)

	if err == nil {
		expirationLog := ""
		if expiration != "" {
			expirationLog = fmt.Sprintf(" with lifetime of %sms", expiration)
		}
		logger.Log(
			"Published message %s to exchange %s%s",
			msgString,
			aot.exchangeName,
			expirationLog,
		)
	}

	return err
}
