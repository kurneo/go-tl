package repository

import (
	"context"
	"github.com/kurneo/go-template/pkg/database"
	errPkg "github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/paginate"
	"github.com/kurneo/go-template/pkg/support/slices"
	"gorm.io/gorm/clause"
)

type Repository[M Model[E], E Entity] struct {
	L logger.Contract
	D database.Contract
}

func (r Repository[M, E]) All(ctx context.Context, vars ...interface{}) ([]E, *paginate.Paginator, errPkg.Contract) {
	var m M
	var list []M
	var count int64

	q := r.D.GetDB(ctx).Table(m.TableName())
	ApplyEagerLoad(q, GetPreload(vars...))
	ApplySelectColumns(q, GetSelectColumns(vars...))
	ApplyScopes(q, GetScopes(vars...))
	ApplyOrder(q, GetOrders(vars...))

	if err := q.Count(&count).Error; err != nil {
		r.L.Error(err)
		return nil, nil, errPkg.NewDatasource(err)
	}

	page := GetPage(vars...)
	limit := GetPerPage(vars...)
	pg := ApplyPaginate(q, page, limit)

	if err := q.Find(&list).Error; err != nil {
		r.L.Error(err)
		return nil, nil, errPkg.NewDatasource(err)
	}

	listE := slices.Map[M, E](list, func(model M) E {
		return *model.ToEntity()
	})

	if !pg {
		return listE, nil, nil
	}

	return listE, paginate.Populate(*page, *limit, count), nil
}

func (r Repository[M, E]) AllBy(ctx context.Context, c Condition, vars ...interface{}) ([]E, *paginate.Paginator, errPkg.Contract) {
	var m M
	var list []M
	var count int64

	q := r.D.GetDB(ctx).Table(m.TableName())
	ApplyCondition(q, c)
	ApplySelectColumns(q, GetSelectColumns(vars...))
	ApplyScopes(q, GetScopes(vars...))

	if err := q.Count(&count).Error; err != nil {
		r.L.Error(err)
		return nil, nil, errPkg.NewDatasource(err)
	}

	ApplyEagerLoad(q, GetPreload(vars...))
	ApplyOrder(q, GetOrders(vars...))

	page := GetPage(vars...)
	limit := GetPerPage(vars...)
	pg := ApplyPaginate(q, page, limit)

	if err := q.Find(&list).Error; err != nil {
		r.L.Error(err)
		return nil, nil, errPkg.NewDatasource(err)
	}

	listE := slices.Map[M, E](list, func(model M) E {
		return *model.ToEntity()
	})

	if !pg {
		return listE, nil, nil
	}

	return listE, paginate.Populate(*page, *limit, count), nil
}

func (r Repository[M, E]) FirstBy(ctx context.Context, c Condition, vars ...interface{}) (*E, errPkg.Contract) {
	var m M
	q := r.D.GetDB(ctx).Table(m.TableName())
	ApplyCondition(q, c)
	ApplyEagerLoad(q, GetPreload(vars...))
	ApplySelectColumns(q, GetSelectColumns(vars...))
	ApplyScopes(q, GetScopes(vars...))

	if err := q.First(&m).Error; err == nil {
		entity := m.ToEntity()
		return entity, nil
	} else {
		if r.D.IsNotFound(err) {
			return nil, nil
		}
		r.L.Error(err)
		return nil, errPkg.NewDatasource(err)
	}
}

func (r Repository[M, E]) FindByID(ctx context.Context, id int, vars ...interface{}) (*E, errPkg.Contract) {
	var m M
	q := r.D.GetDB(ctx).Table(m.TableName())
	ApplyCondition(q, Equal[int]("id", id))
	ApplyEagerLoad(q, GetPreload(vars...))
	ApplySelectColumns(q, GetSelectColumns(vars...))
	ApplyScopes(q, GetScopes(vars...))

	if err := q.First(&m).Error; err == nil {
		entity := m.ToEntity()
		return entity, nil
	} else {
		if r.D.IsNotFound(err) {
			return nil, nil
		}
		r.L.Error(err)
		return nil, errPkg.NewDatasource(err)
	}
}

func (r Repository[M, E]) Insert(ctx context.Context, e *E) errPkg.Contract {
	var m M
	model := m.FromEntity(*e).(*M)
	if err := r.D.GetDB(ctx).Omit(clause.Associations).Create(&model).Error; err != nil {
		r.L.Error(err)
		return errPkg.NewDatasource(err)
	}
	*e = *(*model).ToEntity()
	return nil
}

func (r Repository[M, E]) InsertMany(ctx context.Context, es *[]E) errPkg.Contract {
	var m M
	models := slices.Map[E, *M](*es, func(v E) *M {
		return m.FromEntity(v).(*M)
	})
	if err := r.D.GetDB(ctx).Omit(clause.Associations).Create(&models).Error; err != nil {
		r.L.Error(err)
		return errPkg.NewDatasource(err)
	}

	c := slices.Map[*M, E](models, func(v *M) E {
		return *(*v).ToEntity()
	})
	*es = c
	return nil
}

func (r Repository[M, E]) Update(ctx context.Context, e *E) errPkg.Contract {
	var m M
	model := m.FromEntity(*e).(*M)
	if err := r.D.GetDB(ctx).Omit(clause.Associations).Updates(model).Error; err != nil {
		r.L.Error(err)
		return errPkg.NewDatasource(err)
	}
	*e = *(*model).ToEntity()
	return nil
}

func (r Repository[M, E]) Delete(ctx context.Context, e *E) errPkg.Contract {
	var m M
	model := m.FromEntity(*e).(*M)
	if err := r.D.GetDB(ctx).Omit(clause.Associations).Delete(model).Error; err != nil {
		r.L.Error(err)
		return errPkg.NewDatasource(err)
	}
	return nil
}

func (r Repository[M, E]) Exists(ctx context.Context, id int) (bool, errPkg.Contract) {
	var m M
	var exists bool
	err := r.D.GetDB(ctx).Table(m.TableName()).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).Error
	if err != nil {
		r.L.Error(err)
		return false, errPkg.NewDatasource(err)
	}
	return exists, nil
}

func (r Repository[M, E]) ExistsBy(ctx context.Context, c Condition) (bool, errPkg.Contract) {
	var m M
	var exists bool
	q := r.D.GetDB(ctx).Table(m.TableName()).Select("count(*) > 0")
	ApplyCondition(q, c)
	err := q.Find(&exists).Error
	if err != nil {
		r.L.Error(err)
		return false, errPkg.NewDatasource(err)
	}
	return exists, nil
}