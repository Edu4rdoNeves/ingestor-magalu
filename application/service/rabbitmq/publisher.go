package rabbitmq

import "github.com/streadway/amqp"

type PublishQueue struct {
	Exchange  string
	Key       string
	Mandatory bool
	Immediate bool
}

func NewPublish(publishQueue PublishQueue) *PublishQueue {
	return &PublishQueue{
		Exchange:  publishQueue.Exchange,
		Key:       publishQueue.Key,
		Mandatory: publishQueue.Mandatory,
		Immediate: publishQueue.Immediate,
	}
}

func (r *RabbitMQ) Publish(body []byte) error {
	_, err := r.declareQueue(r.declare)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		r.publish.Exchange,
		r.publish.Key,
		r.publish.Mandatory,
		r.publish.Immediate,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
