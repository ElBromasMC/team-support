package store

import (
	"time"
)

type Image struct {
	Id       int    `json:"id"`
	Filename string `json:"filename"`
}

type Type string

const (
	StoreType    Type = "STORE"
	GarantiaType Type = "GARANTIA"
)

type Category struct {
	Id          int    `json:"id"`
	Type        Type   `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Img         Image  `json:"img"`
	Slug        string `json:"slug"`
}

type Item struct {
	Id              int      `json:"id"`
	Category        Category `json:"category"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	LongDescription string   `json:"longDescription"`
	Img             Image    `json:"img"`
	LargeImg        Image    `json:"largeImg"`
	Slug            string   `json:"slug"`
}

type Product struct {
	Id      int               `json:"id"`
	Item    Item              `json:"item"`
	Name    string            `json:"name"`
	Price   int               `json:"price"` // Stored in USD cents
	Details map[string]string `json:"details"`
}

type Order struct {
	PurchaseOrder int       `json:"purchaseOrder"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Address       string    `json:"address"`
	PhoneNumber   string    `json:"phoneNumber"`
	CreatedAt     time.Time `json:"createdAt"`
}
