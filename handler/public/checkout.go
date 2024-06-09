package public

import (
	"alc/handler/util"
	"alc/model/cart"
	"alc/model/checkout"
	view "alc/view/checkout"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// GET "/checkout"
func (h *Handler) HandleCheckoutShow(c echo.Context) error {
	// Get cart items
	items := cart.GetItems(c.Request().Context())
	if len(items) == 0 {
		return c.Redirect(http.StatusFound, "/store")
	}
	return util.Render(c, http.StatusOK, view.Show(items))
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

	return c.Redirect(http.StatusFound, "/checkout/orders/"+orderID.String()+"/payment")
}

// GET "/checkout/orders/:orderID/payment"
func (h *Handler) HandleCheckoutPaymentShow(c echo.Context) error {
	// Parsing request
	orderID, err := uuid.FromString(c.Param("orderID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Identificador no válido")
	}

	// Query data
	order, err := h.OrderService.GetOrderById(orderID)
	if err != nil {
		return err
	}
	products, err := h.OrderService.GetOrderProducts(order)
	if err != nil {
		return err
	}
	trans, err := h.TransactionService.InsertTransaction(order, checkout.CalculateAmount(products), "IZIPAY")
	if err != nil {
		return err
	}
	formFields := h.PaymentService.GetPaymentData(order, trans)

	return util.Render(c, http.StatusOK, view.PaymentPage(order, products, formFields))
}

// GET "/checkout/orders/:orderID"
func (h *Handler) HandleCheckoutOrderShow(c echo.Context) error {
	// Parsing request
	orderID, err := uuid.FromString(c.Param("orderID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Identificador no válido")
	}

	// Query data
	order, err := h.OrderService.GetOrderById(orderID)
	if err != nil {
		return err
	}
	products, err := h.OrderService.GetOrderProducts(order)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, view.Tracking(order, products))
}

// TODO
func TODOCHECK(c echo.Context) error {
	// Remove cart cookie
	sess, _ := session.Get(cart.SessionName, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		c.Logger().Debug("Error removing cart session: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
