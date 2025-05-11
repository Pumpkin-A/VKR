package paymentManager

import (
	"context"
	"log"
	models "payment_gateway/internal/models"
)

type Producer interface {
	WriteExternalTransactionOperationEvent(ctx context.Context, payment models.ExternalTransactionOperationEvent) error
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
	trOperationEvent := requestData.ConvertToExternalTransactionOperationEvent()

	_ = pm.Producer.WriteExternalTransactionOperationEvent(ctx, trOperationEvent)

	return trOperationEvent.UUID, nil
}

func (pm *PaymentManager) MakeRefund(ctx context.Context, requestData models.MakeRefundRequest) (string, error) {
	trOperationEvent := requestData.ConvertToExternalTransactionOperationEvent()

	_ = pm.Producer.WriteExternalTransactionOperationEvent(ctx, trOperationEvent)

	return trOperationEvent.UUID, nil
}

func (pm *PaymentManager) CancelPayment(ctx context.Context, requestData models.CancelPayment) (string, error) {
	trOperationEvent := requestData.ConvertToExternalTransactionOperationEvent()

	_ = pm.Producer.WriteExternalTransactionOperationEvent(ctx, trOperationEvent)

	return trOperationEvent.UUID, nil
}

func (pm *PaymentManager) GetPayment(ctx context.Context, uuid string) (models.ExternalTransactionOperationEvent, error) {
	log.Printf("Получение транзакции из бд, id: %s\n", uuid)
	// payment, err := pm.DB.GetPayment(uuid)
	// if err != nil {
	// 	return models.Payment{}, err
	// }

	return models.ExternalTransactionOperationEvent{}, nil
}
