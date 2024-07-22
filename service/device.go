package service

import (
	"alc/model/store"
	"context"
	"errors"
	"net/http"

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

func (ds Device) GetDeviceHistory(device store.Device) ([]store.DeviceHistory, error) {
	sql := `SELECT id, issued_by, issued_at
	FROM store_devices_history
	WHERE device_id = $1`
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
