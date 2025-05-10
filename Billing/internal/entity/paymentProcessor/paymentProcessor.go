package paymentProcessor

import (
	models "billing/internal/models"
	"context"
	"log/slog"
)

type Client interface {
	DoPayment(payment models.Payment) (models.ResultOfRequestFromBank, error)
}

type Producer interface {
	WriteEventInternalPaymentResult(ctx context.Context, event models.EventInternalPaymentResult) error
	Close() error
}

type PaymentProcessor struct {
	client   Client
	producer Producer
}

func New(c Client, producer Producer) *PaymentProcessor {
	return &PaymentProcessor{
		client:   c,
		producer: producer,
	}
}

func (pm *PaymentProcessor) DoPayment(ctx context.Context, payment models.Payment) (string, error) {
	result, err := pm.client.DoPayment(payment)
	if err != nil {
		slog.Error("error with http request to bank_example", "err", err)
		return "", err
	}
	slog.Info("result from bank_example was received", "UUID", result.UUID, "operation", "CREATE", "statusResult", result.Status)

	// пишем событие другим
	event := result.ConvertToEventInternalPaymentResult(models.CreateTransactionOperation)
	err = pm.producer.WriteEventInternalPaymentResult(ctx, event)
	if err != nil {
		return "", err
	}

	return payment.UUID, nil
}
