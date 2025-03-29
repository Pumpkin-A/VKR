package paymentmanager

import (
	"context"
	"fmt"
	"log"
	models "payment_gateway/internal/models"
)

type DB interface {
	AddPayment(models.Payment) error
	GetPayment(uuid string) (models.Payment, error)
	AddCardIfNotExist(c models.Card) error
}

type PaymentManager struct {
	DB DB
}

func New(db DB) *PaymentManager {
	return &PaymentManager{DB: db}
}

func (pm *PaymentManager) CreatePayment(ctx context.Context, requestData models.CreatePaymentRequest) (string, error) {
	payment := models.ConvertCreatePaymentRequestToPayment(requestData)
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
