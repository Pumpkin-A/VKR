package broker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"billing/config"
	"billing/internal/models"

	"github.com/segmentio/kafka-go"
)

type PaymentProcessor interface {
	DoPayment(ctx context.Context, payment models.Payment) (string, error)
}

type Consumer struct {
	reader *kafka.Reader
	pp     PaymentProcessor
}

func NewConsumer(ctx context.Context, cfg config.Config, pp PaymentProcessor) (*Consumer, error) {
	consumer := &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:               []string{cfg.Kafka.Broker1Address},
			GroupID:               cfg.Kafka.ConsumerGroup,
			Topic:                 cfg.Kafka.InternalTransactionOperationsTopic,
			WatchPartitionChanges: true,
		}),
		pp: pp,
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

		var event models.EventInternalTransactionOperation
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			slog.Error("error with unmarshal msg:", "err", err.Error())
		}

		slog.Info("msg was fetch from kafka", "partition: ", msg.Partition, "offset: ", msg.Offset, "payment uuid: ", event.UUID)

		payment := event.ConvertToPayment()
		switch event.TransactionOperation {
		case models.CreateTransactionOperation:
			_, err := c.pp.DoPayment(ctx, payment)
			if err != nil {
				slog.Error("error with payment processing", "method", "doPayment", "err", err)
				return
			}
			slog.Info("correct payment processing", "method", "doPayment", "uuid", payment.UUID)
		default:
			slog.Error("unknown transaction operation")
			return
		}

		err = c.reader.CommitMessages(context.Background(), msg)
		if err != nil {
			slog.Error("error with kafka committing msg", "func", "Consumer: Run", "payment uuid:", payment.UUID, "err", err.Error())
			return
		}
	}
}
