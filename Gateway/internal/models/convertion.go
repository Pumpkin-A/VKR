package models

import (
	"time"

	"github.com/google/uuid"
)

func (req *CreatePaymentRequest) ConvertToExternalTransactionOperationEvent() ExternalTransactionOperationEvent {
	return ExternalTransactionOperationEvent{
		UUID:                 uuid.NewString(),
		TransactionOperation: CreateTransactionOperation,
		Status:               InProcessingPaymentStatus,
		Paid:                 false,
		Amount:               req.Amount,
		CreatedAt:            time.Now(),
		Description:          "",
		ExpiresAt:            time.Now().Add(time.Duration(time.Minute * 30)),
		PaymentMethod:        req.PaymentMethod,
		Recipient:            req.Recipient,
		Refundable:           false,
		Test:                 false,
		IncomeAmount: amount{
			Value:    "33.33",
			Currency: RussianRubleCurrency,
		},
	}
}

func (req *MakeRefundRequest) ConvertToExternalTransactionOperationEvent() ExternalTransactionOperationEvent {
	return ExternalTransactionOperationEvent{
		UUID:                 req.UUID,
		TransactionOperation: RefundTransactionOperation,
		Status:               InProcessingPaymentStatus,
		Paid:                 false,
		Amount:               req.Amount,
		CreatedAt:            time.Now(),
		Description:          "",
		ExpiresAt:            time.Now().Add(time.Duration(time.Minute * 30)),
		PaymentMethod:        req.PaymentMethod,
		Recipient:            req.Recipient,
		Refundable:           false,
		Test:                 false,
		IncomeAmount: amount{
			Value:    "33.33",
			Currency: RussianRubleCurrency,
		},
	}
}

func (req *CancelPayment) ConvertToExternalTransactionOperationEvent() ExternalTransactionOperationEvent {
	return ExternalTransactionOperationEvent{
		UUID:                 req.UUID,
		TransactionOperation: CancelTransactionOperation,
		Status:               InProcessingPaymentStatus,
		Paid:                 false,
		Amount:               req.Amount,
		CreatedAt:            time.Now(),
		Description:          "",
		ExpiresAt:            time.Now().Add(time.Duration(time.Minute * 30)),
		PaymentMethod:        req.PaymentMethod,
		Recipient:            req.Recipient,
		Refundable:           false,
		Test:                 false,
		IncomeAmount: amount{
			Value:    "33.33",
			Currency: RussianRubleCurrency,
		},
	}
}
