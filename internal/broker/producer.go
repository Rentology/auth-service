package broker

import amqp "github.com/rabbitmq/amqp091-go"

type Producer struct {
	channel *amqp.Channel
	queue   string
}

func NewProducer(b *Broker, queueName string) (*Producer, error) {
	ch, err := b.CreateChannel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // Устойчивая очередь
		false, // Удалять неиспользуемую очередь
		false, // Неэксклюзивная очередь
		false, // Без ожидания
		nil,   // Дополнительные аргументы
	)
	if err != nil {
		return nil, err
	}

	return &Producer{channel: ch, queue: queueName}, nil
}

func (p *Producer) Publish(body []byte) error {
	err := p.channel.Publish(
		"",      // Exchange
		p.queue, // Routing key
		false,   // Mandatory
		false,   // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.channel.Close()
}
