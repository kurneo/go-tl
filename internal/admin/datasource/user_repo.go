package datasource

import (
	"context"
	"github.com/kurneo/go-template/internal/admin/datasource/models"
	"github.com/kurneo/go-template/internal/admin/entities"
	pkgError "github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/support/repository"
)

type UserRepo struct {
	repository.Repository[models.Admin, entities.Admin]
}

func (repo UserRepo) GetUser(ctx context.Context, email string) (*entities.Admin, pkgError.Contract) {
	u, err := repo.FirstBy(ctx, repository.Equal("email", email))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo UserRepo) GetUserById(ctx context.Context, id int) (*entities.Admin, pkgError.Contract) {
	u, err := repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo UserRepo) UpdateLastLoginTime(ctx context.Context, user *entities.Admin) pkgError.Contract {
	var m models.Admin
	model := m.FromEntity(*user).(*models.Admin)
	err := repo.D.GetDB(ctx).Table((*model).TableName()).
		Where("id = ?", model.ID).
		Update("last_login_at", model.LastLoginAt).Error
	if err != nil {
		return pkgError.NewDatasource(err)
	}
	return nil
}
