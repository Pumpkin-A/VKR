package main

import (
	"context"
	"transaction_service/config"
	"transaction_service/internal/broker"
	"transaction_service/internal/db"
	"transaction_service/internal/paymentManager"
)

// consume

// db := db.New(cfg)
// defer db.DB.Close()

// producer := broker.New(cfg.Kafka.BootstrapServers, cfg.Kafka.Topic)
// defer producer.Close()

func main() {
	ctx := context.Background()
	cfg := config.New()

	db := db.New(cfg)
	defer db.DB.Close()

	pm := paymentManager.New(db)

	consumer := broker.New(cfg.Kafka.Topic, *pm)
	defer consumer.Close()

	consumer.Read(ctx)

}
