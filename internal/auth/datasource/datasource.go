package datasource

import (
	"github.com/kurneo/go-template/internal/auth/datasource/models"
	"github.com/kurneo/go-template/internal/auth/entities"
	"github.com/kurneo/go-template/internal/auth/usecase"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/repository"
)

func NewUserRepo[T repository.PrimaryKey](lg logger.Contract, db database.Contract) usecase.UserRepositoryContract[T] {
	return &UserRepo[T]{
		Repository: repository.Repository[models.User[T], entities.User[T], T]{
			D: db,
			L: lg,
		},
	}
}
