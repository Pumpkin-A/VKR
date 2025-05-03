package paymentmanager

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	models "payment_gateway/internal/models"
)

type Producer interface {
	Write(ctx context.Context, key, value []byte) error
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
	// log.Printf("Добавление транзакции в бд, id: %s\n", payment.UUID)
	// err := pm.DB.AddCardIfNotExist(payment.PaymentMethod.Card)
	// fmt.Println(payment.PaymentMethod.Card.Number)
	// if err != nil {
	// 	return "", err
	// }
	// err = pm.DB.AddPayment(payment)
	// if err != nil {
	// 	return "", err
	// }

	paymentByte, err := json.Marshal(payment)
	if err != nil {
		slog.Error("error with marshal payment with uuid", payment.UUID, err.Error())
	}
	pm.Producer.Write(ctx, []byte(payment.UUID), paymentByte)

	return payment.UUID, nil
}

func (pm *PaymentManager) GetPayment(ctx context.Context, uuid string) (models.Payment, error) {
	log.Printf("Получение транзакции из бд, id: %s\n", uuid)
	// payment, err := pm.DB.GetPayment(uuid)
	// if err != nil {
	// 	return models.Payment{}, err
	// }

	return models.Payment{}, nil
}
