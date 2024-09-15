package service

import (
	"alc/model/currency"
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
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

func (cs Currency) SetExchangeRate(from currency.Currency, to currency.Currency, rate float64) error {
	// Check if exists
	var exists bool
	sql := `SELECT EXISTS (
		SELECT 1 FROM exchange_rates
		WHERE base_currency = $1 AND target_currency = $2
	)`
	if err := cs.db.QueryRow(context.Background(), sql, from, to).Scan(&exists); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !exists {
		sql := `INSERT INTO exchange_rates
		(base_currency, target_currency, rate)
		VALUES ($1, $2, $3)`
		if _, err := cs.db.Exec(context.Background(), sql, from, to, rate); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		sql := `UPDATE exchange_rates
		SET rate = $1
		WHERE base_currency = $2 AND target_currency = $3`
		if _, err := cs.db.Exec(context.Background(), sql, rate, from, to); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return nil
}

func (cs Currency) GetExchangeRate(to currency.Currency) (currency.ExchangeRate, error) {
	sql := `SELECT base_currency, rate
	FROM exchange_rates
	WHERE target_currency = $1`

	rows, err := cs.db.Query(context.Background(), sql, to)
	if err != nil {
		return currency.ExchangeRate{}, echo.NewHTTPError(http.StatusInternalServerError, "Error en el servidor de tasa de cambios")
	}
	defer rows.Close()

	pairs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (currency.Pair, error) {
		var pair currency.Pair
		err := row.Scan(&pair.Curr, &pair.Rate)
		return pair, err
	})
	if err != nil {
		return currency.ExchangeRate{}, echo.NewHTTPError(http.StatusInternalServerError, "Error en el servidor de tasa de cambios")
	}

	table := make(map[currency.Currency]float64, len(pairs)+1)
	for _, pair := range pairs {
		table[pair.Curr] = pair.Rate
	}
	table[to] = 1

	rate := currency.NewExchangeRate(to, table)
	return rate, nil
}
