package entities

import (
	"github.com/kurneo/go-template/pkg/support/repository"
	"time"
)

type (
	User[T repository.PrimaryKey] struct {
		ID          T          `json:"-"`
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

func (a User[T]) ToMap() map[string]interface{} {
	user := map[string]interface{}{
		"name":          a.Name,
		"email":         a.Email,
		"last_login_at": a.LastLoginAt,
	}
	return user
}

func (a User[T]) CheckPassword(password string, checker PasswordCheckerContract) bool {
	return checker.Check(a.Password, password)
}
