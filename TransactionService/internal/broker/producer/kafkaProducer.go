package broker

import (
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

// func (p *Producer) WriteCreatePaymentEvent(ctx context.Context, payment models.CreatePaymentEvent) error {
// 	paymentByte, err := json.Marshal(payment)
// 	if err != nil {
// 		slog.Error("error with marshal payment with uuid", payment.UUID, err.Error())
// 	}

// 	err = p.Writer.WriteMessages(ctx,
// 		kafka.Message{
// 			Key:   []byte(payment.UUID),
// 			Value: paymentByte,
// 		})
// 	if err != nil {
// 		slog.Error("failed to write message:", "err", err.Error())
// 		return err
// 	}
// 	slog.Info("succesful writing message to kafka", "uuid:", payment.UUID)
// 	return nil
// }
