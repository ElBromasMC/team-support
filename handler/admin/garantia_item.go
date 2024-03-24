package admin

import (
	"alc/handler/util"
	"alc/model/store"
	"alc/view/admin/garantia"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GET "/admin/garantia/:categorySlug/:itemSlug"
func (h *Handler) HandleGarantiaItemShow(c echo.Context) error {
	// Parsing request
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	// Query category
	category, err := h.AdminService.GetCategory(store.GarantiaType, categorySlug)
	if err != nil {
		return err
	}

	// Query item
	item, err := h.AdminService.GetItem(category, itemSlug)
	if err != nil {
		return err
	}

	// Query products
	products, err := h.AdminService.GetProducts(item)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, garantia.ItemShow(item, products))
}

// POST "/admin/garantia/:categorySlug/:itemSlug"
func (h *Handler) HandleNewGarantiaProduct(c echo.Context) error {
	// Parsing request
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	var product store.Product
	product.Name = c.FormValue("name")
	priceFloat, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid price")
	}
	product.Price = int(math.Round(priceFloat * 100))

	// Query category
	category, err := h.AdminService.GetCategory(store.GarantiaType, categorySlug)
	if err != nil {
		return err
	}

	// Query and attach item
	item, err := h.AdminService.GetItem(category, itemSlug)
	if err != nil {
		return err
	}
	product.Item = item

	// Insert product into database
	if _, err := h.AdminService.InsertProduct(product); err != nil {
		return err
	}

	// Get updated products
	products, err := h.AdminService.GetProducts(item)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, garantia.ProductTableShow(products))
}

// PUT "/admin/garantia/:categorySlug/:itemSlug"
func (h *Handler) HandleUpdateGarantiaProduct(c echo.Context) error {
	// Parsing request
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	var product store.Product
	product.Name = c.FormValue("name")
	priceFloat, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid price")
	}
	product.Price = int(math.Round(priceFloat * 100))

	// Query category
	category, err := h.AdminService.GetCategory(store.GarantiaType, categorySlug)
	if err != nil {
		return err
	}

	// Query and attach item
	item, err := h.AdminService.GetItem(category, itemSlug)
	if err != nil {
		return err
	}
	product.Item = item

	// Update product
	if err := h.AdminService.UpdateProduct(id, product); err != nil {
		return err
	}

	// Get updated products
	products, err := h.AdminService.GetProducts(item)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, garantia.ProductTableShow(products))
}

// DELETE "/admin/garantia/:categorySlug/:itemSlug"
func (h *Handler) HandleRemoveGarantiaProduct(c echo.Context) error {
	// Parsing request
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Query category
	category, err := h.AdminService.GetCategory(store.GarantiaType, categorySlug)
	if err != nil {
		return err
	}

	// Query item
	item, err := h.AdminService.GetItem(category, itemSlug)
	if err != nil {
		return err
	}

	// Remove product
	if err := h.AdminService.RemoveProduct(id); err != nil {
		return err
	}

	// Get updated products
	products, err := h.AdminService.GetProducts(item)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, garantia.ProductTableShow(products))
}
