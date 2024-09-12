package config

import "alc/model/store"

const (
	STORE_CURRENCY store.Currency = store.PEN
)

var (
	CURRENCY_NUMERIC_CODES = map[store.Currency]string{
		store.PEN: "604",
		store.USD: "840",
	}
)

var (
	CURRENCY_CODES = map[store.Currency]string{
		store.PEN: "PEN",
		store.USD: "US$",
	}
)
