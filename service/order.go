package service

import (
	"alc/model/checkout"
	"alc/model/store"
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Order struct {
	db *pgxpool.Pool
}

func NewOrderService(db *pgxpool.Pool) Order {
	return Order{
		db: db,
	}
}

// Order management

func (os Order) GetOrderById(id uuid.UUID) (checkout.Order, error) {
	var order checkout.Order
	if err := os.db.QueryRow(context.Background(), `SELECT id, purchase_order, email, phone_number,
name, address, city, postal_code, created_at, sync_status, locked_at
FROM store_orders
WHERE id = $1`, id).Scan(&order.Id, &order.PurchaseOrder, &order.Email, &order.Phone,
		&order.Name, &order.Address, &order.City, &order.PostalCode, &order.CreatedAt, &order.SyncStatus, &order.LockedAt); err != nil {
		return checkout.Order{}, echo.NewHTTPError(http.StatusNotFound, "Orden no encontrada")
	}
	return order, nil
}

func (os Order) UpdateOrderStatus(id uuid.UUID, status checkout.OrderSyncStatus) error {
	sql := `UPDATE store_orders SET sync_status = $1 WHERE id = $2`
	c, err := os.db.Exec(context.Background(), sql, status, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if c.RowsAffected() != 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Orden no encontrada")
	}
	return nil
}

// Order products management

func (os Order) GetOrderProducts(order checkout.Order) ([]checkout.OrderProduct, error) {
	rows, err := os.db.Query(context.Background(), `SELECT id, quantity, details, product_type, product_category, product_item,
product_name, product_price, product_details, status, updated_at, product_id, product_part_number, product_currency
FROM order_products
WHERE order_id = $1`, order.Id)
	if err != nil {
		return []checkout.OrderProduct{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (checkout.OrderProduct, error) {
		var product checkout.OrderProduct
		var productID *int

		product.Details = make(map[string]string)
		product.ProductDetails = make(map[string]string)

		var hstoreDetails pgtype.Hstore
		var hstoreProductDetails pgtype.Hstore

		err := row.Scan(&product.Id, &product.Quantity, &hstoreDetails, &product.ProductType, &product.ProductCategory, &product.ProductItem,
			&product.ProductName, &product.ProductPrice, &hstoreProductDetails, &product.Status, &product.UpdatedAt, &productID, &product.ProductPartNumber,
			&product.ProductCurrency)
		product.Order = order
		for key, value := range hstoreDetails {
			if value != nil {
				product.Details[key] = *value
			}
		}
		for key, value := range hstoreProductDetails {
			if value != nil {
				product.ProductDetails[key] = *value
			}
		}

		// Query and attach product
		sql := `SELECT id, name, price, slug, stock, currency
		FROM store_products
		WHERE id = $1`
		var storeProduct store.Product
		if err := os.db.QueryRow(context.Background(), sql, productID).Scan(&storeProduct.Id,
			&storeProduct.Name, &storeProduct.Price, &storeProduct.Slug, &storeProduct.Stock, &storeProduct.Currency); err == nil {
			product.Product = storeProduct
		}

		return product, err
	})
	if err != nil {
		return []checkout.OrderProduct{}, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return products, nil
}

func (os Order) InsertOrderProducts(order checkout.Order, products []checkout.OrderProduct) (uuid.UUID, error) {
	tx, err := os.db.Begin(context.Background())
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer tx.Rollback(context.Background())

	// Insert order
	var orderID uuid.UUID
	if err := tx.QueryRow(context.Background(), `INSERT INTO store_orders (email, phone_number, name, address, city, postal_code)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`, order.Email, order.Phone, order.Name, order.Address, order.City, order.PostalCode).Scan(&orderID); err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError, "Error inserting new order")
	}

	for _, product := range products {
		if product.Product.Id == 0 {
			return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Insert product
		hstoreDetails := make(pgtype.Hstore, len(product.Details))
		for key, val := range product.Details {
			valCopy := val
			hstoreDetails[key] = &valCopy
		}

		hstoreProductDetails := make(pgtype.Hstore, len(product.ProductDetails))
		for key, val := range product.ProductDetails {
			valCopy := val
			hstoreProductDetails[key] = &valCopy
		}

		if _, err := tx.Exec(context.Background(), `INSERT INTO order_products (order_id, quantity, details,
product_type, product_category, product_item, product_name, product_price, product_details, product_id, product_part_number, product_currency)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`, orderID, product.Quantity, hstoreDetails,
			product.ProductType, product.ProductCategory, product.ProductItem, product.ProductName, product.ProductPrice,
			hstoreProductDetails, product.Product.Id, product.ProductPartNumber, product.ProductCurrency); err != nil {
			return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return orderID, nil
}

func (os Order) UpdateProductsStock(order checkout.Order, products []checkout.OrderProduct) error {
	// Start transaction
	tx, err := os.db.Begin(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer tx.Rollback(context.Background())

	// Update stock
	for _, product := range products {
		if product.Product.Id == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "Producto no encontrado")
		}
		if product.Product.Stock != nil {
			sql := `UPDATE store_products SET stock = stock - $1 WHERE id = $2 AND stock - $1 >= 0`
			c, err := os.db.Exec(context.Background(), sql, product.Quantity, product.Product.Id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
			if c.RowsAffected() != 1 {
				return echo.NewHTTPError(http.StatusNotFound, "Stock inválido")
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(context.Background()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}
