package checkout

import (
	"alc/model/store"
	"time"
)

type Order struct {
	PurchaseOrder int
	Email         string
	Phone         string
	Name          string
	Address       string
	City          string
	PostalCode    string
	CreatedAt     time.Time
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
