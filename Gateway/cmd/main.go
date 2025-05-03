package main

import (
	"log"
	"payment_gateway/config"
	"payment_gateway/internal/api"
	"payment_gateway/internal/broker"
	paymentmanager "payment_gateway/internal/businessLogic"
)

func main() {
	// ctx := context.Background()
	cfg := config.New()

	producer := broker.New(cfg.Kafka.BootstrapServers, cfg.Kafka.Topic)
	defer producer.Close()

	pm := paymentmanager.New(producer)

	s, err := api.New(cfg, pm)
	if err != nil {
		log.Panic("server error")
	}

	err = s.RunHTTPServer()
	if err != nil {
		log.Panic("server error")
	}
}
