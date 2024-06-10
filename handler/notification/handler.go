package notification

import (
	"alc/service"
)

type Handler struct {
	TransactionService service.Transaction
	PaymentService     service.Payment
}
