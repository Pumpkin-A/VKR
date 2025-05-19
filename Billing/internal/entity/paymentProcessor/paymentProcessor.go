package paymentProcessor

import (
	models "billing/internal/models"
	"context"
	"log/slog"
	"math/rand/v2"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Client interface {
	DoPayment(ctx context.Context, payment models.Payment) (models.ResultOfRequestFromBank, error)
}

type Producer interface {
	WriteEventInternalPaymentResult(ctx context.Context, event models.EventInternalPaymentResult) error
	Close() error
}

type PaymentProcessor struct {
	client   Client
	producer Producer
	tracer   trace.Tracer
}

func New(c Client, producer Producer) *PaymentProcessor {
	return &PaymentProcessor{
		client:   c,
		producer: producer,
		tracer:   otel.Tracer("paymentProcessor_billing"),
	}
}

func (pm *PaymentProcessor) DoPayment(ctx context.Context, payment models.Payment) (string, error) {
	ctx, sp := pm.tracer.Start(ctx, "paymentProcessor.DoPayment")
	sp.SetAttributes(attribute.String("paymentId", payment.UUID))
	defer sp.End()

	result, err := pm.client.DoPayment(ctx, payment)
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

func (pm *PaymentProcessor) DoRefund(ctx context.Context, payment models.Payment) (string, error) {
	result := models.ResultOfRequestFromBank{
		UUID:      payment.UUID,
		Status:    string(randomStatus()),
		ErrorText: "",
	}

	event := result.ConvertToEventInternalPaymentResult(models.RefundTransactionOperation)
	err := pm.producer.WriteEventInternalPaymentResult(ctx, event)
	if err != nil {
		return "", err
	}

	return payment.UUID, nil
}

func (pm *PaymentProcessor) CancelPayment(ctx context.Context, payment models.Payment) (string, error) {
	result := models.ResultOfRequestFromBank{
		UUID:      payment.UUID,
		Status:    string(randomStatus()),
		ErrorText: "",
	}

	event := result.ConvertToEventInternalPaymentResult(models.CancelTransactionOperation)
	err := pm.producer.WriteEventInternalPaymentResult(ctx, event)
	if err != nil {
		return "", err
	}

	return payment.UUID, nil
}

func randomStatus() models.BankExampleStatus {
	var status models.BankExampleStatus
	randNum := rand.Float64()
	switch {
	case randNum <= 0.6: // 60% chance
		status = models.SuccessedBankExampleStatus
	case randNum <= 0.9: // 30% chance (0.6 + 0.3)
		status = models.FailedBankExampleStatus
	default: // 10% chance
		status = models.ErrorBankExampleStatus
	}
	return status
}
