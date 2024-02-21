package model

type ServiceItem struct {
	Name        string
	Description string
}

type GarantiaItem struct {
	Name  string
	Price int // Stored in PEN centimos
	Slug  string
	Img   string
}
