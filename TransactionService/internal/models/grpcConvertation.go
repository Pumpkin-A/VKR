package models

import (
	pb "transaction_service/pkg/pb/github.com/yourproject/pkg/pb/transaction/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertPaymentToGrpc(payment Payment) *pb.PaymentResponse {
	// Конвертация времени
	createdAt := timestamppb.New(payment.CreatedAt)
	expiresAt := timestamppb.New(payment.ExpiresAt)

	// Конвертация статуса (пример, адаптируйте под ваши значения)
	var status pb.PaymentStatus
	switch payment.Status {
	case InProcessingPaymentStatus:
		status = pb.PaymentStatus_IN_PROCESSING
	case SuccessPaymentStatus:
		status = pb.PaymentStatus_SUCCESS
	case FailedPaymentStatus:
		status = pb.PaymentStatus_FAILED
	case CancelledPaymentStatus:
		status = pb.PaymentStatus_CANCELLED
	case RefundedPaymentStatus:
		status = pb.PaymentStatus_REFUNDED
	default:
		status = pb.PaymentStatus_ERROR
	}

	return &pb.PaymentResponse{
		Uuid:          payment.UUID,
		Status:        status,
		Paid:          payment.Paid,
		Amount:        convertAmountToGrpc(payment.Amount),
		CreatedAt:     createdAt,
		Description:   payment.Description,
		ExpiresAt:     expiresAt,
		PaymentMethod: convertPaymentMethodToGrpc(payment.PaymentMethod),
		Recipient:     convertRecipientToGrpc(payment.Recipient),
		Refundable:    payment.Refundable,
		Test:          payment.Test,
		IncomeAmount:  convertAmountToGrpc(payment.IncomeAmount),
	}
}

func convertAmountToGrpc(a amount) *pb.Amount {
	var currency pb.Currency
	switch a.Currency {
	case "RUB":
		currency = pb.Currency_CURRENCY_RUB
	// ... другие валюты
	default:
		currency = pb.Currency_CURRENCY_UNSPECIFIED
	}

	return &pb.Amount{
		Value:    a.Value,
		Currency: currency,
	}
}

func convertPaymentMethodToGrpc(pm paymentMethod) *pb.PaymentMethod {
	var paymentType pb.PaymentType
	switch pm.Type {
	case "SBP":
		paymentType = pb.PaymentType_PAYMENT_TYPE_SBP
	case "bank_card":
		paymentType = pb.PaymentType_PAYMENT_TYPE_BANK_CARD
	default:
		paymentType = pb.PaymentType_PAYMENT_TYPE_UNSPECIFIED
	}

	return &pb.PaymentMethod{
		Type: paymentType,
		Id:   pm.ID,
		Card: convertCardToGrpc(pm.Card),
	}
}

func convertCardToGrpc(c Card) *pb.Card {
	return &pb.Card{
		Number:        c.Number,
		ExpiryMonth:   int32(c.ExpiryMonth),
		ExpiryYear:    int32(c.ExpiryYear),
		CardType:      c.CardType,
		CardProduct:   convertCardProductToGrpc(c.CardProduct),
		IssuerCountry: c.IssuerCountry,
		IssuerName:    c.IssuerName,
	}
}

func convertCardProductToGrpc(cp struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}) *pb.CardProduct {
	return &pb.CardProduct{
		Code: int32(cp.Code),
		Name: cp.Name,
	}
}

func convertRecipientToGrpc(r recipient) *pb.Recipient {
	return &pb.Recipient{
		AccountNumber: r.AccountNumber,
		Title:         r.Title,
	}
}
