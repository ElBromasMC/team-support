package handler

import (
	"alc/model"
	"alc/view/component"
	"alc/view/user"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleSignupShow(c echo.Context) error {
	return render(c, http.StatusOK, user.SignupShow())
}

func (h *Handler) HandleSignup(c echo.Context) (err error) {
	u := &model.User{}
	if err = c.Bind(u); err != nil {
		return
	}
	return render(c, http.StatusOK, component.Message(fmt.Sprintf("Usuario: %s, Key: %s", u.Email, u.Password)))
}
