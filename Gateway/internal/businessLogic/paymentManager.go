package paymentmanager

import (
	"context"
	"log"
	models "payment_gateway/internal/models"
)

type DB interface {
	AddPayment(models.Payment) error
	GetPayment(id int) models.Payment
}

type PaymentManager struct {
	DB DB
}

func New(db DB) *PaymentManager {
	return &PaymentManager{DB: db}
}

func (pm *PaymentManager) CreatePayment(ctx context.Context, requestData models.CreatePaymentRequest) (string, error) {
	payment := models.ConvertCreatePaymentRequestToPayment(requestData)
	log.Printf("Добавление транзакции в бд, id: %s\n", payment.ID)
	return payment.ID, nil
}

func (pm *PaymentManager) GetPaymentInfo(ctx context.Context, id string) (models.Payment, error) {
	log.Printf("Получение транзакции из бд, id: %s\n", id)
	return models.Payment{ID: id}, nil
}
