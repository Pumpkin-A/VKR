package broker

import (
	"context"
	"encoding/json"
	"log/slog"
	"payment_gateway/config"
	"payment_gateway/internal/models"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func New(cfg config.Config) *Producer {
	return &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Kafka.BootstrapServers),
			Topic:    cfg.Kafka.Topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Close() error {
	if err := p.Writer.Close(); err != nil {
		slog.Error("failed to close writer:", "err", err.Error())
		return err
	}
	return nil
}

func (p *Producer) WriteExternalTransactionOperationEvent(ctx context.Context, payment models.ExternalTransactionOperationEvent) error {
	paymentByte, err := json.Marshal(payment)
	if err != nil {
		slog.Error("error with marshal payment with uuid", payment.UUID, err.Error())
	}

	err = p.Writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(payment.UUID),
			Value: paymentByte,
		})
	if err != nil {
		slog.Error("failed to write message:", "err", err.Error())
		return err
	}
	slog.Info("succesful writing message to kafka (Gateway -> TS)", "uuid:", payment.UUID)
	return nil
}
