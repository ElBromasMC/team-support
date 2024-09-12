package service

import (
	"alc/model/cart"
	"alc/model/checkout"
	"alc/model/store"
	"math"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Currency struct {
	db *pgxpool.Pool
}

func NewCurrencyService(db *pgxpool.Pool) Currency {
	return Currency{
		db: db,
	}
}

func (cs Currency) GetCurrency(slug string) (store.Currency, error) {
	if slug == "USD" {
		return store.USD, nil
	} else if slug == "PEN" {
		return store.PEN, nil
	} else {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Divisa inválida")
	}
}

// Returns p/q value
func (cs Currency) GetExchangeRate(p store.Currency, q store.Currency) (float64, error) {
	if p == q {
		return 1, nil
	}
	if p == store.USD && q == store.PEN {
		return 3.78, nil
	} else if p == store.PEN && q == store.USD {
		return 0.26, nil
	} else {
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}
}

// Order product management

func (cs Currency) CalculateOrderProductSubtotal(p checkout.OrderProduct, to store.Currency) (int, error) {
	rate, err := cs.GetExchangeRate(p.ProductCurrency, to)
	if err != nil {
		return 0, err
	}
	newPrice := float64(p.Quantity) * float64(p.ProductPrice) * rate
	return int(math.Round(newPrice)), nil
}

func (cs Currency) CalculateOrderProductsAmount(products []checkout.OrderProduct, to store.Currency) (int, error) {
	amount := 0
	for _, product := range products {
		subtotal, err := cs.CalculateOrderProductSubtotal(product, to)
		if err != nil {
			return 0, err
		}
		amount += subtotal
	}
	return amount, nil
}

// Cart item management

func (cs Currency) CalculateCartItemSubtotal(item cart.Item, to store.Currency) (int, error) {
	rate, err := cs.GetExchangeRate(item.Product.Currency, to)
	if err != nil {
		return 0, err
	}
	newPrice := float64(item.Quantity) * float64(item.Product.Price) * rate
	return int(math.Round(newPrice)), nil
}

func (cs Currency) CalculateCartItemsAmount(items []cart.Item, to store.Currency) (int, error) {
	amount := 0
	for _, item := range items {
		subtotal, err := cs.CalculateCartItemSubtotal(item, to)
		if err != nil {
			return 0, err
		}
		amount += subtotal
	}
	return amount, nil
}
