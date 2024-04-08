package service

import (
	"alc/model/checkout"
	"alc/model/store"
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Public struct {
	DB *pgxpool.Pool
}

func (ps Public) GetType(slug string) (store.Type, error) {
	var t store.Type
	if slug == store.GarantiaType.ToSlug() {
		t = store.GarantiaType
	} else if slug == store.StoreType.ToSlug() {
		t = store.StoreType
	} else {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Invalid type")
	}
	return t, nil
}

// Check
func (ps Public) GetCategories(t store.Type) ([]store.Category, error) {
	rows, err := ps.DB.Query(context.Background(), `SELECT sc.id, sc.name, sc.description, sc.slug, img.id, img.filename
FROM store_categories AS sc
LEFT JOIN images AS img
ON sc.img_id = img.id
WHERE sc.type = $1
ORDER BY sc.id DESC`, t)
	if err != nil {
		return []store.Category{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	cats, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Category, error) {
		var cat store.Category
		var imgId *int
		var imgFilename *string
		err := row.Scan(&cat.Id, &cat.Name, &cat.Description, &cat.Slug, &imgId, &imgFilename)
		if imgId != nil {
			cat.Img.Id = *imgId
			cat.Img.Filename = *imgFilename
		}
		cat.Type = t
		return cat, err
	})
	if err != nil {
		return []store.Category{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return cats, nil
}

// Check missing (Img)
func (ps Public) GetCategory(t store.Type, slug string) (store.Category, error) {
	var cat store.Category
	if err := ps.DB.QueryRow(context.Background(), `SELECT id, name, description
FROM store_categories
WHERE type = $1 AND slug = $2`, t, slug).Scan(&cat.Id, &cat.Name, &cat.Description); err != nil {
		return store.Category{}, echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}
	cat.Type = t
	cat.Slug = slug
	return cat, nil
}

// Check
func (ps Public) GetCategoryById(id int) (store.Category, error) {
	var cat store.Category
	var imgId *int
	var imgFilename *string
	if err := ps.DB.QueryRow(context.Background(), `SELECT sc.type, sc.name, sc.description, sc.slug, img.id, img.filename
FROM store_categories AS sc
LEFT JOIN images AS img
ON sc.img_id = img.id
WHERE sc.id = $1`, id).Scan(&cat.Type, &cat.Name, &cat.Description, &cat.Slug, &imgId, &imgFilename); err != nil {
		return store.Category{}, echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}
	if imgId != nil {
		cat.Img.Id = *imgId
		cat.Img.Filename = *imgFilename
	}
	cat.Id = id
	return cat, nil
}

func (ps Public) GetAllItemsLike(t store.Type, like string, page int, n int) ([]store.Item, error) {
	var rows pgx.Rows
	var err error
	if len(like) != 0 {
		rows, err = ps.DB.Query(context.Background(), `SELECT sc.slug, si.id, si.name, si.slug, img.id, img.filename
		FROM store_items AS si
		JOIN store_categories AS sc
		ON si.category_id = sc.id
		LEFT JOIN images AS img
		ON si.img_id = img.id
		WHERE sc.type = $1
		AND si.name % $2
		ORDER BY si.name <-> $2
		LIMIT ($4 + 1) OFFSET ($3 - 1) * $4`, t, like, page, n)
	} else {
		rows, err = ps.DB.Query(context.Background(), `SELECT sc.slug, si.id, si.name, si.slug, img.id, img.filename
		FROM store_items AS si
		JOIN store_categories AS sc
		ON si.category_id = sc.id
		LEFT JOIN images AS img
		ON si.img_id = img.id
		WHERE sc.type = $1
		ORDER BY si.id DESC
		LIMIT ($3 + 1) OFFSET ($2 - 1) * $3`, t, page, n)
	}

	if err != nil {
		return []store.Item{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Item, error) {
		var item store.Item
		var imgId *int
		var imgFilename *string
		err := row.Scan(&item.Category.Slug, &item.Id, &item.Name, &item.Slug, &imgId, &imgFilename)
		if imgId != nil {
			item.Img.Id = *imgId
			item.Img.Filename = *imgFilename
		}
		return item, err
	})
	if err != nil {
		return []store.Item{}, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return items, nil
}

func (ps Public) GetItemsLike(cat store.Category, like string, page int, n int) ([]store.Item, error) {
	var rows pgx.Rows
	var err error
	if len(like) != 0 {
		rows, err = ps.DB.Query(context.Background(), `SELECT sc.slug, si.id, si.name, si.slug, img.id, img.filename
		FROM store_items AS si
		JOIN store_categories AS sc
		ON si.category_id = sc.id
		LEFT JOIN images AS img
		ON si.img_id = img.id
		WHERE si.category_id = $1
		AND si.name % $2
		ORDER BY si.name <-> $2
		LIMIT ($4 + 1) OFFSET ($3 - 1) * $4`, cat.Id, like, page, n)
	} else {
		rows, err = ps.DB.Query(context.Background(), `SELECT sc.slug, si.id, si.name, si.slug, img.id, img.filename
		FROM store_items AS si
		JOIN store_categories AS sc
		ON si.category_id = sc.id
		LEFT JOIN images AS img
		ON si.img_id = img.id
		WHERE si.category_id = $1
		ORDER BY si.id DESC
		LIMIT ($3 + 1) OFFSET ($2 - 1) * $3`, cat.Id, page, n)
	}

	if err != nil {
		return []store.Item{}, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Item, error) {
		var item store.Item
		var imgId *int
		var imgFilename *string
		err := row.Scan(&item.Category.Slug, &item.Id, &item.Name, &item.Slug, &imgId, &imgFilename)
		if imgId != nil {
			item.Img.Id = *imgId
			item.Img.Filename = *imgFilename
		}
		return item, err
	})
	if err != nil {
		return []store.Item{}, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return items, nil
}

// Check
func (ps Public) GetItems(cat store.Category) ([]store.Item, error) {
	rows, err := ps.DB.Query(context.Background(), `SELECT si.id, si.name, si.description, si.long_description, si.slug,
img.id, img.filename, largeimg.id, largeimg.filename
FROM store_items AS si
LEFT JOIN images AS img
ON si.img_id = img.id
LEFT JOIN images AS largeimg
ON si.largeimg_id = largeimg.id
WHERE si.category_id = $1
ORDER BY si.id DESC`, cat.Id)
	if err != nil {
		return []store.Item{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Item, error) {
		var item store.Item
		var imgId *int
		var imgFilename *string
		var largeimgId *int
		var largeimgFilename *string
		err := row.Scan(&item.Id, &item.Name, &item.Description, &item.LongDescription, &item.Slug,
			&imgId, &imgFilename, &largeimgId, &largeimgFilename)
		if imgId != nil {
			item.Img.Id = *imgId
			item.Img.Filename = *imgFilename
		}
		if largeimgId != nil {
			item.LargeImg.Id = *largeimgId
			item.LargeImg.Filename = *largeimgFilename
		}
		item.Category = cat
		return item, err
	})
	if err != nil {
		return []store.Item{}, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return items, nil
}

// Check
func (ps Public) GetItem(cat store.Category, slug string) (store.Item, error) {
	var item store.Item
	var imgId *int
	var imgFilename *string
	var largeimgId *int
	var largeimgFilename *string
	if err := ps.DB.QueryRow(context.Background(), `SELECT si.id, si.name, si.description, si.long_description,
img.id, img.filename, largeimg.id, largeimg.filename
FROM store_items AS si
LEFT JOIN images AS img
ON si.img_id = img.id
LEFT JOIN images AS largeimg
ON si.largeimg_id = largeimg.id
WHERE si.category_id = $1 AND si.slug = $2`, cat.Id, slug).Scan(&item.Id, &item.Name, &item.Description, &item.LongDescription,
		&imgId, &imgFilename, &largeimgId, &largeimgFilename); err != nil {
		return store.Item{}, echo.NewHTTPError(http.StatusNotFound, "Item not found")
	}
	if imgId != nil {
		item.Img.Id = *imgId
		item.Img.Filename = *imgFilename
	}
	if largeimgId != nil {
		item.LargeImg.Id = *largeimgId
		item.LargeImg.Filename = *largeimgFilename
	}
	item.Category = cat
	item.Slug = slug
	return item, nil
}

// Check
func (ps Public) GetItemById(id int) (store.Item, error) {
	var item store.Item

	var imgId *int
	var imgFilename *string

	var largeimgId *int
	var largeimgFilename *string

	var catId int
	if err := ps.DB.QueryRow(context.Background(), `SELECT si.name, si.description, si.long_description, si.slug,
img.id, img.filename, largeimg.id, largeimg.filename, si.category_id
FROM store_items AS si
LEFT JOIN images AS img
ON si.img_id = img.id
LEFT JOIN images AS largeimg
ON si.largeimg_id = largeimg.id
WHERE si.id = $1`, id).Scan(&item.Name, &item.Description, &item.LongDescription, &item.Slug,
		&imgId, &imgFilename, &largeimgId, &largeimgFilename, &catId); err != nil {
		return store.Item{}, echo.NewHTTPError(http.StatusNotFound, "Item not found")
	}

	if imgId != nil {
		item.Img.Id = *imgId
		item.Img.Filename = *imgFilename
	}

	if largeimgId != nil {
		item.LargeImg.Id = *largeimgId
		item.LargeImg.Filename = *largeimgFilename
	}
	item.Id = id

	// Query and attach category
	cat, _ := ps.GetCategoryById(catId)
	item.Category = cat

	return item, nil
}

// Check
func (ps Public) GetProducts(item store.Item) ([]store.Product, error) {
	rows, err := ps.DB.Query(context.Background(), `SELECT id, name, price, details, slug
FROM store_products
WHERE item_id = $1
ORDER BY id ASC`, item.Id)
	if err != nil {
		return []store.Product{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Product, error) {
		var product store.Product
		product.Details = make(map[string]string)

		var detailsHstore pgtype.Hstore
		err := row.Scan(&product.Id, &product.Name, &product.Price, &detailsHstore, &product.Slug)
		for key, value := range detailsHstore {
			if value != nil {
				product.Details[key] = *value
			}
		}
		product.Item = item
		return product, err
	})
	if err != nil {
		return []store.Product{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return products, nil
}

// Check
func (ps Public) GetProduct(i store.Item, slug string) (store.Product, error) {
	var product store.Product
	product.Details = make(map[string]string)

	var detailsHstore pgtype.Hstore
	if err := ps.DB.QueryRow(context.Background(), `SELECT id, name, price, details
FROM store_products
WHERE item_id = $1 AND slug = $2`, i.Id, slug).Scan(&product.Id, &product.Name, &product.Price, &detailsHstore); err != nil {
		return store.Product{}, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	for key, value := range detailsHstore {
		if value != nil {
			product.Details[key] = *value
		}
	}
	product.Item = i
	product.Slug = slug

	return product, nil
}

// Check
func (ps Public) GetProductById(id int) (store.Product, error) {
	var product store.Product
	product.Details = make(map[string]string)

	var itemId int
	var detailsHstore pgtype.Hstore
	if err := ps.DB.QueryRow(context.Background(), `SELECT item_id, name, price, details, slug
FROM store_products
WHERE id = $1`, id).Scan(&itemId, &product.Name, &product.Price, &detailsHstore, &product.Slug); err != nil {
		return store.Product{}, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	for key, value := range detailsHstore {
		if value != nil {
			product.Details[key] = *value
		}
	}
	product.Id = id

	// Query and attach item
	item, _ := ps.GetItemById(itemId)
	product.Item = item

	return product, nil
}

func (ps Public) InsertOrderProducts(order checkout.Order, products []checkout.OrderProduct) error {
	tx, err := ps.DB.Begin(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer tx.Rollback(context.Background())

	// Insert order
	var orderId int
	if err := tx.QueryRow(context.Background(), `INSERT INTO store_orders (email, phone_number, name, address, city, postal_code)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING purchase_order`, order.Email, order.Phone, order.Name, order.Address, order.City, order.PostalCode).Scan(&orderId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error inserting new order")
	}

	if _, err := tx.CopyFrom(
		context.Background(),
		pgx.Identifier{"order_products"},
		[]string{"order_id", "quantity", "product_type", "product_category", "product_item", "product_name", "product_price"},
		pgx.CopyFromSlice(len(products), func(i int) ([]any, error) {
			// TODO: hstoreDetails, hstoreProductDetails
			hstoreDetails := make(pgtype.Hstore, len(products[i].Details))
			for key, val := range products[i].Details {
				valCopy := val
				hstoreDetails[key] = &valCopy
			}

			hstoreProductDetails := make(pgtype.Hstore, len(products[i].ProductDetails))
			for key, val := range products[i].ProductDetails {
				valCopy := val
				hstoreProductDetails[key] = &valCopy
			}

			return []any{orderId, products[i].Quantity, products[i].ProductType, products[i].ProductCategory,
				products[i].ProductItem, products[i].ProductName, products[i].ProductPrice}, nil
		}),
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error inserting order products")
	}

	if err := tx.Commit(context.Background()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}
