package db

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"transaction_service/config"
	models "transaction_service/internal/models"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func New(cfg config.Config) *PostgresDB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.DbHost, cfg.DB.DbPort, cfg.DB.DbUser, cfg.DB.DbPassword, cfg.DB.DbName, cfg.DB.SSLmode)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		slog.Error("error with opening db", "err", err.Error())
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		slog.Error("error with pinging db", "err", err.Error())
		os.Exit(1)
	}

	slog.Info("db connection success", "user", cfg.DB.DbUser, "dbname", cfg.DB.DbName)
	return &PostgresDB{
		DB: db,
	}
}

func (db *PostgresDB) AddPayment(p models.Payment) error {
	insertQuery := `insert into payment (uuid, status_id, paid, amount, currency_id, created_at, expired_at, description, 
	paymnent_type_id, card_number, recepient_account_number, refundable, test, income)
	values ($1, (SELECT id from status WHERE status.status=$2), $3, $4, (SELECT id from currency WHERE title=$5), 
	$6, $7,  $8, (SELECT id from payment_type WHERE payment_type.type=$9), $10, $11, $12, $13, $14)`
	_, err := db.DB.Exec(insertQuery, p.UUID, p.Status, p.Paid, p.Amount.Value, p.Amount.Currency, p.CreatedAt, p.ExpiresAt,
		p.Description, p.PaymentMethod.Type, p.PaymentMethod.Card.Number, p.Recipient.AccountNumber, p.Refundable, p.Test, p.IncomeAmount.Value)
	if err != nil {
		slog.Error("error with adding payment", "err", err.Error())
		return err
	}
	slog.Info("payment was successfully added to DB", "uuid:", p.UUID)
	return nil
}

func (db *PostgresDB) UpdatePaymentStatus(uuid, status string) error {
	updateQuery := `UPDATE payment
	SET status_id = (SELECT id from status WHERE status.status=$1)
	WHERE uuid = $2;`
	_, err := db.DB.Exec(updateQuery, status, uuid)
	if err != nil {
		slog.Error("error with updating payment status", "err", err.Error())
		return err
	}
	slog.Info("payment status was successfully updated in DB", "uuid:", uuid)
	return nil
}

func (db *PostgresDB) AddCardIfNotExist(c models.Card) error {
	insertQuery := `insert into card (number, expiry_month, expiry_year, card_type, code, name, issuer_country, issuer_name_id)
	select $1, $2, $3, $4, $5, $6, $7, (select id from issuer_name where issuer_name.issuer_name=$8)
	where NOT EXISTS (SELECT number FROM card WHERE number = $9);`
	_, err := db.DB.Exec(insertQuery, c.Number, c.ExpiryMonth, c.ExpiryYear, c.CardType, c.CardProduct.Code, c.CardProduct.Name,
		c.IssuerCountry, c.IssuerName, c.Number)
	if err != nil {
		slog.Error("error with adding card", "err", err.Error())
		return err
	}
	slog.Info("card was successfully added to DB", "card_number:", c.Number)
	return nil
}

func (db *PostgresDB) GetPayment(uuid string) (models.Payment, error) {
	selectQuery := `select p.uuid, s.status, p.paid, p.amount, cu.title, p.created_at, p.expired_at,
       p.description, pt.type, c.number, c.expiry_month, c.expiry_year, c.card_type,
       c.code, c.name, c.issuer_country, i.issuer_name,
       r.title, p.refundable, p.test, p.income from payment as p
		join card as c on p.card_number = c.number
		join payment_type as pt on p.paymnent_type_id = pt.id
		join status as s on p.status_id = s.id
		join currency as cu on p.currency_id = cu.id
		join recepient as r on p.recepient_account_number = r.account_number
		join issuer_name AS i ON c.issuer_name_id = i.id
		where uuid = $1;`
	row := db.DB.QueryRow(selectQuery, uuid)
	if err := row.Err(); err != nil {
		slog.Error("error with query row in GetPayment", "err", err.Error())
		return models.Payment{}, err
	}
	log.Println(row)

	p := models.Payment{}
	err := row.Scan(&p.UUID, &p.Status, &p.Paid, &p.Amount.Value, &p.Amount.Currency, &p.CreatedAt, &p.ExpiresAt,
		&p.Description, &p.PaymentMethod.Type, &p.PaymentMethod.Card.Number, &p.PaymentMethod.Card.ExpiryMonth,
		&p.PaymentMethod.Card.ExpiryYear, &p.PaymentMethod.Card.CardType, &p.PaymentMethod.Card.CardProduct.Code,
		&p.PaymentMethod.Card.CardProduct.Name, &p.PaymentMethod.Card.IssuerCountry, &p.PaymentMethod.Card.IssuerName,
		&p.Recipient.Title, &p.Refundable, &p.Test, &p.IncomeAmount.Value)
	if err != nil {
		slog.Error("error with scan row in GetPayment", "err", err.Error())
		return models.Payment{}, err
	}
	log.Println(p)

	return p, nil
}

func (db *PostgresDB) GetPaymentStatus(uuid string) (models.PaymentStatus, error) {
	selectQuery := `select s.status from payment as p
	join status as s on p.status_id = s.id
	where uuid = $1;`
	row := db.DB.QueryRow(selectQuery, uuid)
	if err := row.Err(); err != nil {
		slog.Error("error with query row in GetPaymentStatus", "err", err.Error())
		return "", err
	}
	log.Println(row)

	var status models.PaymentStatus
	err := row.Scan(&status)
	if err != nil {
		slog.Error("error with scan row in GetPaymentStatus", "err", err.Error())
		return "", err
	}
	log.Println(status)

	return status, nil
}
