package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type Broker struct {
	conn *amqp.Connection
	log  *slog.Logger
}

func NewBroker(url string) (*Broker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Broker{conn: conn}, nil
}

func (b *Broker) Close() error {
	return b.conn.Close()
}

func (b *Broker) CreateChannel() (*amqp.Channel, error) {
	return b.conn.Channel()
}
