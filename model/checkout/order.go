package checkout

import (
	"alc/model/auth"
	"alc/model/currency"
	"alc/model/store"
	"errors"
	"math"
	"net/mail"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

type OrderSyncStatus string

const (
	Pending   OrderSyncStatus = "PENDING"
	Completed OrderSyncStatus = "COMPLETED"
	Failed    OrderSyncStatus = "FAILED"
)

type Order struct {
	Id            uuid.UUID
	PurchaseOrder int
	Email         string
	Phone         string
	Name          string
	Address       string
	City          string
	PostalCode    string
	AssignedTo    auth.User
	CreatedAt     time.Time
	SyncStatus    OrderSyncStatus
	LockedAt      time.Time
}

type OrderStatus string

const (
	Pendiente    OrderStatus = "PENDIENTE"
	EnProceso    OrderStatus = "EN PROCESO"
	PorConfirmar OrderStatus = "POR CONFIRMAR"
	Entregado    OrderStatus = "ENTREGADO"
	Cancelado    OrderStatus = "CANCELADO"
)

type OrderProduct struct {
	Id                int
	Order             Order
	Quantity          int
	Details           map[string]string
	Product           store.Product
	ProductType       store.Type
	ProductCategory   string
	ProductItem       string
	ProductName       string
	ProductPrice      int
	ProductCurrency   currency.Currency
	ProductDetails    map[string]string
	ProductPartNumber string
	Status            OrderStatus
	UpdatedAt         time.Time
}

func (order Order) Normalize() (Order, error) {
	// Trim the strings
	order.Email = strings.ToLower(strings.TrimSpace(order.Email))
	order.Phone = strings.TrimSpace(order.Phone)
	order.Name = strings.TrimSpace(order.Name)
	order.Address = strings.TrimSpace(order.Address)
	order.City = strings.TrimSpace(order.City)
	order.PostalCode = strings.TrimSpace(order.PostalCode)

	// Validate email
	address, err := mail.ParseAddress(order.Email)
	if err != nil {
		return Order{}, errors.New("invalid email")
	}
	order.Email = address.Address

	// Validate the rest
	if len(order.Name) == 0 {
		return Order{}, errors.New("invalid name")
	}
	if len(order.Address) == 0 {
		return Order{}, errors.New("invalid address")
	}
	if len(order.City) == 0 {
		return Order{}, errors.New("invalid city")
	}
	if len(order.PostalCode) == 0 {
		return Order{}, errors.New("invalid postalcode")
	}

	return order, nil
}

// Order product management

func (p OrderProduct) CalculateIndividualPrice(r currency.ExchangeRate) (int, error) {
	rate, ok := r.Get(p.ProductCurrency)
	if !ok {
		return 0, errors.New("invalid currency")
	}
	newPrice := float64(p.ProductPrice) * rate
	return int(math.Round(newPrice)), nil
}

func (p OrderProduct) CalculateSubtotal(r currency.ExchangeRate) (int, error) {
	rate, ok := r.Get(p.ProductCurrency)
	if !ok {
		return 0, errors.New("invalid currency")
	}
	newPrice := float64(p.Quantity) * float64(p.ProductPrice) * rate
	return int(math.Round(newPrice)), nil
}

func CalculateAmount(r currency.ExchangeRate, products []OrderProduct) (int, error) {
	amount := 0
	for _, product := range products {
		subtotal, err := product.CalculateSubtotal(r)
		if err != nil {
			return 0, err
		}
		amount += subtotal
	}
	return amount, nil
}
