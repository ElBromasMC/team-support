package store

import (
	"alc/config"
	"alc/handler/util"
	"alc/model/currency"
	"alc/model/store"
	"alc/view/admin/store/bulk"
	"alc/view/admin/store/bulk/asus"
	"alc/view/admin/store/bulk/product"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleBulkLoaderShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, bulk.Show(t))
}

// Products

func (h *Handler) HandleBulkLoaderProductsShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, product.Show(t))
}

func (h *Handler) HandleBulkLoaderProductsInsertion(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	productsJson := c.FormValue("products")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}

	var products []store.Product
	if t == store.GarantiaType {
	} else {
		return nil
	}

	// Decode the products
	if err := json.Unmarshal([]byte(productsJson), &products); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Productos no válidos")
	}

	// Insert the products
	productErrors := make([]store.Product, 0, len(products))
	errors := make([]error, 0, len(products))
	for _, product := range products {
		// Normalize data
		product, err := product.Normalize()
		if err != nil {
			productErrors = append(productErrors, product)
			errors = append(errors, err)
			continue
		}
		if len(product.Item.Name) == 0 {
			productErrors = append(productErrors, product)
			errors = append(errors, fmt.Errorf("debe proporcionar el item"))
			continue
		}

		// Insert item if not exists
		itemId, err := h.AdminService.InsertItemIfNotExists(product.Item.Category, slug.Make(product.Item.Name), product.Item.Name)
		if err != nil {
			productErrors = append(productErrors, product)
			errors = append(errors, err)
			continue
		}

		// Attach data
		product.Item.Id = itemId
		product.Currency = currency.USD

		if _, err := h.AdminService.InsertProduct(product); err != nil {
			productErrors = append(productErrors, product)
			errors = append(errors, err)
			continue
		}
	}

	return util.Render(c, http.StatusOK, product.ErrorsShow(t, productErrors, errors))
}

func (h *Handler) HandleBulkLoaderProductsPreview(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	file, err := c.FormFile("products")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Debe proporcionar los productos")
	}

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}

	// Data source
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al abrir los productos")
	}
	defer src.Close()

	reader := csv.NewReader(src)

	var products []store.Product
	if t == store.GarantiaType {
		reader.FieldsPerRecord = 7

		records, err := reader.ReadAll()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Error al leer los productos")
		}

		products = make([]store.Product, 0, len(records))
		for _, row := range records {
			// Parsing data
			categorySlug := strings.ToLower(strings.TrimSpace(row[0]))
			itemName := strings.TrimSpace(row[1])

			var product store.Product
			product.Name = row[2]
			product.PartNumber = row[3]

			acceptBefore := strings.ToUpper(strings.TrimSpace(row[4]))
			acceptAfter := strings.ToUpper(strings.TrimSpace(row[5]))
			if acceptBefore == "SI" {
				product.AcceptBeforeSixMonths = true
			} else {
				product.AcceptBeforeSixMonths = false
			}
			if acceptAfter == "SI" {
				product.AcceptAfterSixMonths = true
			} else {
				product.AcceptAfterSixMonths = false
			}

			priceFloat, err := strconv.ParseFloat(row[6], 64)
			if err != nil {
				continue
			}
			product.Price = int(math.Round(priceFloat * 100))

			// Query data
			cat, err := h.AdminService.GetCategory(t, categorySlug)
			if err != nil {
				continue
			}
			i := store.Item{
				Category: cat,
				Name:     itemName,
			}

			// Attach data
			product.Item = i

			// Normalize data
			product, err = product.Normalize()
			if err != nil {
				continue
			}
			if len(product.Item.Name) == 0 {
				continue
			}

			products = append(products, product)
		}
	} else {
		return nil
	}

	// Encode the products
	encProducts, err := json.Marshal(products)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al codificar los productos")
	}

	return util.Render(c, http.StatusOK, product.Preview(t, products, encProducts))
}

func (h *Handler) HandleBulkLoaderProductsInsertionFormShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return nil
	}

	if t == store.GarantiaType {
		return util.Render(c, http.StatusOK, product.InsertionForm(t))
	} else {
		return nil
	}
}

// ASUS

type servicePair struct {
	Key   string
	Value int
}

func parseServiceContent(input string) []servicePair {
	// Remove spaces and uppercase the input
	input = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToUpper(r)
	}, input)
	var result []servicePair
	pairs := strings.Split(input, ";")
	for _, pair := range pairs {
		parts := strings.Split(pair, "-")
		if len(parts) != 2 {
			continue
		}
		tag := parts[0]
		valueStr := strings.TrimSuffix(parts[1], "M")
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			continue
		}
		result = append(result, servicePair{
			Key:   tag,
			Value: value,
		})
	}
	return result
}

func (h *Handler) HandleBulkLoaderAsusShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return nil
	}

	if t == store.GarantiaType {
		return util.Render(c, http.StatusOK, asus.Show(t))
	} else {
		return nil
	}
}

