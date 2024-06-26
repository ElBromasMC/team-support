package notification

import (
	"strings"

	"alc/model/payment"
	"alc/model/transaction"

	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIzipayPayNotification(c echo.Context) error {
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

	// Check context mode
	if form.Get("vads_ctx_mode") != string(h.PaymentService.GetMode()) {
		return c.NoContent(http.StatusOK)
	}

	switch form.Get("vads_url_check_src") {
	case "PAY", "BATCH_AUTO":
		// Get order uuid
		orderId, err := uuid.FromString(form.Get("vads_order_id"))
		if err != nil {
			return c.NoContent(http.StatusOK)
		}

		// Get transaction data
		transStatus := form.Get("vads_trans_status")
		transId := strings.ToLower(form.Get("vads_trans_id"))
		// transDate := form.Get("vads_trans_date")
		transUuid := form.Get("vads_trans_uuid")
		// captureDelay, err := strconv.Atoi(form.Get("vads_capture_delay"))
		// if err != nil {
		// 	return c.NoContent(http.StatusOK)
		// }

		// TODO: logs
		switch transStatus {
		case "AUTHORISED", "CAPTURED":
			// Set transaction status to 'COMPLETED'
			h.TransactionService.UpdateTransaction(orderId, transId, transUuid, transaction.Completed)

			// Send notifications to the customer and the store

		case "REFUSED", "EXPIRED":
			// Set transaction status to 'FAILED'
			h.TransactionService.UpdateTransaction(orderId, transId, transUuid, transaction.Failed)

		case "ABANDONED", "CANCELLED":
			// Set transaction status to 'CANCELLED'
			h.TransactionService.UpdateTransaction(orderId, transId, transUuid, transaction.Cancelled)
		}
	}

	return c.NoContent(http.StatusOK)
}
