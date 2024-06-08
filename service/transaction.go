package service

import "github.com/jackc/pgx/v5/pgxpool"

type Transaction struct {
	db *pgxpool.Pool
}

func NewTransactionService(db *pgxpool.Pool) Transaction {
	return Transaction{
		db: db,
	}
}
