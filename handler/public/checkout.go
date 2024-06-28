package public

import (
	"alc/handler/util"
	"alc/model/cart"
	"alc/model/checkout"
	"alc/model/payment"
	"alc/model/transaction"
	view "alc/view/checkout"
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"github.com/wneessen/go-mail"
)

// GET "/checkout?msg"
func (h *Handler) HandleCheckoutShow(c echo.Context) error {
	// Parsing request
	msg := c.QueryParam("msg")

	// Get cart items
	items := cart.GetItems(c.Request().Context())
	if len(items) == 0 {
		return c.Redirect(http.StatusFound, "/store")
	}
	return util.Render(c, http.StatusOK, view.Show(items, msg))
}

// POST "/checkout/orders"
func (h *Handler) HandleCheckoutOrderInsertion(c echo.Context) error {
	// Parsing request
	var order checkout.Order
	order.Email = c.FormValue("email")
	order.Phone = c.FormValue("phone")
	order.Name = c.FormValue("billing-name")
	order.Address = c.FormValue("billing-address")
	order.City = c.FormValue("billing-city")
	order.PostalCode = c.FormValue("billing-zip")

	// Validate order
	order, err := order.Normalize()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get cart items
	items := cart.GetItems(c.Request().Context())
	if len(items) == 0 {
		return c.Redirect(http.StatusFound, "/store")
	}

	// Get order products
	products := make([]checkout.OrderProduct, 0, len(items))
	for _, i := range items {
		p := checkout.OrderProduct{
			Quantity:        i.Quantity,
			Details:         i.Details,
			Product:         i.Product,
			ProductType:     i.Product.Item.Category.Type,
			ProductCategory: i.Product.Item.Category.Name,
			ProductItem:     i.Product.Item.Name,
			ProductName:     i.Product.Name,
			ProductPrice:    i.Product.Price,
			ProductDetails:  i.Product.Details,
		}
		products = append(products, p)
	}

	// Insert order products
	orderID, err := h.OrderService.InsertOrderProducts(order, products)
	if err != nil {
		return err
	}
	order.Id = orderID

	// Remove cart cookie
	// sess, _ := session.Get(cart.SessionName, c)
	// sess.Options = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   -1,
	// 	Secure:   true,
	// 	HttpOnly: false,
	// 	SameSite: http.SameSiteStrictMode,
	// }
	// if err := sess.Save(c.Request(), c.Response()); err != nil {
	// 	c.Logger().Debug("Error removing cart session: ", err)
	// 	return echo.NewHTTPError(http.StatusInternalServerError)
	// }

	// Send confirmation email to client
	msg := mail.NewMsg()
	if err := msg.From("no-reply@teamsupportperu.com"); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if err := msg.To(order.Email); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Email no válido")
	}
	msg.Subject("Order Payment Test")

	body, err := templ.ToGoHTML(context.Background(), view.ConfirmationEmail(order, products, "localhost:8080"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	msg.SetBodyString(mail.TypeTextHTML, string(body))

	if err := h.EmailService.DialAndSend(context.Background(), msg); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusSeeOther, "/checkout/confirmation")
}

// GET "/checkout/orders/:orderID/payment?fail"
func (h *Handler) HandleCheckoutPaymentShow(c echo.Context) error {
	// Parsing request
	orderID, err := uuid.FromString(c.Param("orderID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Identificador no válido")
	}

	fail := false
	if len(c.QueryParam("fail")) > 0 {
		fail = true
	}

	// Query order and products
	order, err := h.OrderService.GetOrderById(orderID)
	if err != nil {
		return err
	}

	products, err := h.OrderService.GetOrderProducts(order)
	if err != nil {
		return err
	}

	// Create transaction if not exists
	trans, err := h.TransactionService.GetTransaction(order)
	if err != nil {
		// Check lock time
		if time.Now().After(order.LockedAt) {
			msg := "Su orden ha expirado"
			return c.Redirect(http.StatusFound, "/checkout?msg="+url.QueryEscape(msg))
		}

		// Check product availability

		// Create and attach transaction
		trans, err = h.TransactionService.InsertTransaction(order, checkout.CalculateAmount(products), "IZIPAY")
		if err != nil {
			return err
		}
	} else {
		// Check transaction status
		switch trans.Status {
		case transaction.Authorised, transaction.Completed:
			return c.Redirect(http.StatusFound, "/checkout/orders/"+order.Id.String())
		case transaction.Failed, transaction.Cancelled:
			msg := "La transacción ha fallado"
			return c.Redirect(http.StatusFound, "/checkout?msg="+url.QueryEscape(msg))
		}

		// Check lock time
		if time.Now().After(order.LockedAt) {
			msg := "Su orden ha expirado"
			return c.Redirect(http.StatusFound, "/checkout?msg="+url.QueryEscape(msg))
		}

		// Check product availability
	}

	// Generate the form fields
	formFields := h.PaymentService.GetPaymentData(order, trans)

	return util.Render(c, http.StatusOK, view.PaymentPage(order, products, formFields, fail))
}

// POST "/checkout/orders/:orderID/preview"
func (h *Handler) HandleCheckoutOrderPreview(c echo.Context) error {
	// Parsing request
	orderID, err := uuid.FromString(c.Param("orderID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Identificador no válido")
	}

	form, err := c.FormParams()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Datos inválidos")
	}

	// Query order
	order, err := h.OrderService.GetOrderById(orderID)
	if err != nil {
		return err
	}

	// Compute and check signature
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

	if form.Get("signature") != signature {
		return echo.NewHTTPError(http.StatusForbidden, "Datos inválidos")
	}

	// Check orderID
	vadsOrderID, err := uuid.FromString(form.Get("vads_order_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Datos inválidos")
	}

	if vadsOrderID != order.Id {
		return echo.NewHTTPError(http.StatusBadRequest, "Datos inválidos")
	}

	// Collect data

	return c.NoContent(http.StatusOK)
}

// GET "/checkout/orders/:orderID"
func (h *Handler) HandleCheckoutOrderShow(c echo.Context) error {
	// Parsing request
	orderID, err := uuid.FromString(c.Param("orderID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Identificador no válido")
	}

	// Query order
	order, err := h.OrderService.GetOrderById(orderID)
	if err != nil {
		return err
	}

	// Create transaction if not exists
	trans, err := h.TransactionService.GetTransaction(order)
	if err != nil {
		return c.Redirect(http.StatusFound, "/checkout/orders/"+order.Id.String()+"/payment")
	}

	// Check transaction status
	switch trans.Status {
	case transaction.Pending, transaction.Failed, transaction.Cancelled:
		return c.Redirect(http.StatusFound, "/checkout/orders/"+order.Id.String()+"/payment")
	}

	// Query products
	products, err := h.OrderService.GetOrderProducts(order)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, view.Tracking(order, products))
}
