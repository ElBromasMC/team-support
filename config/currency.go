package config

import "alc/model/currency"

const (
	STORE_CURRENCY currency.Currency = currency.PEN
)

var (
	CURRENCY_NUMERIC_CODES = map[currency.Currency]string{
		currency.PEN: "604",
		currency.USD: "840",
	}
)

var (
	CURRENCY_CODES = map[currency.Currency]string{
		currency.PEN: "PEN",
		currency.USD: "US$",
	}
)
