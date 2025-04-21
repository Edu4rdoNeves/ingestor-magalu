package rabbitmq

import "github.com/streadway/amqp"

func (r *RabbitMQ) PublishWithNewChannel(body []byte) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		r.declare.Name,
		r.declare.Durable,
		r.declare.AutoDelete,
		r.declare.Exclusive,
		r.declare.NoWait,
		r.declare.Args,
	)
	if err != nil {
		return err
	}

	return ch.Publish(
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
