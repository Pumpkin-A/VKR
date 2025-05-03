package models

import (
	"time"
)

type (
	PaymentStatus string
	Currency      string
	PaymentType   string
)

var (
	CreatedPaymentStatus      PaymentStatus = "created"
	SuccessPaymentStatus      PaymentStatus = "success"
	FailedPaymentStatus       PaymentStatus = "failed"
	InProcessingPaymentStatus PaymentStatus = "inProcessing"
	RussianRubleCurrency      Currency      = "RUB"
	SBPPaymentType            PaymentType   = "SBP"
	CreditCardPaymentType     PaymentType   = "bank_card"
)

type Payment struct {
	UUID   string        `json:"id"`
	Status PaymentStatus `json:"status"`
	Paid   bool          `json:"paid"`
	Amount amount        `json:"amount"`
	// AuthorizationDetails authorizationDetails `json:"authorization_details"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	ExpiresAt   time.Time `json:"expires_at"`
	// Metadata    struct {
	// } `json:"metadata"`
	PaymentMethod paymentMethod `json:"payment_method"`
	Recipient     recipient     `json:"recipient"`
	Refundable    bool          `json:"refundable"`
	Test          bool          `json:"test"`
	IncomeAmount  amount        `json:"income_amount"`
}

type amount struct {
	Value    string   `json:"value"`
	Currency Currency `json:"currency"`
}

type authorizationDetails struct {
	// Rrn      string `json:"rrn"`
	AuthCode string `json:"auth_code"`
	// ThreeDSecure struct {
	// 	Applied bool `json:"applied"`
	// } `json:"three_d_secure"`
}

type paymentMethod struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	// Saved bool   `json:"saved"`
	Card Card `json:"card"`
	// Title string `json:"title"`
}

type Card struct {
	// First6      string `json:"first6"`
	// Last4       string `json:"last4"`
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
