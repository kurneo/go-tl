package entity

import (
	"time"
)

type User struct {
	ID          int64      `json:"-"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Password    string     `json:"-"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (a User) ToMap() map[string]interface{} {
	user := map[string]interface{}{
		"name":          a.Name,
		"email":         a.Email,
		"last_login_at": a.LastLoginAt,
	}
	return user
}
