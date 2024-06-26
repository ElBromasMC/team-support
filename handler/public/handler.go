package public

import (
	"alc/service"
)

type Handler struct {
	PublicService      service.Public
	EmailService       service.Email
	AuthService        service.Auth
	OrderService       service.Order
	TransactionService service.Transaction
	PaymentService     service.Payment
	DeviceService      service.Device
}
