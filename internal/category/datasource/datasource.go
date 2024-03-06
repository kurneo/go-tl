package datasource

import (
	"github.com/kurneo/go-template/internal/category/datasource/models"
	"github.com/kurneo/go-template/internal/category/entities"
	"github.com/kurneo/go-template/internal/category/usecase"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/repository"
)

func NewCatRepo(l logger.Contract, db database.Contract) usecase.CategoryRepositoryContract {
	return &CatRepo{
		repository.Repository[models.Category, entities.Category, int64]{
			D: db,
			L: l,
		},
	}
}
