package amqp

import (
	"fmt"
	"github.com/streadway/amqp"
)

type AmqpMsgAsRawInputMsg struct {
	amqpMsg amqp.Delivery
}

func (arim AmqpMsgAsRawInputMsg) GetPayload() []byte {
	return arim.amqpMsg.Body
}

func (arim AmqpMsgAsRawInputMsg) GetId() string {
	return fmt.Sprint(arim.amqpMsg.DeliveryTag)
}
