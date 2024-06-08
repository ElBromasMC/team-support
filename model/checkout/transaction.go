package checkout

import "time"

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
	Order     Order
	Status    TransactionStatus
	Amount    int
	Platform  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
