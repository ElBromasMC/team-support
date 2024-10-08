package currency

import (
	"context"
	"fmt"
)

// Currency management

type Currency string

const (
	USD Currency = "USD"
	PEN Currency = "PEN"
)

func GetCurrencies() []Currency {
	currencies := []Currency{
		USD,
		PEN,
	}
	return currencies
}

type Pair struct {
	Curr Currency
	Rate float64
}

type ExchangeRate struct {
	to    Currency
	table map[Currency]float64
}

func NewExchangeRate(to Currency, table map[Currency]float64) ExchangeRate {
	return ExchangeRate{
		to:    to,
		table: table,
	}
}

func (r ExchangeRate) To() Currency {
	return r.to
}

func (r ExchangeRate) Get(from Currency) (float64, bool) {
	rate, ok := r.table[from]
	return rate, ok
}

func (r ExchangeRate) GetTable() []Pair {
	pairs := make([]Pair, 0, len(r.table))
	for key, value := range r.table {
		pair := Pair{
			Curr: key,
			Rate: value,
		}
		pairs = append(pairs, pair)
	}
	return pairs
}

type RateKey struct{}

func GetExchangeRate(ctx context.Context) ExchangeRate {
	rate, ok := ctx.Value(RateKey{}).(ExchangeRate)
	if ok {
		return rate
	} else {
		return ExchangeRate{}
	}
}

// Miscellaneous

func DisplayPrice(codes map[Currency]string, value int, currency Currency) string {
	code, ok := codes[currency]
	if !ok {
		code = ""
	}

	return fmt.Sprintf("%s %.2f", code, float64(value)/100.0)
}
