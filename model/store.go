package model

import (
	"encoding/json"

	"github.com/gofrs/uuid/v5"
)

type StoreSubCategory struct {
	Name        string
	Description string
	Img         string
	Slug        string
}

type StoreItem struct {
	Uuid             uuid.UUID `json:"uuid"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	SubCategory      string    `json:"subcategory"`
	Price            int       `json:"price"` // Stored in USD cents
	BriefDescription string    `json:"briefDescription"`
	Description      string    `json:"-"`
	Slug             string    `json:"slug"`
	Img              string    `json:"img"`
	LargeImg         string    `json:"-"`
	Quantity         int       `json:"quantity,omitempty"`
}

func (i StoreItem) ToJSON() string {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
