package model

type StoreItem struct {
	Category         string
	Name             string
	Price            int // Stored in PEN centimos
	BriefDescription string
	Description      string
	Slug             string
	Img              string
	LargeImg         string
	Quantity         int
}
