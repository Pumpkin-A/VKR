package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"transaction_service/config"
	consumer "transaction_service/internal/broker/consumer"
	producer "transaction_service/internal/broker/producer"
	"transaction_service/internal/db"
	"transaction_service/internal/entity/paymentManager"

	"golang.org/x/sync/errgroup"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.New()

	db := db.New(cfg)
	defer func() {
		if err := db.DB.Close(); err != nil {
			slog.Error("failed to close DB", "err", err)
		} else {
			slog.Info("DB was closed")
		}
	}()

	producer := producer.NewProducer(cfg)
	defer func() {
		if err := producer.Close(); err != nil {
			slog.Error("error with closing producer")
		} else {
			slog.Info("producer was closed")
		}
	}()

	pm := paymentManager.New(db, producer)

	consumer, _ := consumer.NewConsumer(mainCtx, cfg, pm)

	g, _ := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		defer slog.Info("GatewayConsumer was closed")
		defer stop()

		consumer.ReadFromGateway(mainCtx)
		return nil
	})

	g.Go(func() error {
		defer slog.Info("BillingConsumer was closed")
		defer stop()

		consumer.ReadFromBilling(mainCtx)
		return nil
	})

	if err := g.Wait(); err != nil {
		slog.Info("service exit reason", "err", err.Error())
	}
	slog.Info("servcice exiting")

}
