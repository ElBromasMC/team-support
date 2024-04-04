package store

import (
	"alc/handler/util"
	"alc/model/store"
	"alc/view/admin/store/product"
	"math"
	"net/http"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

// GET "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products"
func (h *Handler) HandleProductsShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}
	i, err := h.AdminService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}

	// Query products
	products, err := h.AdminService.GetProducts(i)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, product.Show(i, products))
}

// POST "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products"
func (h *Handler) HandleProductInsertion(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	var p store.Product
	p.Name = c.FormValue("name")
	priceFloat, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid price")
	}
	p.Price = int(math.Round(priceFloat * 100))

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}
	i, err := h.AdminService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}

	// Attach data
	p.Item = i
	p.Slug = slug.Make(p.Name)

	// Insert product into database
	if _, err := h.AdminService.InsertProduct(p); err != nil {
		return err
	}

	// Get updated products
	products, err := h.AdminService.GetProducts(i)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, product.Table(i, products))
}

// PUT "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/:productSlug"
func (h *Handler) HandleProductUpdate(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")
	productSlug := c.Param("productSlug")

	var p store.Product
	p.Name = c.FormValue("name")
	priceFloat, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid price")
	}
	p.Price = int(math.Round(priceFloat * 100))

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}
	i, err := h.AdminService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}
	oldProduct, err := h.AdminService.GetProduct(i, productSlug)
	if err != nil {
		return err
	}

	// Attach data
	p.Item = i
	p.Slug = slug.Make(p.Name)

	// Update product
	if err := h.AdminService.UpdateProduct(oldProduct.Id, p); err != nil {
		return err
	}

	// Get updated products
	products, err := h.AdminService.GetProducts(i)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, product.Table(i, products))
}

// DELETE "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/:productSlug"
func (h *Handler) HandleProductDeletion(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")
	productSlug := c.Param("productSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}
	i, err := h.AdminService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}
	p, err := h.AdminService.GetProduct(i, productSlug)
	if err != nil {
		return err
	}

	// Remove product
	if err := h.AdminService.RemoveProduct(p.Id); err != nil {
		return err
	}

	// Get updated products
	products, err := h.AdminService.GetProducts(i)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, product.Table(i, products))
}

// GET "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/insert"
func (h *Handler) HandleProductInsertionFormShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}
	i, err := h.AdminService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, product.InsertionForm(i))
}

// GET "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/:productSlug/update"
func (h *Handler) HandleProductUpdateFormShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")
	productSlug := c.Param("productSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}
	i, err := h.AdminService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}
	p, err := h.AdminService.GetProduct(i, productSlug)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, product.UpdateForm(p))
}

// GET "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/:productSlug/delete"
func (h *Handler) HandleProductDeletionFormShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")
	productSlug := c.Param("productSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}
	i, err := h.AdminService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}
	p, err := h.AdminService.GetProduct(i, productSlug)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, product.DeletionForm(p))
}
