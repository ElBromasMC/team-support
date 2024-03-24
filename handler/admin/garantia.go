package admin

import (
	"alc/handler/util"
	"alc/model/store"
	"alc/view/admin/garantia"
	"math"
	"net/http"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

// GET "/admin/garantia"
func (h *Handler) HandleGarantiaShow(c echo.Context) error {
	cats, err := h.AdminService.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.Show(cats))
}

// POST "/admin/garantia"
func (h *Handler) HandleNewGarantiaCategory(c echo.Context) error {
	// Parsing request
	var cat store.Category
	cat.Type = store.GarantiaType
	cat.Name = c.FormValue("name")
	cat.Description = c.FormValue("description")

	img, imgErr := c.FormFile("img")

	// Generate slug from name
	cat.Slug = slug.Make(cat.Name)

	// Insert and attach image if present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		cat.Img = newImg
	}

	// Insert it into database
	if _, err := h.AdminService.InsertCategory(cat); err != nil {
		return err
	}

	// Get updated categories
	cats, err := h.AdminService.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.CategoryTableShow(cats))
}

// PUT "/admin/garantia"
func (h *Handler) HandleUpdateGarantiaCategory(c echo.Context) error {
	// Parsing request
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	var cat store.Category
	cat.Type = store.GarantiaType
	cat.Name = c.FormValue("name")
	cat.Description = c.FormValue("description")

	img, imgErr := c.FormFile("img")

	// Generate slug from name
	cat.Slug = slug.Make(cat.Name)

	// Insert and attach image if present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		cat.Img = newImg
	}

	// Update category
	if err := h.AdminService.UpdateCategory(id, cat); err != nil {
		return err
	}

	// Get updated categories
	cats, err := h.AdminService.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.CategoryTableShow(cats))
}

// GET "/admin/garantia/:slug"
func (h *Handler) HandleGarantiaCategoryShow(c echo.Context) error {
	// Parsing request
	catSlug := c.Param("slug")

	cat, err := h.AdminService.GetCategory(store.GarantiaType, catSlug)
	if err != nil {
		return err
	}

	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.CategoryShow(cat, items))
}

// POST "/admin/garantia/:slug"
func (h *Handler) HandleNewGarantiaItem(c echo.Context) error {
	// Parsing request
	catSlug := c.Param("slug")

	var item store.Item
	item.Name = c.FormValue("name")
	item.Description = c.FormValue("description")

	img, imgErr := c.FormFile("img")

	// Query and attach category
	cat, err := h.AdminService.GetCategory(store.GarantiaType, catSlug)
	if err != nil {
		return err
	}
	item.Category = cat

	// Generate slug from name
	item.Slug = slug.Make(item.Name)

	// Insert and attach image if present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		item.Img = newImg
	}

	// Insert it into database
	if _, err := h.AdminService.InsertItem(item); err != nil {
		return err
	}

	// Get updated items
	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.ItemTableShow(items))
}

// PUT "/admin/garantia/:slug"
func (h *Handler) HandleUpdateGarantiaItem(c echo.Context) error {
	// Parsing request
	catSlug := c.Param("slug")

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	var item store.Item
	item.Name = c.FormValue("name")
	item.Description = c.FormValue("description")

	img, imgErr := c.FormFile("img")

	// Generate slug from name
	item.Slug = slug.Make(item.Name)

	// Query and attach category
	cat, err := h.AdminService.GetCategory(store.GarantiaType, catSlug)
	if err != nil {
		return err
	}
	item.Category = cat

	// Insert and attach image if present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		item.Img = newImg
	}

	// Update item
	if err := h.AdminService.UpdateItem(id, item); err != nil {
		return err
	}

	// Get updated items
	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.ItemTableShow(items))
}

// DELETE "/admin/garantia/:slug"
func (h *Handler) HandleRemoveGarantiaItem(c echo.Context) error {
	// Parsing request
	catSlug := c.Param("slug")

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	// Query category
	cat, err := h.AdminService.GetCategory(store.GarantiaType, catSlug)
	if err != nil {
		return err
	}

	// Remove item
	if err := h.AdminService.RemoveItem(id); err != nil {
		return err
	}

	// Get updated items
	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.ItemTableShow(items))
}

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
