package public

import (
	"alc/config"
	"alc/handler/util"
	"alc/model/store"
	"alc/view/garantia"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"unicode"

	"github.com/labstack/echo/v4"
)

// GET "/garantia"
func (h *Handler) HandleGarantiaShow(c echo.Context) error {
	cats, err := h.PublicService.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.Show(cats))
}

// GET "/garantia/:slug"
func (h *Handler) HandleGarantiaCategoryShow(c echo.Context) error {
	slug := c.Param("slug")

	cat, err := h.PublicService.GetCategory(store.GarantiaType, slug)
	if err != nil {
		return err
	}

	items, err := h.PublicService.GetItems(cat)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, garantia.ShowCategory(cat, items))
}

// GET "/garantia/:categorySlug/:itemSlug"
func (h *Handler) HandleGarantiaItemShow(c echo.Context) error {
	// Parse request
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")
	productStr := c.QueryParam("productId")

	productId := 0
	if len(productStr) != 0 {
		var err error
		productId, err = strconv.Atoi(productStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Id de producto inválido")
		}
	}

	cat, err := h.PublicService.GetCategory(store.GarantiaType, categorySlug)
	if err != nil {
		return err
	}

	item, err := h.PublicService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}

	products, err := h.PublicService.GetProducts(item)
	if err != nil {
		return err
	}

	// Get exchange rate
	rate, err := h.CurrencyService.GetExchangeRate(config.STORE_CURRENCY)
	if err != nil {
		return err
	}

	// Get selected index
	defaultIndex := 0
	for n, p := range products {
		if p.Id == productId {
			defaultIndex = n
			break
		}
	}

	return util.Render(c, http.StatusOK, garantia.ShowItem(item, products, defaultIndex, rate))
}

func (h *Handler) HandleGarantiaPartNumberRedirection(c echo.Context) error {
	// Parse request
	partNumber := c.FormValue("PartNumber")

	// Remove spaces and uppercase part number
	partNumber = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToUpper(r)
	}, partNumber)

	// Query data
	p, err := h.PublicService.GetProductByPartNumber(partNumber)
	if err != nil {
		return err
	}

	// Redirect to product
	redirectPath := path.Join("/garantia", p.Item.Category.Slug, p.Item.Slug) + fmt.Sprintf("?productId=%d", p.Id)
	_, ok := c.Request().Header[http.CanonicalHeaderKey("HX-Request")]
	if !ok {
		return c.Redirect(http.StatusFound, redirectPath)
	}
	c.Response().Header().Set("HX-Redirect", redirectPath)
	return c.NoContent(http.StatusOK)
}
