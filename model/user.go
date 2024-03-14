package model

import "github.com/gofrs/uuid/v5"

type UserRole string

const (
	AdminRole  UserRole = "ADMIN"
	NormalRole UserRole = "NORMAL"
)

type User struct {
	Name     string    `form:"name"`
	Email    string    `form:"email"`
	Password string    `form:"password"`
	Role     UserRole  `form:"-"`
	Session  uuid.UUID `form:"-"`
}
