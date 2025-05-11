package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"statistics/config"
	consumer "statistics/internal/broker/consumer"
	client "statistics/internal/client/bankExampleClient"
	statisticsManager "statistics/internal/entity/statisticsManager"
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

	statisticsManager := statisticsManager.New(client)

	consumer, _ := consumer.NewConsumer(mainCtx, cfg, statisticsManager)

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
