package paymentmanager

import (
	"context"
	"log"
	models "payment_gateway/internal/models"
)

type Producer interface {
	WriteCreatePaymentEvent(ctx context.Context, payment models.CreatePaymentEvent) error
	Close() error
}

type PaymentManager struct {
	Producer Producer
}

func New(producer Producer) *PaymentManager {
	return &PaymentManager{
		Producer: producer,
	}
}

func (pm *PaymentManager) CreatePayment(ctx context.Context, requestData models.CreatePaymentRequest) (string, error) {
	payment := models.ConvertCreatePaymentRequestToPayment(requestData)

	pm.Producer.WriteCreatePaymentEvent(ctx, payment)

	return payment.UUID, nil
}

func (pm *PaymentManager) GetPayment(ctx context.Context, uuid string) (models.CreatePaymentEvent, error) {
	log.Printf("Получение транзакции из бд, id: %s\n", uuid)
	// payment, err := pm.DB.GetPayment(uuid)
	// if err != nil {
	// 	return models.Payment{}, err
	// }

	return models.CreatePaymentEvent{}, nil
}
