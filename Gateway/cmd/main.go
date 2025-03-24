package main

import (
	"log"
	"payment_gateway/config"
	"payment_gateway/internal/api"
	paymentmanager "payment_gateway/internal/businessLogic"
	"payment_gateway/internal/db"
)

func main() {
	// ctx := context.Background()
	cfg := config.New()

	db := db.New("db connection string")

	pm := paymentmanager.New(db)

	s, err := api.New(cfg, pm)
	if err != nil {
		log.Panic("server error")
	}

	err = s.RunHTTPServer()
	if err != nil {
		log.Panic("server error")
	}
}
