package broker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"statistics/config"
	"statistics/internal/models"

	"github.com/segmentio/kafka-go"
)

type statisticsManager interface {
	Analize(ctx context.Context)
}

type Consumer struct {
	reader *kafka.Reader
	sm     statisticsManager
}

func NewConsumer(ctx context.Context, cfg config.Config, sm statisticsManager) (*Consumer, error) {
	consumer := &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:               []string{cfg.Kafka.Broker1Address},
			GroupID:               cfg.Kafka.ConsumerGroup,
			Topic:                 cfg.Kafka.Topic,
			WatchPartitionChanges: true,
		}),
		sm: sm,
	}

	return consumer, nil
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		slog.Error("failed to close reader:", "err", err.Error())
		return err
	}
	slog.Info("reader was successfully closed")
	return nil
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				slog.Info("Consumer context canceled. Exiting...")
				return
			}
			slog.Error("could not fetch message", "func", "Consumer: Run", "err", err.Error())
			return
		}

		var event models.EventExternalPaymentResult
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			slog.Error("error with unmarshal msg:", "err", err.Error())
		}

		slog.Info("msg was fetch from kafka", "partition: ", msg.Partition, "offset: ", msg.Offset, "EVENT", event)

		// обработка сообщения с отправкой синхронных запросов сторонним сервисам
		c.sm.Analize(ctx)

		err = c.reader.CommitMessages(context.Background(), msg)
		if err != nil {
			slog.Error("error with kafka committing msg", "func", "Consumer: Run", "payment uuid:", event.UUID, "err", err.Error())
			return
		}
	}
}
