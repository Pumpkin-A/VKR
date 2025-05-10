package models

func ConvertExternalTransactionOperationEventToPayment(event ExternalTransactionOperationEvent) Payment {
	return Payment{
		UUID:          event.UUID,
		Status:        event.Status,
		Paid:          event.Paid,
		Amount:        event.Amount,
		CreatedAt:     event.CreatedAt,
		Description:   event.Description,
		ExpiresAt:     event.ExpiresAt,
		PaymentMethod: event.PaymentMethod,
		Recipient:     event.Recipient,
		Refundable:    event.Refundable,
		Test:          event.Test,
		IncomeAmount:  event.IncomeAmount,
	}
}

func ConvertPaymentToInternalTrasactionOperationEvent(p Payment) InternalTransactionOperationEvent {
	return InternalTransactionOperationEvent{
		UUID:          p.UUID,
		Amount:        p.Amount,
		CreatedAt:     p.CreatedAt,
		ExpiresAt:     p.ExpiresAt,
		PaymentMethod: p.PaymentMethod,
		Recipient:     p.Recipient,
		Refundable:    p.Refundable,
		Test:          p.Test,
		IncomeAmount:  p.IncomeAmount,
	}
}
