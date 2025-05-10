package models

func (event *EventExternalTransactionOperation) ConvertToPayment() Payment {
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

func (p *Payment) ConvertToInternalTrasactionOperationEvent(operation TransactionOperation) EventInternalTransactionOperation {
	return EventInternalTransactionOperation{
		UUID:                 p.UUID,
		TransactionOperation: operation,
		Amount:               p.Amount,
		CreatedAt:            p.CreatedAt,
		ExpiresAt:            p.ExpiresAt,
		PaymentMethod:        p.PaymentMethod,
		Recipient:            p.Recipient,
		Refundable:           p.Refundable,
		Test:                 p.Test,
		IncomeAmount:         p.IncomeAmount,
	}
}

func (event *EventInternalPaymentResult) ConvertToPaymentResult() PaymentResult {
	return PaymentResult{
		UUID:                 event.UUID,
		TransactionOperation: event.TransactionOperation,
		Status:               event.Status,
		Error:                event.Error,
	}
}

func (res *PaymentResult) ConvertToEventExternalPaymentResult(paymentStatus PaymentStatus) EventExternalPaymentResult {
	return EventExternalPaymentResult{
		UUID:   res.UUID,
		Error:  res.Error,
		Status: paymentStatus,
	}
}
