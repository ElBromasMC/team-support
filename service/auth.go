package service

import (
	"alc/model/auth"
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	db *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) Auth {
	return Auth{
		db: db,
	}
}

func (us Auth) GetUserByUuid(id uuid.UUID) (auth.User, error) {
	var u auth.User
	if err := us.db.QueryRow(context.Background(), `SELECT u.name, u.email, u.role
FROM users AS u
JOIN sessions AS s
ON u.user_id = s.user_id
WHERE s.session_id = $1`, id).Scan(&u.Name, &u.Email, &u.Role); err != nil {
		return auth.User{}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}
	u.Id = id
	return u, nil
}

func (us Auth) GetUserIdAndHpassByEmail(email string) (uuid.UUID, []byte, error) {
	var id uuid.UUID
	var hpass string

	if err := us.db.QueryRow(context.Background(), `SELECT user_id, hashed_password
FROM users WHERE email = $1`, email).Scan(&id, &hpass); err != nil {
		return uuid.UUID{}, []byte{}, echo.NewHTTPError(http.StatusUnauthorized, "Email no encontrado")
	}
	return id, []byte(hpass), nil
}

func (us Auth) InsertUser(u auth.User, hpass []byte) error {
	if _, err := us.db.Exec(context.Background(), `INSERT INTO users (name, email, hashed_password)
VALUES ($1, $2, $3)`, u.Name, u.Email, string(hpass)); err != nil {
		// TODO: Test and handle unique email condition
		return echo.NewHTTPError(http.StatusConflict, "Ya existe una cuenta con el email proporcionado")
	}
	return nil
}

func (us Auth) InsertSession(userId uuid.UUID) (uuid.UUID, error) {
	var session uuid.UUID
	if err := us.db.QueryRow(context.Background(), `INSERT INTO sessions (user_id)
VALUES ($1) RETURNING session_id`, userId).Scan(&session); err != nil {
		return uuid.UUID{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return session, nil
}

func (us Auth) DeleteSession(id uuid.UUID) error {
	if _, err := us.db.Exec(context.Background(), `DELETE FROM sessions
WHERE session_id = $1`, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
