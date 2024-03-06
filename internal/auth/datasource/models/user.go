package models

import (
	"github.com/kurneo/go-template/internal/auth/entities"
	"github.com/kurneo/go-template/pkg/support/repository"
	"time"
)

type User[T repository.PrimaryKey] struct {
	ID          T          `gorm:"primaryKey"`
	Name        string     `gorm:"column:name"`
	Email       string     `gorm:"column:email"`
	Password    string     `gorm:"column:password"`
	LastLoginAt *time.Time `gorm:"column:last_login_at"`
	CreatedAt   *time.Time `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (a User[T]) TableName() string {
	return "users"
}

func (a User[T]) ToEntity() *entities.User[T] {
	return &entities.User[T]{
		ID:          a.ID,
		Name:        a.Name,
		Email:       a.Email,
		Password:    a.Password,
		LastLoginAt: a.LastLoginAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func (a User[T]) FromEntity(e entities.User[T]) interface{} {
	return &User[T]{
		ID:          e.ID,
		Name:        e.Name,
		Email:       e.Email,
		Password:    e.Password,
		LastLoginAt: e.LastLoginAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
