package model

import (
	"github.com/kurneo/go-template/internal/auth/domain/entity"
	"time"
)

type User struct {
	ID          int64      `gorm:"primaryKey"`
	Name        string     `gorm:"column:name"`
	Email       string     `gorm:"column:email"`
	Password    string     `gorm:"column:password"`
	LastLoginAt *time.Time `gorm:"column:last_login_at"`
	CreatedAt   *time.Time `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (a User) TableName() string {
	return "users"
}

func (a User) ToEntity() *entity.User {
	return &entity.User{
		ID:          a.ID,
		Name:        a.Name,
		Email:       a.Email,
		Password:    a.Password,
		LastLoginAt: a.LastLoginAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func (a User) FromEntity(e entity.User) interface{} {
	return &User{
		ID:          e.ID,
		Name:        e.Name,
		Email:       e.Email,
		Password:    e.Password,
		LastLoginAt: e.LastLoginAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
