package device

import (
	"alc/handler/util"
	"alc/model/auth"
	"alc/model/store"
	"alc/view/admin/device"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	devices, err := h.DeviceService.GetDevices(true)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, device.Show(devices))
}

func (h *Handler) HandleInsertion(c echo.Context) error {
	// Parsing request
	var dev store.Device
	dev.Serie = c.FormValue("serial")
	dev.Valid = true
	if c.FormValue("is-before") == "SI" {
		dev.IsBeforeSixMonths = true
		dev.IsAfterSixMonths = false
	} else {
		dev.IsBeforeSixMonths = false
		dev.IsAfterSixMonths = true
	}
	dev, err := dev.Normalize()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "La serie debe tener entre 12 y 15 caracteres")
	}

	// Query data
	user, ok := auth.GetUser(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Insert device
	h.DeviceService.InsertDevice(user.Email, dev)

	// Get updated devices
	devices, err := h.DeviceService.GetDevices(true)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, device.Table(devices))
}

func (h *Handler) HandleDeactivation(c echo.Context) error {
	return nil
}

func (h *Handler) HandleInsertionFormShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, device.InsertionForm())
}

func (h *Handler) HandleHistoryShow(c echo.Context) error {
	return nil
}

func (h *Handler) HandleDeactivationFormShow(c echo.Context) error {
	return nil
}
