package datasource

import (
	"github.com/kurneo/go-template/internal/admin/datasource/models"
	"github.com/kurneo/go-template/internal/admin/entities"
	"github.com/kurneo/go-template/internal/admin/usecase"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/repository"
)

func NewUserRepo(lg logger.Contract, db database.Contract) usecase.AdminRepositoryContract {
	return &UserRepo{
		Repository: repository.Repository[models.Admin, entities.Admin]{
			D: db,
			L: lg,
		},
	}
}

func NewTokenRepo(lg logger.Contract, db database.Contract) usecase.AdminTokenRepositoryContract {
	return &TokenRepo{
		Repository: repository.Repository[models.AdminAccessToken, entities.AdminAccessToken]{
			D: db,
			L: lg,
		},
	}
}
