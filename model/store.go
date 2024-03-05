package model

import (
	"encoding/json"

	"github.com/gofrs/uuid/v5"
)

type StoreItem struct {
	Uuid             uuid.UUID `json:"uuid"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	Price            int       `json:"price"` // Stored in PEN centimos
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
