package transaction

import (
	"alc/model/checkout"
	"alc/model/currency"
	"time"
)

type TransactionStatus string

const (
	Pending    TransactionStatus = "PENDING"
	Authorised TransactionStatus = "AUTHORISED"
	Completed  TransactionStatus = "COMPLETED"
	Failed     TransactionStatus = "FAILED"
	Cancelled  TransactionStatus = "CANCELLED"
)

type Transaction struct {
	Id        int
	TransId   string
	TransUuid *string
	Order     checkout.Order
	Status    TransactionStatus
	Amount    int
	Currency  currency.Currency
	Platform  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
