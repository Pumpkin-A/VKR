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
	WriterInternalTransactionOperations *kafka.Writer
	WriterExternalPaymentResult         *kafka.Writer
}

func NewProducer(cfg config.Config) *Producer {
	return &Producer{
		WriterInternalTransactionOperations: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Kafka.Broker1Address),
			Topic:    cfg.Kafka.TopicInternalTransactionOperations,
			Balancer: &kafka.LeastBytes{},
		},
		WriterExternalPaymentResult: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Kafka.Broker1Address),
			Topic:    cfg.Kafka.TopicTransactionFinishStatus,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Close() error {
	if err := p.WriterInternalTransactionOperations.Close(); err != nil {
		slog.Error("[WriterInternalTransactionOperations] failed to close writer:", "err", err.Error())
		return err
	}
	if err := p.WriterExternalPaymentResult.Close(); err != nil {
		slog.Error("[WriterTransactionFinishStatus] failed to close writer:", "err", err.Error())
		return err
	}
	return nil
}

func (p *Producer) WriteEventInternalTransactionOperation(ctx context.Context, event models.EventInternalTransactionOperation) error {
	eventByte, err := json.Marshal(event)
	if err != nil {
		slog.Error("[WriterInternalTransactionOperations] error with marshal payment with uuid", event.UUID, err.Error())
	}

	err = p.WriterInternalTransactionOperations.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(event.UUID),
			Value: eventByte,
		})
	if err != nil {
		slog.Error("[WriterInternalTransactionOperations] failed to write message:", "err", err.Error())
		return err
	}
	slog.Info("[WriterInternalTransactionOperations] succesful writing message to kafka (TS -> Billing)", "uuid:", event.UUID)
	return nil
}

func (p *Producer) WriteEventExternalPaymentResult(ctx context.Context, event models.EventExternalPaymentResult) error {
	eventByte, err := json.Marshal(event)
	if err != nil {
		slog.Error("[WriteEventExternalPaymentResult] error with marshal payment with uuid", event.UUID, err.Error())
	}

	err = p.WriterExternalPaymentResult.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(event.UUID),
			Value: eventByte,
		})
	if err != nil {
		slog.Error("[WriteEventExternalPaymentResult] failed to write message:", "err", err.Error())
		return err
	}
	slog.Info("[WriteEventExternalPaymentResult] succesful writing message to kafka (TS -> ALL)", "uuid:", event.UUID, "status", event.Status)
	return nil
}
