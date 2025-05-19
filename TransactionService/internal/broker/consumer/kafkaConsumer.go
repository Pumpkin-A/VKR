package broker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"transaction_service/config"
	"transaction_service/internal/models"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type PaymentManager interface {
	CreatePayment(ctx context.Context, payment models.Payment) (string, error)
	ResultProcessing(ctx context.Context, res models.PaymentResult) error
	MakeRefund(ctx context.Context, payment models.Payment) error
	CancelPayment(ctx context.Context, payment models.Payment) error
}

type Consumer struct {
	gatewayReader *kafka.Reader
	billingReader *kafka.Reader
	pm            PaymentManager
	tracer        trace.Tracer
}

func NewConsumer(ctx context.Context, cfg config.Config, pm PaymentManager) (*Consumer, error) {
	consumer := &Consumer{
		gatewayReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:               []string{cfg.Kafka.Broker1Address},
			GroupID:               cfg.Kafka.ConsumerGroup,
			Topic:                 cfg.Kafka.TopicExternalTransactionOperations,
			WatchPartitionChanges: true,
		}),
		billingReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:               []string{cfg.Kafka.Broker1Address},
			GroupID:               cfg.Kafka.ConsumerGroup,
			Topic:                 cfg.Kafka.TopicInternalPaymentResult,
			WatchPartitionChanges: true,
		}),
		pm:     pm,
		tracer: otel.Tracer("kafka_transaction_service"),
	}

	return consumer, nil
}

func (c *Consumer) Close() error {
	if err := c.gatewayReader.Close(); err != nil {
		slog.Error("failed to close reader:", "err", err.Error())
		return err
	}
	if err := c.billingReader.Close(); err != nil {
		slog.Error("failed to close reader:", "err", err.Error())
		return err
	}
	slog.Info("readers was successfully closed")
	return nil
}

func (c *Consumer) ReadFromGateway(ctx context.Context) {
	for {
		func() {
			msg, err := c.gatewayReader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					slog.Info("[gatewayReader] Consumer context canceled. Exiting...")
					return
				}
				slog.Error("[gatewayReader] could not fetch message", "func", "Consumer: Run", "err", err.Error())
				return
			}

			// Экстрактим tracing контекст из заголовков
			propagator := otel.GetTextMapPropagator()
			carrier := propagation.MapCarrier{}
			for _, header := range msg.Headers {
				carrier.Set(string(header.Key), string(header.Value))
			}

			// Создаем новый контекст с tracing информацией
			ctx = propagator.Extract(ctx, carrier)

			// Создаем span для обработки сообщения
			ctx, span := c.tracer.Start(ctx, "KafkaConsumer.ReadFromGateway")
			defer span.End()

			// Добавляем атрибуты
			span.SetAttributes(
				attribute.String("kafkaTopic", msg.Topic),
				attribute.Int("kafkaPartition", msg.Partition),
				attribute.Int64("kafkaOffset", msg.Offset),
				attribute.String("paymentUUID", string(msg.Key)),
			)

			var event models.EventExternalTransactionOperation
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to unmarshal message")
				slog.Error("[gatewayReader] error with unmarshal msg:", "err", err.Error())
			}

			slog.Info("[gatewayReader] msg was fetch from kafka", "partition: ", msg.Partition, "offset: ", msg.Offset, "payment uuid: ", event.UUID)

			payment := event.ConvertToPayment()
			switch event.TransactionOperation {
			case models.CreateTransactionOperation:
				_, err = c.pm.CreatePayment(ctx, payment)
				if err != nil {
					slog.Error("[gatewayReader] error with createPayment:", "err", err.Error())
				}
			case models.RefundTransactionOperation:
				err = c.pm.MakeRefund(ctx, payment)
				if err != nil {
					slog.Error("[gatewayReader] error with make refund:", "err", err.Error())
				}
			case models.CancelTransactionOperation:
				err = c.pm.CancelPayment(ctx, payment)
				if err != nil {
					slog.Error("[gatewayReader] error with make cancel:", "err", err.Error())
				}
			default:
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				slog.Error("[gatewayReader] unknown transaction operation", "payment uuid:", payment.UUID, "err", err.Error())
				return
			}
			err = c.gatewayReader.CommitMessages(ctx, msg)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				slog.Error("[gatewayReader] error with kafka committing msg", "func", "Consumer: Run", "payment uuid:", payment.UUID, "err", err.Error())
				return
			}
		}()
	}
}

func (c *Consumer) ReadFromBilling(ctx context.Context) {
	for {
		msg, err := c.billingReader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				slog.Info("[billingReader] Consumer context canceled. Exiting...")
				return
			}
			slog.Error("[billingReader] could not fetch message", "func", "Consumer: Run", "err", err.Error())
			return
		}

		var eventResult models.EventInternalPaymentResult
		err = json.Unmarshal(msg.Value, &eventResult)
		if err != nil {
			slog.Error("error with unmarshal msg:", "err", err.Error())
			return
		}

		slog.Info("[billingReader] msg was fetch from kafka", "partition: ", msg.Partition, "offset: ", msg.Offset, "payment uuid: ", eventResult.UUID)

		res := eventResult.ConvertToPaymentResult()
		err = c.pm.ResultProcessing(ctx, res)
		if err != nil {
			slog.Error("error with ResultProcessing:", "err", err.Error())
			return
		}

		err = c.billingReader.CommitMessages(context.Background(), msg)
		if err != nil {
			slog.Error("error with kafka committing msg", "func", "Consumer: Run", "payment uuid:", res.UUID, "err", err.Error())
			return
		}
		// }
	}
}