func (h *Handler) HandleBulkLoaderAsusInsertion(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	productsJson := c.FormValue("products")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}

	var products []store.Product
	if t == store.GarantiaType {
	} else {
		return nil
	}

	// Decode the products
	if err := json.Unmarshal([]byte(productsJson), &products); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Productos no válidos")
	}

	// Insert the products
	productErrors := make([]store.Product, 0, len(products))
	errors := make([]error, 0, len(products))
	for _, product := range products {
		// Normalize data
		product, err := product.Normalize()
		if err != nil {
			productErrors = append(productErrors, product)
			errors = append(errors, err)
			continue
		}
		if len(product.Item.Name) == 0 {
			productErrors = append(productErrors, product)
			errors = append(errors, fmt.Errorf("debe proporcionar el item"))
			continue
		}

		// Insert category if not exists
		categoryId, err := h.AdminService.InsertCategoryIfNotExists(t, product.Item.Category.Slug, product.Item.Category.Name)
		if err != nil {
			productErrors = append(productErrors, product)
			errors = append(errors, err)
			continue
		}
		product.Item.Category.Id = categoryId

		// Insert item if not exists
		itemId, err := h.AdminService.InsertItemIfNotExists(product.Item.Category, product.Item.Slug, product.Item.Name)
		if err != nil {
			productErrors = append(productErrors, product)
			errors = append(errors, err)
			continue
		}
		product.Item.Id = itemId
		product.Currency = currency.USD

		// Insert the product
		if _, err := h.AdminService.InsertProductWithSlug(product); err != nil {
			productErrors = append(productErrors, product)
			errors = append(errors, err)
			continue
		}
	}

	return util.Render(c, http.StatusOK, asus.ErrorsShow(t, productErrors, errors))
}

func (h *Handler) HandleBulkLoaderAsusPreview(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	file, err := c.FormFile("products")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Debe proporcionar los productos")
	}

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}

	// Data source
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al abrir los productos")
	}
	defer src.Close()

	reader := csv.NewReader(src)

	var products []store.Product
	if t == store.GarantiaType {
		reader.FieldsPerRecord = 7

		records, err := reader.ReadAll()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Error al leer los productos")
		}

		products = make([]store.Product, 0, len(records))
	RowLoop:
		for _, row := range records {
			// Parsing data
			productCode := strings.ToUpper(strings.TrimSpace(row[0]))
			serviceContent := strings.TrimSpace(row[1])
			stdPeriodStr := strings.TrimSuffix(strings.ToUpper(strings.TrimSpace(row[6])), "M")
			stdPeriod, err := strconv.Atoi(stdPeriodStr)
			if err != nil {
				continue
			}

			var product store.Product
			acceptBefore := strings.ToUpper(strings.TrimSpace(row[2]))
			acceptAfter := strings.ToUpper(strings.TrimSpace(row[3]))
			if acceptBefore == "SI" {
				product.AcceptBeforeSixMonths = true
			} else if acceptBefore == "NO" {
				product.AcceptBeforeSixMonths = false
			} else {
				continue
			}
			if acceptAfter == "SI" {
				product.AcceptAfterSixMonths = true
			} else if acceptAfter == "NO" {
				product.AcceptAfterSixMonths = false
			} else {
				continue
			}
			product.PartNumber = row[4]
			priceFloat, err := strconv.ParseFloat(row[5], 64)
			if err != nil {
				continue
			}
			product.Price = int(math.Round(priceFloat * 100))

			// Get category
			catName, ok := config.ASUS_CODES[productCode]
			if !ok {
				continue
			}
			product.Item.Category.Type = t
			product.Item.Category.Slug = slug.Make(productCode)
			product.Item.Category.Name = catName

			// Get item and product
			serviceContentMap := parseServiceContent(serviceContent)
			if len(serviceContentMap) == 0 {
				continue
			}
			itemSlug := make([]string, 0, len(serviceContentMap))
			itemName := make([]string, 0, len(serviceContentMap))
			productSlug := make([]string, 0, len(serviceContentMap))
			productName := make([]string, 0, len(serviceContentMap))
			for i, content := range serviceContentMap {
				key := content.Key
				value := content.Value
				longName, ok := config.ASUS_LONGNAMES[key]
				if !ok {
					continue RowLoop
				}
				shortName, ok := config.ASUS_SHORTNAMES[key]
				if !ok {
					continue RowLoop
				}
				// Item
				itemSlug = append(itemSlug, key)
				if i == 0 {
					itemName = append(itemName, longName)
				} else {
					itemName = append(itemName, shortName)
				}
				// Product
				productSlug = append(productSlug, fmt.Sprintf("%d", value))
				productName = append(productName, fmt.Sprintf("%s %d meses", shortName, value))
			}
			product.Item.Slug = slug.Make(strings.Join(itemSlug, "-"))
			product.Item.Name = strings.Join(itemName, " + ")
			product.Slug = slug.Make(strings.Join(productSlug, "-") + fmt.Sprintf("_%d", stdPeriod))
			product.Name = strings.Join(productName, " + ") + fmt.Sprintf(" (Total: %d meses)", stdPeriod)

			// Normalize data
			product, err = product.Normalize()
			if err != nil {
				continue
			}

			products = append(products, product)
		}
	} else {
		return nil
	}

	// Encode the products
	encProducts, err := json.Marshal(products)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al codificar los productos")
	}

	return util.Render(c, http.StatusOK, asus.Preview(t, products, encProducts))
}

func (h *Handler) HandleBulkLoaderAsusInsertionFormShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return nil
	}

	if t == store.GarantiaType {
		return util.Render(c, http.StatusOK, asus.InsertionForm(t))
	} else {
		return nil
	}
}
