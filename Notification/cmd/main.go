package main

import (
	"context"
	"log/slog"
	"notifications/config"
	consumer "notifications/internal/broker/consumer"
	client "notifications/internal/client/bankExampleClient"
	"notifications/internal/entity/notificationsManager"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.New()

	client := client.New(cfg)

	notificationsManager := notificationsManager.New(client)

	consumer, _ := consumer.NewConsumer(mainCtx, cfg, notificationsManager)

	g, _ := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		defer slog.Info("consumer was closed")
		defer stop()

		consumer.Run(mainCtx)
		return nil
	})

	if err := g.Wait(); err != nil {
		slog.Info("service exit reason", "err", err.Error())
	}
	slog.Info("servcice exiting")

}
