package broker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"billing/config"
	"billing/internal/models"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type PaymentProcessor interface {
	DoPayment(ctx context.Context, payment models.Payment) (string, error)
	DoRefund(ctx context.Context, payment models.Payment) (string, error)
	CancelPayment(ctx context.Context, payment models.Payment) (string, error)
}

type Consumer struct {
	reader *kafka.Reader
	pp     PaymentProcessor
	tracer trace.Tracer
}

func NewConsumer(ctx context.Context, cfg config.Config, pp PaymentProcessor) (*Consumer, error) {
	consumer := &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:               []string{cfg.Kafka.Broker1Address},
			GroupID:               cfg.Kafka.ConsumerGroup,
			Topic:                 cfg.Kafka.InternalTransactionOperationsTopic,
			WatchPartitionChanges: true,
		}),
		pp:     pp,
		tracer: otel.Tracer("kafka_consumer_billing"),
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
		if err := func() error {
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					slog.Info("Consumer context canceled. Exiting...")
					return err
				}
				slog.Error("could not fetch message", "func", "Consumer: Run", "err", err.Error())
				return err
			}

			ctx, span := c.interceptorForKafkaComsumer(ctx, msg)
			defer span.End()

			var event models.EventInternalTransactionOperation
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to unmarshal message")
				slog.Error("error with unmarshal msg:", "err", err.Error())
			}

			slog.Info("msg was fetch from kafka", "partition: ", msg.Partition, "offset: ", msg.Offset, "payment uuid: ", event.UUID)
			payment := event.ConvertToPayment()
			go func() {
				switch event.TransactionOperation {
				case models.CreateTransactionOperation:
					_, err := c.pp.DoPayment(ctx, payment)
					if err != nil {
						slog.Error("error with payment processing", "method", "doPayment", "err", err)
						return
					}
					slog.Info("correct payment processing", "method", "doPayment", "uuid", payment.UUID)
				case models.RefundTransactionOperation:
					_, err := c.pp.DoRefund(ctx, payment)
					if err != nil {
						slog.Error("error with payment processing", "method", "doRefund", "err", err)
						return
					}
					slog.Info("correct payment processing", "method", "doRefund", "uuid", payment.UUID)
				case models.CancelTransactionOperation:
					_, err := c.pp.CancelPayment(ctx, payment)
					if err != nil {
						slog.Error("error with payment processing", "method", "doRefund", "err", err)
						return
					}
					slog.Info("correct payment processing", "method", "CancelPayment", "uuid", payment.UUID)
				default:
					slog.Error("unknown transaction operation")
					return
				}

			}()

			err = c.reader.CommitMessages(context.Background(), msg)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to unmarshal message")
				slog.Error("error with kafka committing msg", "func", "Consumer: Run", "payment uuid:", payment.UUID, "err", err.Error())
				return err
			}
			return nil
		}(); err != nil {
			return
		}
	}
}

func (c *Consumer) interceptorForKafkaComsumer(ctx context.Context, msg kafka.Message) (context.Context, trace.Span) {
	// Экстрактим tracing контекст из заголовков
	propagator := otel.GetTextMapPropagator()
	carrier := propagation.MapCarrier{}
	for _, header := range msg.Headers {
		carrier.Set(string(header.Key), string(header.Value))
	}

	// Создаем новый контекст с tracing информацией
	ctx = propagator.Extract(ctx, carrier)

	// Создаем span для обработки сообщения
	ctx, span := c.tracer.Start(ctx, "kafkaConsumer.ReadFromTransactionService")

	// Добавляем атрибуты
	span.SetAttributes(
		attribute.String("kafkaTopic", msg.Topic),
		attribute.Int("kafkaPartition", msg.Partition),
		attribute.Int64("kafkaOffset", msg.Offset),
		attribute.String("paymentUUID", string(msg.Key)),
	)

	return ctx, span
}
