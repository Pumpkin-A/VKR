package models

import (
	"time"
)

type (
	PaymentStatus        string
	Currency             string
	PaymentType          string
	TransactionOperation string
)

var (
	SuccessPaymentStatus       PaymentStatus        = "success"
	FailedPaymentStatus        PaymentStatus        = "failed"
	InProcessingPaymentStatus  PaymentStatus        = "inProcessing"
	ErrorPaymentStatus         PaymentStatus        = "error"
	RefundedPaymentStatus      PaymentStatus        = "refunded"
	CancelledPaymentStatus     PaymentStatus        = "cancelled"
	RussianRubleCurrency       Currency             = "RUB"
	SBPPaymentType             PaymentType          = "SBP"
	CreditCardPaymentType      PaymentType          = "bank_card"
	CreateTransactionOperation TransactionOperation = "create"
	RefundTransactionOperation TransactionOperation = "refund"
	CancelTransactionOperation TransactionOperation = "cancel"
)

type CreatePaymentRequest struct {
	Amount        Amount        `json:"amount"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Recipient     Recipient     `json:"recipient"`
}

type MakeRefundRequest struct {
	UUID          string        `json:"id"`
	Amount        Amount        `json:"amount"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Recipient     Recipient     `json:"recipient"`
}

type CancelPayment struct {
	UUID          string        `json:"id"`
	Amount        Amount        `json:"amount"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Recipient     Recipient     `json:"recipient"`
}

type Payment struct {
	UUID          string        `json:"id"`
	Status        PaymentStatus `json:"status"`
	Paid          bool          `json:"paid"`
	Amount        Amount        `json:"amount"`
	CreatedAt     time.Time     `json:"created_at"`
	Description   string        `json:"description"`
	ExpiresAt     time.Time     `json:"expires_at"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Recipient     Recipient     `json:"recipient"`
	Refundable    bool          `json:"refundable"`
	Test          bool          `json:"test"`
	IncomeAmount  Amount        `json:"income_amount"`
}

type ExternalTransactionOperationEvent struct {
	UUID                 string               `json:"id"`
	TransactionOperation TransactionOperation `json:"transactionOperation"`
	Status               PaymentStatus        `json:"status"`
	Paid                 bool                 `json:"paid"`
	Amount               Amount               `json:"amount"`
	CreatedAt            time.Time            `json:"created_at"`
	Description          string               `json:"description"`
	ExpiresAt            time.Time            `json:"expires_at"`
	PaymentMethod        PaymentMethod        `json:"payment_method"`
	Recipient            Recipient            `json:"recipient"`
	Refundable           bool                 `json:"refundable"`
	Test                 bool                 `json:"test"`
	IncomeAmount         Amount               `json:"income_amount"`
}

type Amount struct {
	Value    string   `json:"value"`
	Currency Currency `json:"currency"`
}

type PaymentMethod struct {
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

type Recipient struct {
	AccountNumber string `json:"account_number"`
	Title         string `json:"title"`
}
