package main

import (
	"context"
	"transaction_service/config"
	broker "transaction_service/internal/broker/consumer"
	"transaction_service/internal/db"
	"transaction_service/internal/paymentManager"
)

func main() {
	ctx := context.Background()
	cfg := config.New()

	db := db.New(cfg)
	defer db.DB.Close()

	pm := paymentManager.New(db)

	consumer := broker.New(ctx, cfg.Kafka.Topic, pm)
	defer consumer.Close()

	for {
	}

}
