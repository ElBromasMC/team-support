package cart

import (
	"alc/model/currency"
	"alc/model/store"
	"context"
	"errors"
	"math"
	"strings"
	"unicode"
)

const SessionName = "cart"

type ItemsKey struct{}

type Item struct {
	Product  store.Product
	Quantity int
	Details  map[string]string
}

type ItemRequest struct {
	ProductId int
	Quantity  int
	Details   map[string]string
}

func (item Item) Normalize() (Item, error) {
	if item.Product.Item.Category.Type == store.GarantiaType {
		if item.Quantity != 1 {
			return Item{}, errors.New("invalid quantity for warranty")
		}
		serie, ok := item.Details["Serie"]
		if !ok {
			return Item{}, errors.New("missing 'Serie' for warranty")
		}
		// Remove spaces and uppercase serial
		serie = strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return unicode.ToUpper(r)
		}, serie)

		// Validate serial
		if !(12 <= len(serie) && len(serie) <= 15) {
			return Item{}, errors.New("invalid 'Serie' for warranty")
		}

		// Attach normalized serial
		item.Details["Serie"] = serie
	} else {
		if item.Quantity < 1 {
			return Item{}, errors.New("invalid quantity for store item")
		}
		if item.Product.Stock != nil {
			if item.Quantity > *item.Product.Stock {
				return Item{}, errors.New("quantity exceeds current stock")
			}
		}
	}

	return item, nil
}

func (item Item) ToRequest() ItemRequest {
	return ItemRequest{
		ProductId: item.Product.Id,
		Quantity:  item.Quantity,
		Details:   item.Details,
	}
}

func GetItems(ctx context.Context) []Item {
	items, _ := ctx.Value(ItemsKey{}).([]Item)
	return items
}

// Cart item management

func (item Item) CalculateIndividualPrice(r currency.ExchangeRate) (int, error) {
	return item.Product.CalculateIndividualPrice(r)
}

func (item Item) CalculateSubtotal(r currency.ExchangeRate) (int, error) {
	rate, ok := r.Get(item.Product.Currency)
	if !ok {
		return 0, errors.New("invalid currency")
	}
	newPrice := float64(item.Quantity) * float64(item.Product.Price) * rate
	return int(math.Round(newPrice)), nil
}

func CalculateAmount(r currency.ExchangeRate, items []Item) (int, error) {
	amount := 0
	for _, item := range items {
		subtotal, err := item.CalculateSubtotal(r)
		if err != nil {
			return 0, err
		}
		amount += subtotal
	}
	return amount, nil
}
