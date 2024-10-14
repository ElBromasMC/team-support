package service

import (
	"alc/model/store"
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Device struct {
	db *pgxpool.Pool
}

func NewDeviceService(db *pgxpool.Pool) Device {
	return Device{
		db: db,
	}
}

// Device management

func (ds Device) GetDevices(valid bool) ([]store.Device, error) {
	sql := `SELECT id, serie, created_at, updated_at, is_before_six_months, is_after_six_months
	FROM store_devices
	WHERE valid = $1`
	rows, err := ds.db.Query(context.Background(), sql, valid)
	if err != nil {
		return []store.Device{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	devices, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Device, error) {
		var device store.Device
		if err := row.Scan(&device.Id, &device.Serie, &device.CreatedAt, &device.UpdatedAt,
			&device.IsBeforeSixMonths, &device.IsAfterSixMonths); err != nil {
			return store.Device{}, err
		}
		device.Valid = valid
		return device, nil
	})
	if err != nil {
		return []store.Device{}, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return devices, nil
}

func (ds Device) GetDevice(serial string) (store.Device, error) {
	var device store.Device
	sql := `SELECT id, serie, created_at, updated_at, is_before_six_months, is_after_six_months, valid
	FROM store_devices
	WHERE serie = $1`
	if err := ds.db.QueryRow(context.Background(), sql, serial).Scan(&device.Id, &device.Serie, &device.CreatedAt,
		&device.UpdatedAt, &device.IsBeforeSixMonths, &device.IsAfterSixMonths, &device.Valid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return store.Device{}, echo.NewHTTPError(http.StatusNotFound, "La serie no se encuentra registrada")
		} else {
			return store.Device{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return device, nil
}

func (ds Device) GetDeviceById(id int) (store.Device, error) {
	var device store.Device
	sql := `SELECT id, serie, created_at, updated_at, is_before_six_months, is_after_six_months, valid
	FROM store_devices
	WHERE id = $1`
	if err := ds.db.QueryRow(context.Background(), sql, id).Scan(&device.Id, &device.Serie, &device.CreatedAt,
		&device.UpdatedAt, &device.IsBeforeSixMonths, &device.IsAfterSixMonths, &device.Valid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return store.Device{}, echo.NewHTTPError(http.StatusNotFound, "La serie no se encuentra registrada")
		} else {
			return store.Device{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return device, nil
}

func (ds Device) InsertDevice(email string, device store.Device) error {
	// Check if exists
	var exists bool
	sql := `SELECT EXISTS (
		SELECT 1 FROM store_devices
		WHERE serie = $1
	)`
	if err := ds.db.QueryRow(context.Background(), sql, device.Serie).Scan(&exists); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Insert or change it
	var device_id int
	if !exists {
		sql := `INSERT INTO store_devices (serie, valid, is_before_six_months, is_after_six_months)
		VALUES ($1, $2, $3, $4) RETURNING id`
		if err := ds.db.QueryRow(context.Background(), sql, device.Serie, device.Valid,
			device.IsBeforeSixMonths, device.IsAfterSixMonths).Scan(&device_id); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		sql := `UPDATE store_devices
		SET valid = $1, is_before_six_months = $2, is_after_six_months = $3
		WHERE serie = $4 RETURNING id`
		if err := ds.db.QueryRow(context.Background(), sql, device.Valid,
			device.IsBeforeSixMonths, device.IsAfterSixMonths, device.Serie).Scan(&device_id); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	// Append history
	sql1 := `INSERT INTO store_devices_history (device_id, issued_by) VALUES ($1, $2)`
	if _, err := ds.db.Exec(context.Background(), sql1, device_id, email); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

func (ds Device) DesactivateDevice(device store.Device) error {
	sql := `UPDATE store_devices SET valid = FALSE WHERE id = $1`
	c, err := ds.db.Exec(context.Background(), sql, device.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if c.RowsAffected() < 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Dispositivo no encontrado")
	}
	return nil
}

func sameDate(t1, t2 time.Time, loc *time.Location) bool {
	t1InLoc := t1.In(loc)
	t2InLoc := t2.In(loc)

	return t1InLoc.Year() == t2InLoc.Year() &&
		t1InLoc.Month() == t2InLoc.Month() &&
		t1InLoc.Day() == t2InLoc.Day()
}

func (ds Device) CheckSerialValidity(product store.Product, serial string) error {
	dev, err := ds.GetDevice(serial)
	if err != nil {
		return err
	}
	if !dev.Valid {
		return echo.NewHTTPError(http.StatusBadRequest, "La serie no se encuentra registrada")
	}
	loc, _ := time.LoadLocation("America/Lima")
	if !sameDate(dev.UpdatedAt, time.Now(), loc) {
		return echo.NewHTTPError(http.StatusBadRequest, "La serie registrada ha expirado")
	}
	if !((dev.IsAfterSixMonths && product.AcceptAfterSixMonths) ||
		(dev.IsBeforeSixMonths && product.AcceptBeforeSixMonths)) {
		return echo.NewHTTPError(http.StatusBadRequest, "La serie no aplica para esta garantía")
	}
	return nil
}

// Device history management

func (ds Device) GetDeviceHistory(device store.Device) ([]store.DeviceHistory, error) {
	sql := `SELECT id, issued_by, issued_at
	FROM store_devices_history
	WHERE device_id = $1
	ORDER BY issued_at DESC`
	rows, err := ds.db.Query(context.Background(), sql, device.Id)
	if err != nil {
		return []store.DeviceHistory{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	history, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.DeviceHistory, error) {
		var h store.DeviceHistory
		if err := row.Scan(&h.Id, &h.IssuedBy, &h.IssuedAt); err != nil {
			return store.DeviceHistory{}, err
		}
		h.Device = device
		return h, nil
	})
	if err != nil {
		return []store.DeviceHistory{}, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return history, nil
}

// Device data management

func (ds Device) GetDeviceData(serial string) (dd store.DeviceData, err error) {
	sql := `
	SELECT
		id,
		product_serial,
		product_type,
		part_no_model,
		warranty_start_date,
		warranty_end_date
	FROM
		store_devices_data
	WHERE
		product_serial = $1
	`
	err = ds.db.QueryRow(context.Background(), sql, serial).Scan(
		&dd.Id,
		&dd.ProductSerial,
		&dd.ProductType,
		&dd.PartNoModel,
		&dd.WarrantyStartDate,
		&dd.WarrantyEndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = echo.NewHTTPError(http.StatusNotFound, "La serie no se encuentra registrada")
			return
		} else {
			err = echo.NewHTTPError(http.StatusInternalServerError)
			return
		}
	}
	return
}

func (ds Device) InsertDeviceData(data []store.DeviceData) (int, error) {
	sql := `TRUNCATE TABLE temp_store_devices_data`
	_, err := ds.db.Exec(context.Background(), sql)
	if err != nil {
		return 0, err
	}

	// Copy to temp table
	_, err = ds.db.CopyFrom(
		context.Background(),
		pgx.Identifier{"temp_store_devices_data"},
		[]string{
			"product_serial",
			"product_type",
			"part_no_model",
			"warranty_start_date",
			"warranty_end_date",
		},
		pgx.CopyFromSlice(len(data), func(i int) ([]any, error) {
			return []any{
				data[i].ProductSerial,
				data[i].ProductType,
				data[i].PartNoModel,
				data[i].WarrantyStartDate,
				data[i].WarrantyEndDate,
			}, nil
		}),
	)
	if err != nil {
		return 0, err
	}

	// Merge with final table
	sql1 := `
	INSERT INTO store_devices_data (product_serial, product_type, part_no_model, warranty_start_date, warranty_end_date)
	SELECT DISTINCT ON (product_serial)
		product_serial, product_type, part_no_model, warranty_start_date, warranty_end_date
	FROM temp_store_devices_data
	ORDER BY product_serial, warranty_end_date DESC
	ON CONFLICT (product_serial) DO UPDATE SET
		product_type = EXCLUDED.product_type,
		part_no_model = EXCLUDED.part_no_model,
		warranty_start_date = EXCLUDED.warranty_start_date,
		warranty_end_date = EXCLUDED.warranty_end_date
	`
	c, err := ds.db.Exec(context.Background(), sql1)
	if err != nil {
		return 0, err
	}

	return int(c.RowsAffected()), nil
}

func diffMinRight(s string) (int, error) {
	parts := strings.Split(s, "_")
	if len(parts) != 2 {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Producto inválido")
	}

	rightNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Producto inválido")
	}

	leftNumsStr := strings.Split(parts[0], "-")
	if len(leftNumsStr) == 0 {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Producto inválido")
	}

	minNum := math.MaxInt
	for _, numStr := range leftNumsStr {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, echo.NewHTTPError(http.StatusBadRequest, "Producto inválido")
		}
		if num < minNum {
			minNum = num
		}
	}
	if rightNum-minNum <= 0 {
		return rightNum, nil
	} else {
		return rightNum - minNum, nil
	}
}

func containsSubstrings(input string) bool {
	substrings := []string{"accidental", "Disco", "Domicilio"}
	for _, substr := range substrings {
		if strings.Contains(input, substr) {
			return true
		}
	}
	return false
}

func (ds Device) CheckSerialValidityByDeviceData(product store.Product, serial string) (err error) {
	dd, err := ds.GetDeviceData(serial)
	if err != nil {
		return
	}
	if !strings.EqualFold(product.Item.Category.Slug, dd.ProductType) {
		err = echo.NewHTTPError(http.StatusBadRequest, "La serie no aplica para esta garantía")
		return
	}
	days_contracted := int(dd.WarrantyEndDate.Sub(dd.WarrantyStartDate).Hours() / 24)
	days_interval := int(time.Since(dd.WarrantyStartDate).Hours() / 24)
	days_diff := days_contracted - days_interval
	months_contracted := int(math.Round(float64(days_contracted)/360)) * 12
	std_warranty_period, err := diffMinRight(product.Slug)
	if err != nil {
		return
	}

	// Validation
	if days_diff < 0 {
		err = echo.NewHTTPError(http.StatusBadRequest, "Su periodo de garantía ha expirado")
		return
	}
	if months_contracted != std_warranty_period {
		err = echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Debe tener %d meses de garantía contratada", std_warranty_period))
		return
	}

	if len(product.Name) == 0 {
		err = echo.NewHTTPError(http.StatusInternalServerError)
		return
	}

	if !containsSubstrings(product.Name) {
		return nil
	} else if days_interval <= 180 {
		return nil
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "Está fuera del periodo de 180 días")
	}
}
