package broker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"transaction_service/config"
	"transaction_service/internal/models"

	"github.com/segmentio/kafka-go"
)

type PaymentManager interface {
	CreatePayment(ctx context.Context, payment models.Payment) (string, error)
	ResultProcessing(ctx context.Context, res models.PaymentResult) error
}

type Consumer struct {
	gatewayReader *kafka.Reader
	billingReader *kafka.Reader
	pm            PaymentManager
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
		pm: pm,
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
		msg, err := c.gatewayReader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				slog.Info("[gatewayReader] Consumer context canceled. Exiting...")
				return
			}
			slog.Error("[gatewayReader] could not fetch message", "func", "Consumer: Run", "err", err.Error())
			return
		}

		var event models.EventExternalTransactionOperation
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
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
		default:
			slog.Error("[gatewayReader] unknown transaction operation", "payment uuid:", payment.UUID, "err", err.Error())
			return
		}
		err = c.gatewayReader.CommitMessages(context.Background(), msg)
		if err != nil {
			slog.Error("[gatewayReader] error with kafka committing msg", "func", "Consumer: Run", "payment uuid:", payment.UUID, "err", err.Error())
			return
		}
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
