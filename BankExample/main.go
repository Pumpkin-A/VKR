package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/dopayment", DoPayment) // Регистрируем обработчик

	fmt.Println("Server listening on port 9999")
	err := http.ListenAndServe(":9999", nil) // Start server
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// DoPayment - обработчик запроса на проведение платежа.
func DoPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payment Payment
	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %s", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Генерируем рандомно статус
	randNum := rand.Float64()

	var status string
	switch {
	case randNum <= 0.6: // 60% chance
		status = "success"
	case randNum <= 0.9: // 30% chance (0.6 + 0.3)
		status = "failed"
	default: // 10% chance
		status = "error"
	}

	// Формируем результат
	result := PaymentResult{
		UUID:   payment.UUID,
		Status: status,
	}

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// Payment - структура данных для запроса платежа.
type Payment struct {
	UUID          string        `json:"id"`
	Amount        Amount        `json:"amount"`
	CreatedAt     time.Time     `json:"created_at"`
	ExpiresAt     time.Time     `json:"expires_at"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Recipient     Recipient     `json:"recipient"`
	Refundable    bool          `json:"refundable"`
	Test          bool          `json:"test"`
	IncomeAmount  Amount        `json:"income_amount"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
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

// PaymentResult - структура данных для ответа с результатом платежа.
type PaymentResult struct {
	UUID      string `json:"UUID"`
	Status    string `json:"status"` // successed, failed, error
	ErrorText string `json:"errorText"`
}
