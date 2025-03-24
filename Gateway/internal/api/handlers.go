package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	models "payment_gateway/internal/models"
)

func HandleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func (s *Server) HandleCreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("error"))
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

	paymentID, err := s.PaymentManager.CreatePayment(context.Background(), paymentReq)
	if err != nil {
		log.Println("error with creating payment")
	}

	w.Write([]byte(paymentID + "success"))
}
