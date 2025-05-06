package broker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"transaction_service/config"
	"transaction_service/internal/models"

	"github.com/segmentio/kafka-go"
)

type PaymentManager interface {
	CreatePayment(ctx context.Context, payment models.Payment) (string, error)
}

type Consumer struct {
	reader *kafka.Reader
	pm     PaymentManager
}

func New(ctx context.Context, cfg config.Config, pm PaymentManager) (*Consumer, error) {
	consumer := &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:               []string{cfg.Kafka.Broker1Address},
			GroupID:               cfg.Kafka.ConsumerGroup,
			Topic:                 cfg.Kafka.Topic,
			WatchPartitionChanges: true,
		}),
		pm: pm,
	}

	go consumer.Run(ctx)

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
		// the `FetchMessage` method blocks until we receive the next event
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				slog.Info("Consumer context canceled. Exiting...")
				return
			}
			slog.Error("could not fetch message", "func", "Consumer: Run", "err", err.Error())
			return
		}

		var payment models.Payment
		err = json.Unmarshal(msg.Value, &payment)
		if err != nil {
			slog.Error("error with unmarshal msg:", "err", err.Error())
		}

		slog.Info("msg was fetch from kafka", "partition: ", msg.Partition, "offset: ", msg.Offset, "payment uuid: ", payment.UUID)

		_, err = c.pm.CreatePayment(ctx, payment)
		if err != nil {
			slog.Error("error with createPayment:", "err", err.Error())
		}

		err = c.reader.CommitMessages(context.Background(), msg)
		if err != nil {
			slog.Error("error with kafka committing msg", "func", "Consumer: Run", "payment uuid:", payment.UUID, "err", err.Error())
			return
		}
	}
}
