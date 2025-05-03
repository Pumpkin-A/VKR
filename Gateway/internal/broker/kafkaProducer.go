package broker

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func New(BootstrapServers, topic string) *Producer {
	return &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(BootstrapServers),
			Topic:    topic,
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

func (p *Producer) Write(ctx context.Context, key, value []byte) error {
	err := p.Writer.WriteMessages(ctx,
		kafka.Message{
			Key:   key,
			Value: value,
		})
	if err != nil {
		slog.Error("failed to write message:", "err", err.Error())
		return err
	}
	slog.Info("succesful writing message to kafka", "uuid:", string(key))
	return nil
}
