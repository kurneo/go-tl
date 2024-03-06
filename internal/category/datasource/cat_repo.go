package datasource

import (
	"context"
	"github.com/kurneo/go-template/internal/category/datasource/models"
	"github.com/kurneo/go-template/internal/category/entities"
	pkgErr "github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/support/paginate"
	"github.com/kurneo/go-template/pkg/support/repository"
	"strings"
)

type CatRepo struct {
	repository.Repository[models.Category, entities.Category]
}

func (r CatRepo) List(
	ctx context.Context,
	filters map[string]string,
	sort map[string]string,
	page,
	perPage int,
) ([]entities.Category, *paginate.Paginator, pkgErr.Contract) {
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

func (r CatRepo) Store(ctx context.Context, cat *entities.Category) pkgErr.Contract {
	return r.Insert(ctx, cat)
}

func (r CatRepo) UpdateDefault(ctx context.Context, except *entities.Category) pkgErr.Contract {
	var m models.Category
	err := r.D.GetDB(ctx).Table(m.TableName()).Where("id != ?", except.ID).Update("is_default", false).Error

	if err != nil {
		r.L.Error(err)
		return pkgErr.NewDatasource(err)
	}
	return nil
}

func (r CatRepo) Get(ctx context.Context, id int) (*entities.Category, pkgErr.Contract) {
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

func (r CatRepo) Update(ctx context.Context, cat *entities.Category) pkgErr.Contract {
	return r.Repository.Update(ctx, cat)
}

func (r CatRepo) Delete(ctx context.Context, cat *entities.Category) pkgErr.Contract {
	return r.Repository.Delete(ctx, cat)
}
