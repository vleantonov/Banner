package rabbitmq

import (
	"banner/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMq struct {
	l *zap.Logger
	c *amqp.Connection
}

func New(c *amqp.Connection) *RabbitMq {
	return &RabbitMq{
		c: c,
	}
}

func (r *RabbitMq) Delete(ctx context.Context, tagID, featureID *int) error {
	ch, err := r.c.Channel()
	if err != nil {
		return fmt.Errorf("can't create rmq channel: %w", err)
	}
	defer ch.Close()

	body, err := json.Marshal(
		domain.DeleteBodyQueue{
			TagID:     tagID,
			FeatureID: featureID,
		},
	)

	if err != nil {
		return fmt.Errorf("can't marshal params for queue publishing: %w", err)
	}

	err = ch.Publish(
		"", domain.DeleteQueue,
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("can't publish msg: %w", err)
	}

	return nil
}
