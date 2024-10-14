package store

import (
	"alc/config"
	"alc/model/auth"
	"alc/model/currency"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
	"unicode"
)

// Store management

type Image struct {
	Id       int    `json:"id"`
	Filename string `json:"filename"`
	Index    int    `json:"index"`
}

type Type string

const (
	StoreType    Type = "STORE"
	GarantiaType Type = "GARANTIA"
)

func (t Type) ToSlug() string {
	if t == GarantiaType {
		return "garantia"
	} else {
		return "store"
	}
}

func (t Type) ToTitle() string {
	if t == GarantiaType {
		return "Garantías"
	} else {
		return "Tienda"
	}
}

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
	VendorLink      string   `json:"vendorLink"`
	Img             Image    `json:"img"`
	LargeImg        Image    `json:"largeImg"`
	Slug            string   `json:"slug"`
}

type Product struct {
	Id                    int               `json:"id"`
	Item                  Item              `json:"item"`
	Name                  string            `json:"name"`
	Price                 int               `json:"price"`
	Currency              currency.Currency `json:"currency"`
	Stock                 *int              `json:"stock"`
	Details               map[string]string `json:"details"`
	PartNumber            string            `json:"partNumber"`
	AcceptBeforeSixMonths bool              `json:"acceptBefore"`
	AcceptAfterSixMonths  bool              `json:"acceptAfter"`
	Slug                  string            `json:"slug"`
}

func (product Product) Normalize() (Product, error) {
	// Trim name
	product.Name = strings.TrimSpace(product.Name)

	// Remove spaces and uppercase part number
	product.PartNumber = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToUpper(r)
	}, product.PartNumber)

	if len(product.Name) == 0 {
		return Product{}, errors.New("invalid name")
	}
	if product.Price < 0 {
		return Product{}, errors.New("invalid price")
	}
	if product.Stock != nil {
		if *product.Stock < 0 {
			return Product{}, errors.New("invalid stock")
		}
	}
	// if len(product.PartNumber) == 0 {
	// 	return Product{}, errors.New("invalid part number")
	// }

	return product, nil
}

type ProductDiscount struct {
	Id            int
	Product       Product
	DiscountValue int
	ValidFrom     time.Time
	ValidUntil    time.Time
	CouponCode    *string
	MinimumAmount *int
	MaximumAmount *int
}

// Comment management

type ItemComment struct {
	Id          int
	Item        Item
	CommentedBy auth.User
	Title       string
	Message     string
	Rating      int
	UpVotes     int
	DownVotes   int
	IsEdited    bool
	CreatedAt   time.Time
	EditedAt    time.Time
}

// Serial management

type Device struct {
	Id                int
	Serie             string
	Valid             bool
	IsBeforeSixMonths bool
	IsAfterSixMonths  bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type DeviceHistory struct {
	Id       int
	Device   Device
	IssuedBy string
	IssuedAt time.Time
}

func (device Device) Normalize() (Device, error) {
	// Remove spaces and uppercase serial
	device.Serie = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToUpper(r)
	}, device.Serie)

	// Validate serial
	if !(12 <= len(device.Serie) && len(device.Serie) <= 15) {
		return Device{}, errors.New("invalid serial length")
	}

	return device, nil

}

func (product Product) CalculateIndividualPrice(r currency.ExchangeRate) (int, error) {
	rate, ok := r.Get(product.Currency)
	if !ok {
		return 0, errors.New("invalid currency")
	}
	newPrice := float64(product.Price) * rate
	return int(math.Round(newPrice)), nil
}

type DeviceData struct {
	Id                int
	ProductSerial     string
	ProductType       string
	PartNoModel       string
	WarrantyStartDate time.Time
	WarrantyEndDate   time.Time
}

func toUpperAndRemoveSpaces(r rune) rune {
	if unicode.IsSpace(r) {
		return -1
	}
	return unicode.ToUpper(r)
}

func extractString(s string) (string, error) {
	startIndex := strings.Index(s, "[")
	if startIndex == -1 {
		return "", errors.New("bad product type format")
	}
	endIndex := strings.Index(s[startIndex:], "]")
	if endIndex == -1 {
		return "", errors.New("bad product type format")
	}
	endIndex += startIndex
	res := s[startIndex+1 : endIndex]
	return strings.Map(toUpperAndRemoveSpaces, res), nil
}

// Parse csv
// Header: RMA_SERIAL_NO, RMA_PRODUCT_TYPE_DESC, RMA_PART_NO_MODEL, WARRANTY_START_DATE, WARRANTY_END_DATE
func ParseDeviceDataRow(row []string) (dd DeviceData, err error) {
	// Query data
	rma_serial_no := row[0]
	rma_product_type_desc := row[1]
	rma_part_no_model := row[2]
	warranty_start_date := row[3]
	warranty_end_date := row[4]

	pType, err := extractString(rma_product_type_desc)
	if err != nil {
		return
	}
	pType, ok := config.ASUS_DEVICE_TYPE[pType]
	if !ok {
		err = errors.New("unknown product type")
		return
	}
	dd.ProductSerial = strings.Map(toUpperAndRemoveSpaces, rma_serial_no)
	dd.ProductType = pType
	dd.PartNoModel = strings.Map(toUpperAndRemoveSpaces, rma_part_no_model)
	dd.WarrantyStartDate, err = time.Parse("2006-01-02", warranty_start_date)
	if err != nil {
		return
	}
	dd.WarrantyEndDate, err = time.Parse("2006-01-02", warranty_end_date)
	if err != nil {
		return
	}
	return
}

func ParseDeviceDataCSV(data [][]string) ([]DeviceData, []error) {
	dds := make([]DeviceData, 0, len(data))
	errs := make([]error, 0, len(data))
	for i, row := range data {
		dd, err := ParseDeviceDataRow(row)
		if err != nil {
			errs = append(errs, fmt.Errorf("parse error in row %d: %w", i+1, err))
		} else {
			dds = append(dds, dd)
		}
	}
	return dds, errs
}
