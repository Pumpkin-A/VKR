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
	GetPaymentInfo(ctx context.Context, id string) (models.Payment, error)
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
}

func (s *Server) RunHTTPServer() error {
	slog.Info("starting http listener", "address", fmt.Sprintf("http://%s", s.Srv.Addr))
	return s.Srv.ListenAndServe()
}
