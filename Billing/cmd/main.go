package main

import (
	"billing/config"
	consumer "billing/internal/broker/consumer"
	producer "billing/internal/broker/producer"
	client "billing/internal/client/bankExampleClient"
	"billing/internal/entity/paymentProcessor"
	"billing/internal/tracer"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	jaegerEndpoint := "localhost:4317"
	tp, err := tracer.InitTracer("Billing", jaegerEndpoint)
	if err != nil {
		slog.Error("failed to initialize tracer", "err", err.Error())
		stop()
	}
	defer tracer.ShutdownTracer(context.Background(), tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	cfg := config.New()

	producer := producer.NewProducer(cfg)
	defer func() {
		if err := producer.Close(); err != nil {
			slog.Error("error with closing producer")
		} else {
			slog.Info("producer was closed")
		}
	}()

	client := client.New(cfg)

	paymentProcessor := paymentProcessor.New(client, producer)

	consumer, _ := consumer.NewConsumer(mainCtx, cfg, paymentProcessor)

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
