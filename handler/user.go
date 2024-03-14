package handler

import (
	"alc/model"
	"alc/view/user"
	"context"
	"net/http"
	"net/mail"
	"os"
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

	if !(8 <= len(u.Password) && len(u.Password) <= 30) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid password")
	}

	if !(0 < len(tname) && len(tname) <= 225) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid name")
	}

	// Hash password
	hpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	// Save user
	if _, err := h.DB.Exec(context.Background(), `INSERT INTO users (name, email, hashed_password)
VALUES ($1, $2, $3)`, tname, temail, string(hpass)); err != nil {
		// ToDo: Test and handle unique email condition
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusFound, "/login")
}

// Login
func (h *Handler) HandleLoginShow(c echo.Context) error {
	return render(c, http.StatusOK, user.LoginShow(c.QueryParam("to")))
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
	var userid uuid.UUID
	var hpass string
	if err := h.DB.QueryRow(context.Background(), `SELECT user_id, hashed_password
FROM users WHERE email = $1`, temail).Scan(&userid, &hpass); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Email not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hpass), []byte(u.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Wrong password")
	}

	// Start new session
	var session uuid.UUID
	if err := h.DB.QueryRow(context.Background(), `INSERT INTO sessions (user_id)
VALUES ($1) RETURNING session_id`, userid).Scan(&session); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Write a cookie
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = session.String()
	cookie.Expires = time.Now().AddDate(0, 1, 0)
	if os.Getenv("HTTPS") == "true" {
		cookie.Secure = true
	}
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)

	if to := c.QueryParam("to"); len(to) > 0 {
		return c.Redirect(http.StatusFound, to)
	}
	return c.Redirect(http.StatusFound, "/")
}

// Logout
func (h *Handler) HandleLogout(c echo.Context) error {
	defer RemoveCookie(c, "session")
	user, ok := c.Request().Context().Value("user").(model.User)
	if !ok {
		return c.Redirect(http.StatusFound, "/")
	}
	if _, err := h.DB.Exec(context.Background(), `DELETE FROM sessions
WHERE session_id = $1`, user.Session); err != nil {
		return c.Redirect(http.StatusFound, "/")
	}
	return c.Redirect(http.StatusFound, "/")
}
