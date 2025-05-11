package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"payment_gateway/config"
	"payment_gateway/internal/models"
)

type Server struct {
	PaymentManager PaymentManager
	Mux            *http.ServeMux
	Srv            *http.Server
}

type PaymentManager interface {
	CreatePayment(ctx context.Context, requestData models.CreatePaymentRequest) (string, error)
	GetPayment(ctx context.Context, uuid string) (models.Payment, error)
	CancelPayment(ctx context.Context, requestData models.CancelPayment) (string, error)
	MakeRefund(ctx context.Context, requestData models.MakeRefundRequest) (string, error)
}

func New(cfg config.Config, pm PaymentManager) (*Server, error) {
	s := &Server{
		Mux:            http.NewServeMux(),
		PaymentManager: pm,
	}
	listenAddress := fmt.Sprintf("%s:%d", "0.0.0.0", cfg.Server.Port)
	s.Srv = &http.Server{
		Addr:    listenAddress,
		Handler: s.Mux,
	}

	s.registerHandlers()
	return s, nil
}

func (s *Server) registerHandlers() {
	s.Mux.HandleFunc("/", HandleHello)
	s.Mux.HandleFunc("/payments", s.HandleCreatePayment)
	s.Mux.HandleFunc("/refundPayment", s.HandleMakeRefund)
	s.Mux.HandleFunc("/cancelPayment", s.HandleCancelPayment)
	s.Mux.HandleFunc("/getPayment", s.HandleGetPayment)
}

func (s *Server) RunHTTPServer() error {
	slog.Info("starting http listener", "address", fmt.Sprintf("http://%s", s.Srv.Addr))
	return s.Srv.ListenAndServe()
}
