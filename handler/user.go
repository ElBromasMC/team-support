package handler

import (
	"alc/model"
	"alc/view/user"
	"context"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Signup
func (h *Handler) HandleSignupShow(c echo.Context) error {
	return render(c, http.StatusOK, user.SignupShow())
}

func (h *Handler) HandleSignup(c echo.Context) error {
	// Bind
	var u model.User
	if err := c.Bind(&u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid format")
	}

	// Trim name and email
	temail := strings.TrimSpace(u.Email)
	tname := strings.TrimSpace(u.Name)

	// Validate
	if _, err := mail.ParseAddress(temail); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email")
	}

	if !(6 <= len(u.Password) && len(u.Password) <= 30) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid password")
	}

	if !(0 < len(tname) && len(tname) <= 225) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid name")
	}

	// Hash password
	var hpass []byte
	var err error
	if hpass, err = bcrypt.GenerateFromPassword([]byte(u.Password), 14); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	// Save user
	if _, err := h.DB.Exec(context.Background(), "INSERT INTO users (name, email, hashed_password, session) VALUES ($1, $2, $3, ARRAY[]::UUID[])", tname, temail, string(hpass)); err != nil {
		// ToDo: Test and handle unique email condition
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusFound, "/login")
}

// Login
func (h *Handler) HandleLoginShow(c echo.Context) error {
	return render(c, http.StatusOK, user.LoginShow())
}

func (h *Handler) HandleLogin(c echo.Context) error {
	// Bind
	var u model.User
	if err := c.Bind(&u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid format")
	}

	// Trim email
	temail := strings.TrimSpace(u.Email)

	// Verify credentials
	var hpass string
	if err := h.DB.QueryRow(context.Background(), "SELECT hashed_password FROM users WHERE email = $1", temail).Scan(&hpass); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Email not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hpass), []byte(u.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Wrong password")
	}

	// Start new session
	var session uuid.UUID
	if err := h.DB.QueryRow(context.Background(), "UPDATE users SET session = array_append(session, uuid_generate_v4()) WHERE email = $1 RETURNING session[array_length(session, 1)]", temail).Scan(&session); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Write a cookie
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = session.String()
	cookie.Expires = time.Now().Add(24 * time.Hour)
	// cookie.Secure = true
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}
