package repository

import (
	"context"
	"github.com/kurneo/go-template/internal/auth/domain/entity"
	"github.com/kurneo/go-template/pkg/error"
	"time"
)

type UserRepositoryContact interface {
	GetUser(ctx context.Context, e string) (*entity.User, error.Contract)
	GetUserById(ctx context.Context, id int64) (*entity.User, error.Contract)
	UpdateLastLoginTime(ctx context.Context, u *entity.User, time time.Time) error.Contract
}
