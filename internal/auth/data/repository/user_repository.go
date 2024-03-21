package repository

import (
	"context"
	"github.com/kurneo/go-template/internal/auth/data/datasource"
	"github.com/kurneo/go-template/internal/auth/domain/entity"
	"github.com/kurneo/go-template/internal/auth/domain/repository"
	"github.com/kurneo/go-template/pkg/error"
	"time"
)

type UserRepo struct {
	u *datasource.UserDatasource
}

func (repo UserRepo) GetUser(ctx context.Context, email string) (*entity.User, error.Contract) {
	u, err := repo.u.GetUser(ctx, email)
	if err != nil {
		return nil, error.NewDatasource(err)
	}
	return u, nil
}

func (repo UserRepo) GetUserById(ctx context.Context, id int64) (*entity.User, error.Contract) {
	u, err := repo.u.GetUserById(ctx, id)
	if err != nil {
		return nil, error.NewDatasource(err)
	}
	return u, nil
}

func (repo UserRepo) UpdateLastLoginTime(ctx context.Context, user *entity.User, time time.Time) error.Contract {
	err := repo.u.UpdateLastLoginTime(ctx, user, time)
	if err != nil {
		return error.NewDatasource(err)
	}
	return nil
}

func NewUserRepo(userDatasource *datasource.UserDatasource) repository.UserRepositoryContact {
	return &UserRepo{
		u: userDatasource,
	}
}
