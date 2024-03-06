package datasource

import (
	"context"
	"github.com/kurneo/go-template/internal/auth/datasource/models"
	"github.com/kurneo/go-template/internal/auth/entities"
	pkgError "github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/support/repository"
)

type UserRepo[T repository.PrimaryKey] struct {
	repository.Repository[models.User[T], entities.User[T], T]
}

func (repo UserRepo[T]) GetUser(ctx context.Context, email string) (*entities.User[T], pkgError.Contract) {
	u, err := repo.FirstBy(ctx, repository.Equal("email", email))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo UserRepo[T]) GetUserById(ctx context.Context, id T) (*entities.User[T], pkgError.Contract) {
	u, err := repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo UserRepo[T]) UpdateLastLoginTime(ctx context.Context, user *entities.User[T]) pkgError.Contract {
	var m models.User[T]
	model := m.FromEntity(*user).(*models.User[T])
	err := repo.D.GetDB(ctx).Table((*model).TableName()).
		Where("id = ?", model.ID).
		Update("last_login_at", model.LastLoginAt).Error
	if err != nil {
		return pkgError.NewDatasource(err)
	}
	return nil
}
