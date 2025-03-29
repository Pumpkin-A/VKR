package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	models "payment_gateway/internal/models"

	"github.com/google/uuid"
)

func HandleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func (s *Server) HandleCreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var bodyBytes []byte
	var err error

	if r.Body != nil {
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
			return
		}
		defer r.Body.Close()
	}

	// fmt.Printf("Headers: %+v\n", r.Header)

	var paymentReq models.CreatePaymentRequest
	if len(bodyBytes) > 0 {
		if err = json.Unmarshal(bodyBytes, &paymentReq); err != nil {
			fmt.Printf("JSON parse error: %v", err)
			return
		}
		fmt.Println(paymentReq)
	} else {
		fmt.Printf("Body: No Body Supplied\n")
	}

	paymentUUID, err := s.PaymentManager.CreatePayment(context.Background(), paymentReq)
	if err != nil {
		log.Println("error with creating payment")
	}

	w.Write([]byte(paymentUUID + "success"))
}

func (s *Server) HandleGetPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	paymentUUID := r.URL.Query().Get("paymentUUID")
	if paymentUUID == "" || uuid.Validate(paymentUUID) != nil {
		http.Error(w, "Invalid query parameter", http.StatusBadRequest)
		return
	}

	payment, err := s.PaymentManager.GetPayment(context.Background(), paymentUUID)
	if err != nil {
		http.Error(w, "Error in business logic", http.StatusBadRequest)
		return
	}

	jsonPayment, err := json.Marshal(payment)
	if err != nil {
		slog.Error("Error marshaling JSON:", err.Error())
		http.Error(w, "Error in business logic", http.StatusBadRequest)
		return
	}

	slog.Info("payment:", string(jsonPayment))
	w.Write(jsonPayment)
}
