package datasource

import (
	"context"
	"github.com/kurneo/go-template/internal/category/data/model"
	"github.com/kurneo/go-template/internal/category/domain/entity"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/support/paginate"
	"github.com/kurneo/go-template/pkg/support/repository"
	"strings"
)

type CatDatasource struct {
	repository.Repository[model.Category, entity.Category, int64]
}

func (r CatDatasource) List(
	ctx context.Context,
	filters map[string]string,
	sort map[string]string,
	page,
	perPage int,
) ([]entity.Category, *paginate.Paginator, error) {
	var c []repository.Condition
	c = append(c, repository.Contains("name", filters["name"]))

	if filters["status"] != "" {
		c = append(c, repository.Equal[string]("status", filters["status"]))
	}

	if filters["created_at"] != "" {
		createdAt := strings.Split(filters["created_at"], ",")
		if len(createdAt) == 2 {
			c = append(c, repository.Between[string]("created_at", createdAt[0], createdAt[1]))
		}
	}

	return r.AllBy(
		ctx,
		repository.And(c...),
		nil,
		nil,
		[]string{
			"*",
		},
		sort,
		page,
		perPage,
	)
}

func (r CatDatasource) Store(ctx context.Context, cat *entity.Category) error {
	return r.Insert(ctx, cat)
}

func (r CatDatasource) UpdateDefault(ctx context.Context, except *entity.Category) error {
	var m model.Category
	err := r.D.GetDB(ctx).Table(m.TableName()).Where("id != ?", except.ID).Update("is_default", false).Error

	if err != nil {
		return err
	}
	return nil
}

func (r CatDatasource) Get(ctx context.Context, id int64) (*entity.Category, error) {
	return r.FindByID(
		ctx,
		id,
		nil,
		nil,
		[]string{
			"*",
		},
	)
}

func (r CatDatasource) Update(ctx context.Context, cat *entity.Category) error {
	return r.Repository.Update(ctx, cat)
}

func (r CatDatasource) Delete(ctx context.Context, cat *entity.Category) error {
	return r.Repository.Delete(ctx, cat)
}

func NewCatDatasource(db database.Contract) *CatDatasource {
	return &CatDatasource{
		repository.Repository[model.Category, entity.Category, int64]{
			D: db,
		},
	}
}
