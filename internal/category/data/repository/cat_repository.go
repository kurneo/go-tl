package repository

import (
	"context"
	"github.com/kurneo/go-template/internal/category/data/datasource"
	"github.com/kurneo/go-template/internal/category/domain/entity"
	"github.com/kurneo/go-template/internal/category/domain/repository"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/support/paginate"
)

type CatRepository struct {
	d *datasource.CatDatasource
}

func (c CatRepository) List(ctx context.Context, filter map[string]string, sort map[string]string, page, perPage int) ([]entity.Category, *paginate.Paginator, error.Contract) {
	l, p, err := c.d.List(ctx, filter, sort, page, perPage)
	if err != nil {
		return nil, nil, error.NewDatasource(err)
	}
	return l, p, nil
}
func (c CatRepository) Store(ctx context.Context, cat *entity.Category) error.Contract {
	err := c.d.Store(ctx, cat)
	if err != nil {
		return error.NewDatasource(err)
	}
	return nil
}
func (c CatRepository) Get(ctx context.Context, id int64) (*entity.Category, error.Contract) {
	cat, err := c.d.Get(ctx, id)
	if err != nil {
		return nil, error.NewDatasource(err)
	}
	return cat, nil
}
func (c CatRepository) Update(ctx context.Context, cat *entity.Category) error.Contract {
	err := c.d.Update(ctx, cat)
	if err != nil {
		return error.NewDatasource(err)
	}
	return nil
}
func (c CatRepository) UpdateDefault(ctx context.Context, except *entity.Category) error.Contract {
	err := c.d.UpdateDefault(ctx, except)
	if err != nil {
		return error.NewDatasource(err)
	}
	return nil
}
func (c CatRepository) Delete(ctx context.Context, cat *entity.Category) error.Contract {
	err := c.d.Delete(ctx, cat)
	if err != nil {
		return error.NewDatasource(err)
	}
	return nil
}

func NewCatRepo(d *datasource.CatDatasource) repository.CategoryRepositoryContract {
	return &CatRepository{
		d: d,
	}
}
