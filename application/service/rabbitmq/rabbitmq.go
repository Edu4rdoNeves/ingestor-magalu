package rabbitmq

import "github.com/streadway/amqp"

type IRabbitMQ interface {
	PublishWithNewChannel(body []byte) error
	Consumer(handler func([]byte)) error
	Publish(body []byte) error
	Close()
}

type RabbitMQ struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	declare  QueueDeclare
	consumer QueueConsumer
	publish  PublishQueue
}

func NewRabbitMQ(url string, declare QueueDeclare, consumer QueueConsumer, publish PublishQueue) (IRabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:     conn,
		channel:  ch,
		declare:  declare,
		consumer: consumer,
		publish:  publish,
	}, nil
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}

}
