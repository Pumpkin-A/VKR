package grpcClient

import (
	"time"

	"payment_gateway/internal/models"
	pb "payment_gateway/pkg/pb/github.com/yourproject/pkg/pb/transaction/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToDomainPayment конвертирует PaymentResponse gRPC в доменную модель Payment
func ToDomainPayment(grpcPayment *pb.PaymentResponse) *models.Payment {
	if grpcPayment == nil {
		return nil
	}

	return &models.Payment{
		UUID:          grpcPayment.GetUuid(),
		Status:        toDomainPaymentStatus(grpcPayment.GetStatus()),
		Paid:          grpcPayment.GetPaid(),
		Amount:        toDomainAmount(grpcPayment.GetAmount()),
		CreatedAt:     toDomainTime(grpcPayment.GetCreatedAt()),
		Description:   grpcPayment.GetDescription(),
		ExpiresAt:     toDomainTime(grpcPayment.GetExpiresAt()),
		PaymentMethod: toDomainPaymentMethod(grpcPayment.GetPaymentMethod()),
		Recipient:     toDomainRecipient(grpcPayment.GetRecipient()),
		Refundable:    grpcPayment.GetRefundable(),
		Test:          grpcPayment.GetTest(),
		IncomeAmount:  toDomainAmount(grpcPayment.GetIncomeAmount()),
	}
}

func toDomainPaymentStatus(status pb.PaymentStatus) models.PaymentStatus {
	switch status {
	case pb.PaymentStatus_SUCCESS:
		return models.SuccessPaymentStatus
	case pb.PaymentStatus_FAILED:
		return models.FailedPaymentStatus
	case pb.PaymentStatus_IN_PROCESSING:
		return models.InProcessingPaymentStatus
	case pb.PaymentStatus_ERROR:
		return models.ErrorPaymentStatus
	case pb.PaymentStatus_REFUNDED:
		return models.RefundedPaymentStatus
	case pb.PaymentStatus_CANCELLED:
		return models.CancelledPaymentStatus
	default:
		return ""
	}
}

func toDomainAmount(grpcAmount *pb.Amount) models.Amount {
	if grpcAmount == nil {
		return models.Amount{}
	}

	return models.Amount{
		Value:    grpcAmount.GetValue(),
		Currency: toDomainCurrency(grpcAmount.GetCurrency()),
	}
}

func toDomainCurrency(currency pb.Currency) models.Currency {
	switch currency {
	case pb.Currency_CURRENCY_RUB:
		return models.RussianRubleCurrency
	default:
		return ""
	}
}

func toDomainTime(t *timestamppb.Timestamp) time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.AsTime()
}

func toDomainPaymentMethod(method *pb.PaymentMethod) models.PaymentMethod {
	if method == nil {
		return models.PaymentMethod{}
	}

	return models.PaymentMethod{
		Type: toDomainPaymentType(method.GetType()),
		ID:   method.GetId(),
		Card: toDomainCard(method.GetCard()),
	}
}

func toDomainPaymentType(paymentType pb.PaymentType) string {
	switch paymentType {
	case pb.PaymentType_PAYMENT_TYPE_SBP:
		return "SBP"
	case pb.PaymentType_PAYMENT_TYPE_BANK_CARD:
		return "bank_card"
	default:
		return ""
	}
}

func toDomainCard(grpcCard *pb.Card) models.Card {
	if grpcCard == nil {
		return models.Card{}
	}

	return models.Card{
		Number:      grpcCard.GetNumber(),
		ExpiryMonth: int(grpcCard.GetExpiryMonth()),
		ExpiryYear:  int(grpcCard.GetExpiryYear()),
		CardType:    grpcCard.GetCardType(),
		CardProduct: struct {
			Code int    `json:"code"`
			Name string `json:"name"`
		}{
			Code: int(grpcCard.GetCardProduct().GetCode()),
			Name: grpcCard.GetCardProduct().GetName(),
		},
		IssuerCountry: grpcCard.GetIssuerCountry(),
		IssuerName:    grpcCard.GetIssuerName(),
	}
}

func toDomainRecipient(grpcRecipient *pb.Recipient) models.Recipient {
	if grpcRecipient == nil {
		return models.Recipient{}
	}

	return models.Recipient{
		AccountNumber: grpcRecipient.GetAccountNumber(),
		Title:         grpcRecipient.GetTitle(),
	}
}
