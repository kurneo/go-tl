package datasource

import (
	"context"
	"github.com/kurneo/go-template/internal/auth/data/model"
	"github.com/kurneo/go-template/internal/auth/domain/entity"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/support/db_repository"
	"time"
)

type UserDatasource struct {
	db_repository.Repository[model.User, entity.User, int64]
}

func (repo UserDatasource) GetUser(ctx context.Context, email string) (*entity.User, error) {
	u, err := repo.FirstBy(ctx, db_repository.Param{Condition: db_repository.Equal("email", email)})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo UserDatasource) GetUserById(ctx context.Context, id int64) (*entity.User, error) {
	u, err := repo.FindByID(ctx, id, db_repository.Param{})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo UserDatasource) UpdateLastLoginTime(ctx context.Context, user *entity.User, time time.Time) error {
	var m model.User
	userModel := m.FromEntity(*user).(*model.User)
	err := repo.D.GetDB(ctx).Table((*userModel).TableName()).
		Where("id = ?", userModel.ID).
		Update("last_login_at", time).Error
	if err != nil {
		return err
	}
	user.LastLoginAt = &time
	return nil
}

func NewUserDataSource(db database.Contract) *UserDatasource {
	return &UserDatasource{
		db_repository.Repository[model.User, entity.User, int64]{
			D: db,
		},
	}
}
