package service

import (
	"alc/model/checkout"
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Transaction struct {
	db *pgxpool.Pool
}

func NewTransactionService(db *pgxpool.Pool) Transaction {
	return Transaction{
		db: db,
	}
}

func (ts Transaction) InsertTransaction(order checkout.Order, amount int, platform string) (checkout.Transaction, error) {
	var trans checkout.Transaction
	sql := `INSERT INTO store_transactions
	(order_id, amount, platform)
	VALUES
	($1, $2, $3)
	RETURNING
	id, status, amount, platform, created_at, updated_at, trans_id`
	if err := ts.db.QueryRow(context.Background(), sql, order.Id, amount, platform).Scan(&trans.Id, &trans.Status,
		&trans.Amount, &trans.Platform, &trans.CreatedAt, &trans.UpdatedAt, &trans.TransId); err != nil {
		return checkout.Transaction{}, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	trans.Order = order
	return trans, nil
}
