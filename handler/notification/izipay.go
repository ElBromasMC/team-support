package notification

import (
	"strings"

	"alc/model/payment"

	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIzipayNotification(c echo.Context) error {
	// Get form data
	form, err := c.FormParams()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// Check if 'vads_hash' is present
	if !form.Has("vads_hash") {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// Compute signature
	vadsFields := make([]payment.FormData, 0, len(form))
	for k, v := range form {
		if strings.HasPrefix(k, "vads_") && len(v) > 0 {
			vadsFields = append(vadsFields, payment.FormData{
				Key:   k,
				Value: v[0],
			})
		}
	}
	signature := h.PaymentService.ComputeSignature(vadsFields)

	// Compare signatures
	if form.Get("signature") != signature {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	switch form.Get("vads_url_check_src") {
	case "PAY":

	}

	return c.NoContent(http.StatusOK)
}
