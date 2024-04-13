package rabbitmq

import (
	"banner/internal/domain"
	"banner/internal/repository/rabbitmq"
	"fmt"
	"github.com/streadway/amqp"
)

func SetupRMQ(url string) (*rabbitmq.RabbitMq, error) {
	rmqConn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("can't create connection to rabbitmq: %w", err)
	}

	ch, err := rmqConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("can't create channel for queue declaring %w", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		domain.DeleteQueue,
		false, false, false, false, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("can't declare queue for deleting: %w", err)
	}

	return rabbitmq.New(rmqConn), nil
}
