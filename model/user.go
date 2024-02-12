package model

type User struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password,omitempty" form:"password"`
	Name     string `json:"name,omitempty" form:"name"`
}
