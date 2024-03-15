package store

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Type string

const (
	StoreType    Type = "STORE"
	GarantiaType Type = "GARANTIA"
)

type Category struct {
	Uuid        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Type        Type      `json:"type"`
	Description string    `json:"description"`
	Img         string    `json:"img"`
	Slug        string    `json:"slug"`
}

type Item struct {
	Uuid            uuid.UUID `json:"uuid"`
	Name            string    `json:"name"`
	Category        Category  `json:"category"`
	Description     string    `json:"description"`
	LongDescription string    `json:"-"`
	Slug            string    `json:"slug"`
	Img             string    `json:"img"`
	LargeImg        string    `json:"-"`
}

type Product struct {
	Uuid  uuid.UUID `json:"uuid"`
	Name  string    `json:"name"`
	Item  Item      `json:"item"`
	Price int       `json:"price"` // Stored in USD cents
}

type Order struct {
	Uuid    uuid.UUID `json:"uuid"`
	Email   string    `json:"email"`
	Address string    `json:"address"`
	// Status    string    `json:"status"` // "Pending", "Shipped", "Delivered"
	CreatedAt time.Time `json:"createdAt"`
}

type OrderProduct struct {
	ProductID uuid.UUID         `json:"productID"`
	Quantity  int               `json:"quantity"`
	Details   map[string]string `json:"details"`
}

type PurchaseRequest struct {
	Order         Order          `json:"order"`
	OrderProducts []OrderProduct `json:"orderProducts"`
}

func (i Product) ToJSON() string {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
