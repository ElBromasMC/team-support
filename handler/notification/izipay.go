package notification

import (
	"context"
	"fmt"
	"strings"

	"alc/config"
	"alc/model/checkout"
	"alc/model/payment"
	"alc/model/transaction"
	view "alc/view/checkout"

	"net/http"

	"github.com/a-h/templ"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"github.com/wneessen/go-mail"
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

	// TODO: logs
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
		transDate := form.Get("vads_trans_date")
		transUuid := form.Get("vads_trans_uuid")
		// captureDelay, err := strconv.Atoi(form.Get("vads_capture_delay"))
		// if err != nil {
		// 	return c.NoContent(http.StatusOK)
		// }

		switch transStatus {
		case "CAPTURED", "AUTHORISED":
			// Set transaction status to 'COMPLETED' or 'AUTHORISED'
			var status transaction.TransactionStatus
			if transStatus == "CAPTURED" {
				status = transaction.Completed
			} else {
				status = transaction.Authorised
			}
			err := h.TransactionService.UpdateTransaction(orderId, transId, transUuid, status)
			if err != nil {
				return c.NoContent(http.StatusOK)
			}

			// Query order and products
			order, err := h.OrderService.GetOrderById(orderId)
			if err != nil {
				return c.NoContent(http.StatusOK)
			}
			products, err := h.OrderService.GetOrderProducts(order)
			if err != nil {
				return c.NoContent(http.StatusOK)
			}

			// Get exchange rate
			rate, err := h.CurrencyService.GetExchangeRate(config.STORE_CURRENCY)
			if err != nil {
				// TODO
			}

			// Proceed only if the sync status is 'PENDING'
			if order.SyncStatus != checkout.Pending {
				return c.NoContent(http.StatusOK)
			}

			// Sync stock
			if err := h.OrderService.UpdateProductsStock(order, products); err != nil {
				if err := h.OrderService.UpdateOrderStatus(order.Id, checkout.Failed); err != nil {
					return c.NoContent(http.StatusOK)
				}
				// Mail to client
				msg1 := mail.NewMsg()
				if err := msg1.From(h.EmailService.GetSenderEmail()); err != nil {
					return c.NoContent(http.StatusOK)
				}
				if err := msg1.To(order.Email); err != nil {
					return c.NoContent(http.StatusOK)
				}
				msg1.Subject(fmt.Sprintf("Error al procesar la orden %d", order.PurchaseOrder))
				body1, err := templ.ToGoHTML(context.Background(), view.ClientErrorNotification(order, products, h.EmailService.GetWebHostname(), rate))
				if err != nil {
					return c.NoContent(http.StatusOK)
				}
				msg1.SetBodyString(mail.TypeTextHTML, string(body1))
				// Mail to company
				msg2 := mail.NewMsg()
				if err := msg2.From(h.EmailService.GetSenderEmail()); err != nil {
					return c.NoContent(http.StatusOK)
				}
				if err := msg2.To(h.EmailService.GetCompanyEmail()); err != nil {
					return c.NoContent(http.StatusOK)
				}
				msg2.Subject(fmt.Sprintf("Error al procesar la orden %d", order.PurchaseOrder))
				body2, err := templ.ToGoHTML(context.Background(), view.CompanyErrorNotification(order, products, h.EmailService.GetWebHostname(),
					transUuid, transDate, rate))
				if err != nil {
					return c.NoContent(http.StatusOK)
				}
				msg2.SetBodyString(mail.TypeTextHTML, string(body2))
				// Send warning notifications to the customer and the store
				if err := h.EmailService.DialAndSend(context.Background(), msg1, msg2); err != nil {
					return c.NoContent(http.StatusOK)
				}
			} else {
				if err := h.OrderService.UpdateOrderStatus(order.Id, checkout.Completed); err != nil {
					return c.NoContent(http.StatusOK)
				}
				// Mail to client
				msg1 := mail.NewMsg()
				if err := msg1.From(h.EmailService.GetSenderEmail()); err != nil {
					return c.NoContent(http.StatusOK)
				}
				if err := msg1.To(order.Email); err != nil {
					return c.NoContent(http.StatusOK)
				}
				msg1.Subject(fmt.Sprintf("Confirmación de la orden %d", order.PurchaseOrder))
				body1, err := templ.ToGoHTML(context.Background(), view.ClientSuccessNotification(order, products, h.EmailService.GetWebHostname(), rate))
				if err != nil {
					return c.NoContent(http.StatusOK)
				}
				msg1.SetBodyString(mail.TypeTextHTML, string(body1))
				// Mail to company
				msg2 := mail.NewMsg()
				if err := msg2.From(h.EmailService.GetSenderEmail()); err != nil {
					return c.NoContent(http.StatusOK)
				}
				if err := msg2.To(h.EmailService.GetCompanyEmail()); err != nil {
					return c.NoContent(http.StatusOK)
				}
				msg2.Subject(fmt.Sprintf("Orden %d procesada exitosamente", order.PurchaseOrder))
				body2, err := templ.ToGoHTML(context.Background(), view.CompanySuccessNotification(order, products, h.EmailService.GetWebHostname(),
					transUuid, transDate, rate))
				if err != nil {
					return c.NoContent(http.StatusOK)
				}
				msg2.SetBodyString(mail.TypeTextHTML, string(body2))
				// Send successful notifications to the customer and the store
				if err := h.EmailService.DialAndSend(context.Background(), msg1, msg2); err != nil {
					return c.NoContent(http.StatusOK)
				}
			}

		case "REFUSED", "EXPIRED":
			// Set transaction status to 'FAILED'
			err := h.TransactionService.UpdateTransaction(orderId, transId, transUuid, transaction.Failed)
			if err != nil {
				return c.NoContent(http.StatusOK)
			}
		case "CANCELLED":
			// Set transaction status to 'CANCELLED'
			err := h.TransactionService.UpdateTransaction(orderId, transId, transUuid, transaction.Cancelled)
			if err != nil {
				return c.NoContent(http.StatusOK)
			}
		}
	}

	return c.NoContent(http.StatusOK)
}
