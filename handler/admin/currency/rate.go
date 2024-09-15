package currency

import (
	"alc/config"
	"alc/handler/util"
	model "alc/model/currency"
	"alc/view/admin/currency"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleRateShow(c echo.Context) error {
	// Get exchange rate
	rate, err := h.CurrencyService.GetExchangeRate(config.STORE_CURRENCY)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, currency.Show(rate))
}

func (h *Handler) HandleRateUpdate(c echo.Context) error {
	baseCurrency, err := h.CurrencyService.GetCurrency(c.FormValue("base_currency"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Divisa base inválida")
	}
	targetCurrency, err := h.CurrencyService.GetCurrency(c.FormValue("target_currency"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Divisa objetivo inválida")
	}
	rateFloat, err := strconv.ParseFloat(c.FormValue("rate"), 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Tasa de cambio inválida")
	}

	if err := h.CurrencyService.SetExchangeRate(baseCurrency, targetCurrency, rateFloat); err != nil {
		return err
	}

	// Get exchange rate
	rate, err := h.CurrencyService.GetExchangeRate(config.STORE_CURRENCY)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, currency.Table(rate))
}

func (h *Handler) HandleRateUpdateFormShow(c echo.Context) error {
	// Get exchange rate
	rate, err := h.CurrencyService.GetExchangeRate(config.STORE_CURRENCY)
	if err != nil {
		return err
	}

	currencies := model.GetCurrencies()

	return util.Render(c, http.StatusOK, currency.UpdateForm(rate, currencies))
}
