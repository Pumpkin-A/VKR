package paymentManager

import (
	"context"
	"fmt"
	"log"
	models "transaction_service/internal/models"
)

type DB interface {
	AddPayment(models.Payment) error
	GetPayment(uuid string) (models.Payment, error)
	AddCardIfNotExist(c models.Card) error
}

type Producer interface {
	WriteInternalTransactionOperationEvent(ctx context.Context, event models.InternalTransactionOperationEvent) error
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
	event := models.ConvertPaymentToInternalTrasactionOperationEvent(payment)
	event.TransactionOperation = models.CreateTransactionOperation
	_ = pm.Producer.WriteInternalTransactionOperationEvent(ctx, event)

	return payment.UUID, nil
}

func (pm *PaymentManager) GetPayment(ctx context.Context, uuid string) (models.Payment, error) {
	log.Printf("Получение транзакции из бд, id: %s\n", uuid)
	payment, err := pm.DB.GetPayment(uuid)
	if err != nil {
		return models.Payment{}, err
	}

	return payment, nil
}
