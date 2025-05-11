package paymentManager

import (
	"context"
	"log/slog"
	"payment_gateway/internal/grpcClient"
	models "payment_gateway/internal/models"
	pb "payment_gateway/pkg/pb/github.com/yourproject/pkg/pb/transaction/v1"
)

type Client interface {
	GetPayment(ctx context.Context, paymentID string) (*pb.PaymentResponse, error)
}

type Producer interface {
	WriteExternalTransactionOperationEvent(ctx context.Context, payment models.ExternalTransactionOperationEvent) error
}

type PaymentManager struct {
	Producer Producer
	Client   Client
}

func New(producer Producer, client Client) *PaymentManager {
	return &PaymentManager{
		Producer: producer,
		Client:   client,
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

func (pm *PaymentManager) GetPayment(ctx context.Context, uuid string) (models.Payment, error) {
	resp, err := pm.Client.GetPayment(ctx, uuid)
	if err != nil {
		slog.Error("error with grpc request GetPayment", "err", err.Error())
	}
	slog.Info("successful grpc request GetPayment", "uuid", uuid)

	payment := grpcClient.ToDomainPayment(resp)
	slog.Info("successful parsing payment from pb.PaymentResponse", "payment", payment)
	return *payment, nil
}
