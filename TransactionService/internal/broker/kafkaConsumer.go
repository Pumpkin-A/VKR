package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"transaction_service/internal/models"
	pm "transaction_service/internal/paymentManager"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Reader *kafka.Reader
	PM     *pm.PaymentManager
}

func New(topic string, pm pm.PaymentManager) *Consumer {
	c := &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   topic,
			// Partition: 0,
			MaxBytes: 10e6, // 10MB
		}),
		PM: &pm,
	}
	// c.Reader.SetOffset(42)
	return c
}

func (c *Consumer) Close() error {
	if err := c.Reader.Close(); err != nil {
		slog.Error("failed to close reader:", "err", err.Error())
		return err
	}
	return nil
}

func (c *Consumer) Read(ctx context.Context) error {
	fmt.Print("kkjk")
	for {
		m, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			slog.Error("error with reading from kafka:", "err", err.Error())
			return err
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		payment := models.Payment{}
		err = json.Unmarshal(m.Value, &payment)
		if err != nil {
			slog.Error("error with unmarshal msg:", "err", err.Error())
		}
		_, err = c.PM.CreatePayment(ctx, payment)
		if err != nil {
			slog.Error("error with createPaymnet:", "err", err.Error())
		}
	}
}
