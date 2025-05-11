package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"payment_gateway/config"
	"payment_gateway/internal/api"
	"payment_gateway/internal/broker"
	paymentmanager "payment_gateway/internal/entity/paymentManager"
	"payment_gateway/internal/grpcClient"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.New()

	producer := broker.New(cfg)
	defer func() {
		if err := producer.Close(); err != nil {
			slog.Error("error with closing producer")
		} else {
			slog.Info("producer was closed")
		}
	}()

	grpcClient, err := grpcClient.NewPaymentClient(cfg)
	if err != nil {
		slog.Error("error with grpc client creating")
		stop()
	}
	defer func() {
		if err := grpcClient.Close(); err != nil {
			slog.Error("error with closing grpc client")
		}
	}()

	pm := paymentmanager.New(producer, grpcClient)

	s, err := api.New(cfg, pm)
	if err != nil {
		slog.Error("error with creation api", "err", err.Error())
	}

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		defer slog.Info("server was closed")
		defer stop()

		err = s.RunHTTPServer()
		if err != nil {
			slog.Error("error with run http server", "err", err.Error())
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		return s.Srv.Close()
	})

	if err := g.Wait(); err != nil {
		slog.Info("service exit reason", "err", err.Error())
	}
}
