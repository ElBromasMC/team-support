package service

import (
	"alc/model/currency"
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

func (cs Currency) GetCurrency(slug string) (currency.Currency, error) {
	if slug == "USD" {
		return currency.USD, nil
	} else if slug == "PEN" {
		return currency.PEN, nil
	} else {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Divisa inválida")
	}
}

// TODO: connection with database
func (cs Currency) GetExchangeRate(to currency.Currency) (currency.ExchangeRate, error) {
	if to == currency.PEN {
		table := map[currency.Currency]float64{
			currency.PEN: 1,
			currency.USD: 3.78,
		}
		rate := currency.NewExchangeRate(to, table)
		return rate, nil
	} else {
		return currency.ExchangeRate{}, echo.NewHTTPError(http.StatusBadRequest, "Divisa inválida")
	}
}
