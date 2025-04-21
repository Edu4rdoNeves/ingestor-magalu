package rabbitmq

import "github.com/streadway/amqp"

type QueueConsumer struct {
	QueueName string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

func NewConsumerQueue(queueConsumer QueueConsumer, queueDeclare QueueDeclare) *QueueConsumer {
	return &QueueConsumer{
		QueueName: queueConsumer.QueueName,
		Consumer:  queueConsumer.Consumer,
		AutoAck:   queueConsumer.AutoAck,
		Exclusive: queueConsumer.Exclusive,
		NoLocal:   queueConsumer.NoLocal,
		NoWait:    queueConsumer.NoWait,
		Args:      queueConsumer.Args,
	}
}

func (r *RabbitMQ) Consumer(handler func([]byte)) error {
	queue, err := r.declareQueue(r.declare)
	if err != nil {
		return err
	}

	msgs, err := r.declareConsumer(*queue, r.consumer)
	if err != nil {
		return err
	}

	go func() {
		for d := range *msgs {
			handler(d.Body)
		}
	}()

	return nil
}

func (r *RabbitMQ) declareConsumer(queue amqp.Queue, queueConsumer QueueConsumer) (*<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		queue.Name,
		queueConsumer.Consumer,
		queueConsumer.AutoAck,
		queueConsumer.Exclusive,
		queueConsumer.NoLocal,
		queueConsumer.NoWait,
		queueConsumer.Args,
	)
	if err != nil {
		return nil, err
	}

	return &msgs, nil
}
