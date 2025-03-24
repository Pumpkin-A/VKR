package models

import (
	"time"

	"github.com/google/uuid"
)

func ConvertCreatePaymentRequestToPayment(req CreatePaymentRequest) Payment {
	return Payment{
		ID:     uuid.NewString(),
		Status: CreatedPaymentStatus,
		Paid:   false,
		Amount: req.Amount,
		AuthorizationDetails: authorizationDetails{
			Rrn:      "10000000000",
			AuthCode: "000000",
		},
		CreatedAt:     time.Now(),
		Description:   "",
		ExpiresAt:     time.Now().Add(time.Duration(time.Minute * 30)),
		PaymentMethod: req.PaymentMethod,
		Recipient:     req.Recipient,
		Refundable:    false,
		Test:          false,
		IncomeAmount: amount{
			Value:    "33.33",
			Currency: RussianRubleCurrency,
		},
	}
}
