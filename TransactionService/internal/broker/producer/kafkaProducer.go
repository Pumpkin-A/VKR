package broker

import (
	"context"
	"encoding/json"
	"log/slog"
	"transaction_service/config"
	"transaction_service/internal/models"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(cfg config.Config) *Producer {
	return &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Kafka.Broker1Address),
			Topic:    cfg.Kafka.InternalTransactionOperationsTopic,
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

func (p *Producer) WriteInternalTransactionOperationEvent(ctx context.Context, event models.InternalTransactionOperationEvent) error {
	eventByte, err := json.Marshal(event)
	if err != nil {
		slog.Error("error with marshal payment with uuid", event.UUID, err.Error())
	}

	err = p.Writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(event.UUID),
			Value: eventByte,
		})
	if err != nil {
		slog.Error("failed to write message:", "err", err.Error())
		return err
	}
	slog.Info("succesful writing message to kafka (TS -> Billing)", "uuid:", event.UUID)
	return nil
}
