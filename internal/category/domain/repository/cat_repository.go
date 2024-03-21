package repository

import (
	"context"
	"github.com/kurneo/go-template/internal/category/domain/entity"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/support/page_list"
)

type CategoryRepositoryContract interface {
	List(ctx context.Context, filter map[string]string, sort map[string]string, page, perPage int) (*page_list.PageList[entity.Category], error.Contract)
	Store(ctx context.Context, cat *entity.Category) error.Contract
	Get(ctx context.Context, id int64) (*entity.Category, error.Contract)
	Update(ctx context.Context, cat *entity.Category) error.Contract
	UpdateDefault(ctx context.Context, except *entity.Category) error.Contract
	Delete(ctx context.Context, cat *entity.Category) error.Contract
}
