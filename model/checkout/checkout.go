package checkout

import (
	"alc/model"
	"alc/model/store"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Status string

const (
	Pendiente    Status = "PENDIENTE"
	Asignado     Status = "ASIGNADO"
	EnProceso    Status = "EN PROCESO"
	PorConfirmar Status = "POR CONFIRMAR"
	Realizado    Status = "REALIZADO"
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
	CreatedAt     time.Time
	AssignedTo    model.User
}

type OrderProduct struct {
	Id              int
	Order           Order
	Quantity        int
	Details         map[string]string
	ProductType     store.Type
	ProductCategory string
	ProductItem     string
	ProductName     string
	ProductPrice    int
	ProductDetails  map[string]string
}
