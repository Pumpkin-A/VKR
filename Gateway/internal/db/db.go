package db

import models "payment_gateway/internal/models"

type DB struct {
	connString string
}

func New(connString string) *DB {
	return &DB{
		connString: connString,
	}
}

func (db *DB) AddPayment(models.Payment) error {
	return nil
}

func (db *DB) GetPayment(id int) models.Payment {
	return models.Payment{}
}
