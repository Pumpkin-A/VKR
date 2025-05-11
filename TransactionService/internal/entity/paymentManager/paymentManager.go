package paymentManager

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	models "transaction_service/internal/models"
)

type DB interface {
	AddPayment(models.Payment) error
	GetPayment(uuid string) (models.Payment, error)
	AddCardIfNotExist(c models.Card) error
	UpdatePaymentStatus(uuid, status string) error
	GetPaymentStatus(uuid string) (models.PaymentStatus, error)
}

type Producer interface {
	WriteEventInternalTransactionOperation(ctx context.Context, event models.EventInternalTransactionOperation) error
	WriteEventExternalPaymentResult(ctx context.Context, event models.EventExternalPaymentResult) error
	Close() error
}

type PaymentManager struct {
	DB       DB
	Producer Producer
}

func New(db DB, producer Producer) *PaymentManager {
	return &PaymentManager{
		DB:       db,
		Producer: producer,
	}
}

func (pm *PaymentManager) CreatePayment(ctx context.Context, payment models.Payment) (string, error) {
	log.Printf("Добавление транзакции в бд, id: %s\n", payment.UUID)
	err := pm.DB.AddCardIfNotExist(payment.PaymentMethod.Card)
	fmt.Println(payment.PaymentMethod.Card.Number)
	if err != nil {
		return "", err
	}
	err = pm.DB.AddPayment(payment)
	if err != nil {
		return "", err
	}

	// пишем событие другим
	event := payment.ConvertToInternalTrasactionOperationEvent(models.CreateTransactionOperation)
	_ = pm.Producer.WriteEventInternalTransactionOperation(ctx, event)

	return payment.UUID, nil
}

func (pm *PaymentManager) MakeRefund(ctx context.Context, payment models.Payment) error {
	curStatus, err := pm.DB.GetPaymentStatus(payment.UUID)
	if err != nil {
		slog.Error("error with get transaction status from db", "err", err.Error())
		return err
	}

	if curStatus != models.InProcessingPaymentStatus {
		slog.Info("untimely action", "payment uuid:", payment.UUID, "status", curStatus)
		return nil
	}

	event := payment.ConvertToInternalTrasactionOperationEvent(models.RefundTransactionOperation)
	err = pm.Producer.WriteEventInternalTransactionOperation(ctx, event)
	if err != nil {
		slog.Error("error with WriteEventInternalTransactionOperation", "err", err.Error())
		return err
	}

	return nil
}

func (pm *PaymentManager) CancelPayment(ctx context.Context, payment models.Payment) error {
	curStatus, err := pm.DB.GetPaymentStatus(payment.UUID)
	if err != nil {
		slog.Error("error with get transaction status from db", "err", err.Error())
		return err
	}

	if curStatus != models.InProcessingPaymentStatus {
		slog.Info("untimely action", "payment uuid:", payment.UUID, "status", curStatus)
		return nil
	}

	event := payment.ConvertToInternalTrasactionOperationEvent(models.CancelTransactionOperation)
	err = pm.Producer.WriteEventInternalTransactionOperation(ctx, event)
	if err != nil {
		slog.Error("error with WriteEventInternalTransactionOperation", "err", err.Error())
		return err
	}

	return nil
}

func (pm *PaymentManager) setFinalTrsnsactionStatus(ctx context.Context, event models.EventExternalPaymentResult) error {
	err := pm.DB.UpdatePaymentStatus(event.UUID, string(event.Status))
	if err != nil {
		return err
	}

	// массовая рассылОчка файнал статус
	err = pm.Producer.WriteEventExternalPaymentResult(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (pm *PaymentManager) ResultProcessing(ctx context.Context, res models.PaymentResult) error {
	curStatus, err := pm.DB.GetPaymentStatus(res.UUID)
	if err != nil {
		slog.Error("error with get transaction status from db", "err", err.Error())
		return err
	}

	var event models.EventExternalPaymentResult
	switch {
	case res.TransactionOperation == models.CreateTransactionOperation && curStatus == models.InProcessingPaymentStatus:
		if res.Status == models.SuccessedBankStatus {
			event = res.ConvertToEventExternalPaymentResult(models.SuccessPaymentStatus)
		} else if res.Status == models.FailedBankStatus {
			event = res.ConvertToEventExternalPaymentResult(models.FailedPaymentStatus)
		} else {
			event = res.ConvertToEventExternalPaymentResult(models.ErrorPaymentStatus)
		}

	case res.TransactionOperation == models.RefundTransactionOperation && curStatus == models.InProcessingPaymentStatus && res.Status == models.SuccessedBankStatus:
		event = res.ConvertToEventExternalPaymentResult(models.RefundedPaymentStatus)

	case res.TransactionOperation == models.CancelTransactionOperation && curStatus == models.InProcessingPaymentStatus && res.Status == models.SuccessedBankStatus:
		event = res.ConvertToEventExternalPaymentResult(models.CancelledPaymentStatus)
	default:
		slog.Info("untimely action", "payment uuid:", res.UUID, "curStatus", curStatus, "res.Status", res.Status)
		return nil
	}

	err = pm.setFinalTrsnsactionStatus(ctx, event)
	if err != nil {
		slog.Error("error with SetFinalTrsnsactionStatus", "err", err.Error())
		return err
	}

	return nil
}

func (pm *PaymentManager) GetPayment(ctx context.Context, uuid string) (models.Payment, error) {
	log.Printf("Получение транзакции из бд, id: %s\n", uuid)
	payment, err := pm.DB.GetPayment(uuid)
	if err != nil {
		return models.Payment{}, err
	}

	return payment, nil
}
