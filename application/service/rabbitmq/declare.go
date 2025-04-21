package rabbitmq

import "github.com/streadway/amqp"

type QueueDeclare struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func NewDeclareQueue(queueDeclare QueueDeclare) *QueueDeclare {
	return &QueueDeclare{
		Name:       queueDeclare.Name,
		Durable:    queueDeclare.Durable,
		AutoDelete: queueDeclare.AutoDelete,
		Exclusive:  queueDeclare.Exclusive,
		NoWait:     queueDeclare.NoWait,
		Args:       queueDeclare.Args,
	}
}

func (r *RabbitMQ) declareQueue(queueDeclare QueueDeclare) (*amqp.Queue, error) {
	queue, err := r.channel.QueueDeclare(
		queueDeclare.Name,
		queueDeclare.Durable,
		queueDeclare.AutoDelete,
		queueDeclare.Exclusive,
		queueDeclare.NoWait,
		queueDeclare.Args,
	)
	if err != nil {
		return nil, err
	}

	return &queue, nil
}
