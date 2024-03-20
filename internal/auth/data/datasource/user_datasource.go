package datasource

import (
	"context"
	"github.com/kurneo/go-template/internal/auth/data/model"
	"github.com/kurneo/go-template/internal/auth/domain/entity"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/support/repository"
)

type UserDatasource struct {
	repository.Repository[model.User, entity.User, int64]
}

func (repo UserDatasource) GetUser(ctx context.Context, email string) (*entity.User, error) {
	u, err := repo.FirstBy(ctx, repository.Equal("email", email))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo UserDatasource) GetUserById(ctx context.Context, id int64) (*entity.User, error) {
	u, err := repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo UserDatasource) UpdateLastLoginTime(ctx context.Context, user *entity.User) error {
	var m model.User
	model := m.FromEntity(*user).(*model.User)
	err := repo.D.GetDB(ctx).Table((*model).TableName()).
		Where("id = ?", model.ID).
		Update("last_login_at", model.LastLoginAt).Error
	if err != nil {
		return err
	}
	return nil
}

func NewUserDataSource(db database.Contract) *UserDatasource {
	return &UserDatasource{
		repository.Repository[model.User, entity.User, int64]{
			D: db,
		},
	}
}