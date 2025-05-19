package broker

import (
	"context"
	"encoding/json"
	"log/slog"
	"payment_gateway/config"
	"payment_gateway/internal/models"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Producer struct {
	Writer *kafka.Writer
	tracer trace.Tracer
	// propagator propagation.TextMapPropagator
}

func New(cfg config.Config) *Producer {
	return &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Kafka.BootstrapServers),
			Topic:    cfg.Kafka.Topic,
			Balancer: &kafka.LeastBytes{},
		},
		tracer: otel.Tracer("kafka_gateway"),
		// propagator: otel.GetTextMapPropagator(),
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
	ctx, span := p.tracer.Start(ctx, "KafkaProducer.WriteMessage")
	defer span.End()

	// Добавляем атрибуты для трассировки
	span.SetAttributes(
		attribute.String("kafkaTopic", p.Writer.Topic),
		attribute.String("paymentUUID", payment.UUID),
	)

	paymentByte, err := json.Marshal(payment)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to marshal payment")
		slog.Error("error with marshal payment with uuid", payment.UUID, err.Error())
	}
	// Создаем заголовки для сообщения
	headers := make([]kafka.Header, 0)

	// Инжектим tracing контекст в заголовки Kafka
	carrier := propagation.MapCarrier{}
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, carrier)

	for k, v := range carrier {
		headers = append(headers, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	err = p.Writer.WriteMessages(ctx,
		kafka.Message{
			Key:     []byte(payment.UUID),
			Value:   paymentByte,
			Headers: headers,
		})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to write message")
		slog.Error("failed to write message:", "err", err.Error())
		return err
	}
	slog.Info("succesful writing message to kafka (Gateway -> TS)", "uuid:", payment.UUID)
	return nil
}
