package models

import "fmt"

func (event *EventInternalTransactionOperation) ConvertToPayment() Payment {
	return Payment{
		UUID:          event.UUID,
		Amount:        event.Amount,
		CreatedAt:     event.CreatedAt,
		ExpiresAt:     event.ExpiresAt,
		PaymentMethod: event.PaymentMethod,
		Recipient:     event.Recipient,
		Refundable:    event.Refundable,
		Test:          event.Test,
		IncomeAmount:  event.IncomeAmount,
	}
}

func (res *ResultOfRequestFromBank) ConvertToEventInternalPaymentResult(operation TransactionOperation) EventInternalPaymentResult {
	var status BankExampleStatus
	fmt.Println(res.Status)
	switch res.Status {
	case string(SuccessedBankExampleStatus):
		status = SuccessedBankExampleStatus
	case string(FailedBankExampleStatus):
		status = FailedBankExampleStatus
	case string(ErrorBankExampleStatus):
		status = ErrorBankExampleStatus
	default:
		status = UnknownBankExampleStatus
	}
	return EventInternalPaymentResult{
		UUID:                 res.UUID,
		TransactionOperation: operation,
		Status:               status,
		Error:                res.ErrorText,
	}
}
