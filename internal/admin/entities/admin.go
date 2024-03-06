package entities

import (
	"time"
)

type (
	Admin struct {
		ID          int        `json:"-"`
		Name        string     `json:"name"`
		Email       string     `json:"email"`
		Password    string     `json:"-"`
		LastLoginAt *time.Time `json:"last_login_at"`
		CreatedAt   *time.Time `json:"created_at"`
		UpdatedAt   *time.Time `json:"updated_at"`
	}

	PasswordCheckerContract interface {
		Check(h, p string) bool
	}
)

func (a Admin) ToMap() map[string]interface{} {
	user := map[string]interface{}{
		"name":          a.Name,
		"email":         a.Email,
		"last_login_at": a.LastLoginAt,
	}
	return user
}

func (a Admin) CheckPassword(password string, checker PasswordCheckerContract) bool {
	return checker.Check(a.Password, password)
}
