package usecase

import (
	"context"
	"github.com/kurneo/go-template/internal/category/entities"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/support/paginate"
)

type (
	CategoryUseCaseContract interface {
		List(ctx context.Context, filter map[string]string, sort map[string]string, page, perPage int) ([]entities.Category, *paginate.Paginator, error.Contract)
		Store(ctx context.Context, dto CategoryDTO) (*entities.Category, error.Contract)
		Get(ctx context.Context, id int) (*entities.Category, error.Contract)
		Update(ctx context.Context, cat *entities.Category, dto CategoryDTO) error.Contract
		Delete(ctx context.Context, cat *entities.Category) error.Contract
	}

	CategoryRepositoryContract interface {
		List(ctx context.Context, filter map[string]string, sort map[string]string, page, perPage int) ([]entities.Category, *paginate.Paginator, error.Contract)
		Store(ctx context.Context, cat *entities.Category) error.Contract
		Get(ctx context.Context, id int) (*entities.Category, error.Contract)
		Update(ctx context.Context, cat *entities.Category) error.Contract
		UpdateDefault(ctx context.Context, except *entities.Category) error.Contract
		Delete(ctx context.Context, cat *entities.Category) error.Contract
	}

	CategoryDTO interface {
		GetName() string
		GetDescription() *string
		GetStatus() int
		GetIsDefault() bool
	}
)
