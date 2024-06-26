package service

import (
	"alc/model/checkout"
	"alc/model/transaction"
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"
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

func (ts Transaction) GetTransaction(order checkout.Order) (transaction.Transaction, error) {
	var trans transaction.Transaction
	sql := `SELECT id, status, amount,
	platform, created_at, updated_at, trans_id, trans_uuid
	FROM store_transactions
	WHERE order_id = $1`
	if err := ts.db.QueryRow(context.Background(), sql, order.Id).Scan(&trans.Id, &trans.Status, &trans.Amount,
		&trans.Platform, &trans.CreatedAt, &trans.UpdatedAt, &trans.TransId, &trans.TransUuid); err != nil {
		return transaction.Transaction{}, echo.NewHTTPError(http.StatusNotFound, "Transacción no encontrada")
	}
	trans.Order = order
	return trans, nil
}

func (ts Transaction) InsertTransaction(order checkout.Order, amount int, platform string) (transaction.Transaction, error) {
	var trans transaction.Transaction
	sql := `INSERT INTO store_transactions
	(order_id, amount, platform)
	VALUES
	($1, $2, $3)
	RETURNING
	id, status, amount, platform, created_at, updated_at, trans_id`
	if err := ts.db.QueryRow(context.Background(), sql, order.Id, amount, platform).Scan(&trans.Id, &trans.Status,
		&trans.Amount, &trans.Platform, &trans.CreatedAt, &trans.UpdatedAt, &trans.TransId); err != nil {
		return transaction.Transaction{}, echo.NewHTTPError(http.StatusInternalServerError, "Error desconocido, recargue la página")
	}
	trans.Order = order
	return trans, nil
}

func (ts Transaction) UpdateTransaction(orderId uuid.UUID, transId string, transUuid string, status transaction.TransactionStatus) error {
	sql := `UPDATE store_transactions
	SET trans_uuid = $1, status = $2
	WHERE order_id = $3 AND trans_id = $4`
	c, err := ts.db.Exec(context.Background(), sql, transUuid, status, orderId, transId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if c.RowsAffected() != 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Transaction not found")
	}
	return nil
}
