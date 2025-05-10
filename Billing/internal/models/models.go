package models

import (
	"time"
)

type (
	PaymentStatus        string
	Currency             string
	PaymentType          string
	TransactionOperation string
	BankExampleStatus    string
)

var (
	SuccessPaymentStatus       PaymentStatus        = "success"
	FailedPaymentStatus        PaymentStatus        = "failed"
	InProcessingPaymentStatus  PaymentStatus        = "inProcessing"
	ErrorStatus                PaymentStatus        = "error"
	RefundedStatus             PaymentStatus        = "refunded"
	CancelledStatus            PaymentStatus        = "cancelled"
	RussianRubleCurrency       Currency             = "RUB"
	SBPPaymentType             PaymentType          = "SBP"
	CreditCardPaymentType      PaymentType          = "bank_card"
	CreateTransactionOperation TransactionOperation = "create"
	RefundTransactionOperation TransactionOperation = "refund"
	CancelTransactionOperation TransactionOperation = "cancel"
	SuccessedBankExampleStatus BankExampleStatus    = "success"
	FailedBankExampleStatus    BankExampleStatus    = "failed"
	ErrorBankExampleStatus     BankExampleStatus    = "error"
	UnknownBankExampleStatus   BankExampleStatus    = "unknown"
)

type ResultOfRequestFromBank struct {
	UUID      string `json:"UUID"`
	Status    string `json:"status"` // successed, failed, error
	ErrorText string `json:"errorText"`
}

type Payment struct {
	UUID          string        `json:"id"`
	Amount        amount        `json:"amount"`
	CreatedAt     time.Time     `json:"created_at"`
	ExpiresAt     time.Time     `json:"expires_at"`
	PaymentMethod paymentMethod `json:"payment_method"`
	Recipient     recipient     `json:"recipient"`
	Refundable    bool          `json:"refundable"`
	Test          bool          `json:"test"`
	IncomeAmount  amount        `json:"income_amount"`
}

type EventInternalTransactionOperation struct {
	UUID                 string               `json:"id"`
	TransactionOperation TransactionOperation `json:"transactionOperation"`
	Amount               amount               `json:"amount"`
	CreatedAt            time.Time            `json:"created_at"`
	ExpiresAt            time.Time            `json:"expires_at"`
	PaymentMethod        paymentMethod        `json:"payment_method"`
	Recipient            recipient            `json:"recipient"`
	Refundable           bool                 `json:"refundable"`
	Test                 bool                 `json:"test"`
	IncomeAmount         amount               `json:"income_amount"`
}

type EventInternalPaymentResult struct {
	UUID                 string               `json:"id"`
	TransactionOperation TransactionOperation `json:"transactionOperation"`
	Status               BankExampleStatus    `json:"status"`
	Error                string               `json:"errorText"`
}

type amount struct {
	Value    string   `json:"value"`
	Currency Currency `json:"currency"`
}

type paymentMethod struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Card Card   `json:"card"`
}

type Card struct {
	Number      string `json:"number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CardType    string `json:"card_type"`
	CardProduct struct {
		Code int    `json:"code"`
		Name string `json:"name"`
	} `json:"card_product"`
	IssuerCountry string `json:"issuer_country"`
	IssuerName    string `json:"issuer_name"`
}

type recipient struct {
	AccountNumber string `json:"account_number"`
	Title         string `json:"title"`
}
