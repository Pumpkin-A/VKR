package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

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

func New(ctx context.Context, topic string, pm PaymentManager) *Consumer {
	consumer := &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   topic,
			// Partition: 0,
			MaxBytes: 10e6, // 10MB
		}),
		pm: pm,
	}
	go consumer.Read(ctx)
	return consumer
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		slog.Error("failed to close reader:", "err", err.Error())
		return err
	}
	return nil
}

func (c *Consumer) Read(ctx context.Context) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			slog.Error("error with reading from kafka:", "err", err.Error())
			return err
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		var payment models.Payment
		err = json.Unmarshal(m.Value, &payment)
		if err != nil {
			slog.Error("error with unmarshal msg:", "err", err.Error())
		}
		_, err = c.pm.CreatePayment(ctx, payment)
		if err != nil {
			slog.Error("error with createPaymnet:", "err", err.Error())
		}
	}
}
