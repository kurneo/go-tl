package model

import (
	"github.com/kurneo/go-template/internal/category/domain/entity"
	"time"
)

type Category struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	Description *string
	Status      int
	IsDefault   bool
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func (c Category) TableName() string {
	return "categories"
}

func (c Category) ToEntity() *entity.Category {
	return &entity.Category{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Status:      c.Status,
		IsDefault:   c.IsDefault,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func (c Category) FromEntity(e entity.Category) interface{} {
	return &Category{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Status:      e.Status,
		IsDefault:   e.IsDefault,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
