package rabbitmq

import (
	"banner/internal/config"
	"banner/internal/domain"
	"banner/internal/pkg/database"
	"banner/internal/pkg/logger"
	"banner/internal/repository/postresql"
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
)

const workerType = "rabbitmq"

type Repository interface {
	Delete(ctx context.Context, tagID, featureID *int) error
}

type Worker struct {
	l    *zap.Logger
	r    Repository
	Conn *amqp.Connection
}

func New() *Worker {

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("can't create config: %v", err)
	}

	l, err := logger.New()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	l = l.With(zap.String("worker", workerType))

	db, err := database.CreateDBConnection(cfg.StorageCfg.PGUrl, l)
	if err != nil {
		l.Fatal("can't create db connection")
	}
	r := postresql.New(db)

	conn, err := amqp.Dial(cfg.QueueCfg.RMQUrl)
	if err != nil {
		l.Fatal("can't create connection with rabbitmq", zap.Error(err))
	}

	return &Worker{
		l:    l,
		r:    r,
		Conn: conn,
	}
}

func (w *Worker) MustRun() {

	defer w.Conn.Close()
	w.l.Info("try to run worker")

	ch, err := w.Conn.Channel()
	if err != nil {
		w.l.Fatal("can't create channel", zap.Error(err))
	}
	defer ch.Close()

	msgs, err := ch.Consume(domain.DeleteQueue, "",
		true, false, false, false, nil,
	)

	ctx := context.Background()
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			w.l.Info(
				"received message",
				zap.ByteString("msg", d.Body),
				zap.String("msg_id", d.MessageId),
			)
			var b domain.DeleteBodyQueue
			err := json.Unmarshal(d.Body, &b)
			if err != nil {
				w.l.Error("can't unmarshal message", zap.Error(err), zap.String("msg_id", d.MessageId))
			}

			err = w.r.Delete(ctx, b.TagID, b.FeatureID)
			if err != nil {
				w.l.Error("can't delete records", zap.Error(err), zap.String("msg_id", d.MessageId))
			}
		}
	}()

	w.l.Info("Successfully run worker")
	w.l.Info("[*] - waiting for messages")
	<-forever
}
