package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"transaction_service/config"
	broker "transaction_service/internal/broker/consumer"
	"transaction_service/internal/db"
	"transaction_service/internal/paymentManager"
)

func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.New()

	db := db.New(cfg)
	defer func() {
		db.DB.Close()
		slog.Info("DB was closed")
	}()

	pm := paymentManager.New(db)

	broker.New(mainCtx, cfg, pm)

	for {
	}

}
