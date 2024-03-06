package models

import (
	"github.com/kurneo/go-template/internal/admin/entities"
	"time"
)

type (
	Admin struct {
		ID          int        `gorm:"primaryKey"`
		Name        string     `gorm:"column:name"`
		Email       string     `gorm:"column:email"`
		Password    string     `gorm:"column:password"`
		LastLoginAt *time.Time `gorm:"column:last_login_at"`
		CreatedAt   *time.Time `gorm:"column:created_at"`
		UpdatedAt   *time.Time `gorm:"column:updated_at"`
	}

	AdminAccessToken struct {
		ID          int        `gorm:"primaryKey"`
		AccessToken string     `gorm:"column:token"`
		ExpiredAt   time.Time  `gorm:"column:expired_at"`
		AdminID     int        `gorm:"column:admin_id"`
		CreatedAt   *time.Time `gorm:"column:created_at"`
		Admin       *Admin     `gorm:"foreignKey:AdminID"`
	}
)

func (a Admin) TableName() string {
	return "admins"
}

func (a Admin) ToEntity() *entities.Admin {
	return &entities.Admin{
		ID:          a.ID,
		Name:        a.Name,
		Email:       a.Email,
		Password:    a.Password,
		LastLoginAt: a.LastLoginAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func (a Admin) FromEntity(e entities.Admin) interface{} {
	return &Admin{
		ID:          e.ID,
		Name:        e.Name,
		Email:       e.Email,
		Password:    e.Password,
		LastLoginAt: e.LastLoginAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func (t AdminAccessToken) TableName() string {
	return "admin_access_tokens"
}

func (t AdminAccessToken) ToEntity() *entities.AdminAccessToken {
	var a *entities.Admin
	if t.Admin != nil {
		a = t.Admin.ToEntity()
	}
	return &entities.AdminAccessToken{
		ID:          t.ID,
		AccessToken: t.AccessToken,
		ExpiredAt:   t.ExpiredAt,
		AdminID:     t.AdminID,
		CreatedAt:   t.CreatedAt,
		Admin:       a,
	}
}

func (t AdminAccessToken) FromEntity(e entities.AdminAccessToken) interface{} {
	return &AdminAccessToken{
		ID:          e.ID,
		AccessToken: e.AccessToken,
		ExpiredAt:   e.ExpiredAt,
		AdminID:     e.AdminID,
		CreatedAt:   e.CreatedAt,
	}
}
