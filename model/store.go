package model

type StoreItem struct {
	Category string
	Name     string
	Price    int // Stored in PEN centimos
	Slug     string
	Img      string
}
