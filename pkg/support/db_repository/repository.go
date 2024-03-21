package db_repository

import (
	"context"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/support/page_list"
	"github.com/kurneo/go-template/pkg/support/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	PrimaryKey = interface {
		int64 | string
	}

	Entity[P PrimaryKey] interface {
		ToMap() map[string]interface{}
	}

	Model[P PrimaryKey, E Entity[P]] interface {
		ToEntity() *E
		FromEntity(e E) interface{}
		TableName() string
	}

	Condition interface {
		GetQuery() string
		GetValues() []any
	}

	Preload interface {
		GetRelation() string
		GetCondition() *Condition
		GetSelectColumns() []string
	}

	Scope func(*gorm.DB) *gorm.DB

	Param struct {
		Condition Condition
		Preloads  []Preload
		Selects   []string
		Scopes    []Scope
		Orders    map[string]string
		Page      int
		Limit     int
	}

	Repository[M Model[P, E], E Entity[P], P PrimaryKey] struct {
		D database.Contract
	}
)

func (p Param) GetPreload() []Preload {
	if p.Preloads != nil {
		return p.Preloads
	}
	return []Preload{}
}

func (p Param) GetSelectColumns() []string {
	if p.Selects != nil {
		return p.Selects
	}
	return []string{"*"}
}

func (p Param) GetScopes() []Scope {
	if p.Scopes != nil {
		return p.Scopes
	}
	return []Scope{}
}

func (p Param) GetOrders() map[string]string {
	return p.Orders
}

func (p Param) GetCondition() Condition {
	return p.Condition
}

func (p Param) GetPage() int {
	if p.Page == 0 {
		return 1
	}
	return p.Page
}

func (p Param) GetLimit() int {
	if p.Limit == 0 {
		return 10
	}
	return p.Limit
}

func (r Repository[M, E, P]) All(ctx context.Context, p Param) ([]E, error) {
	var m M
	var list []M

	q := r.D.GetDB(ctx).Table(m.TableName())
	ApplyEagerLoad(q, p.GetPreload())
	ApplySelectColumns(q, p.GetSelectColumns())
	ApplyScopes(q, p.GetScopes())
	ApplyOrder(q, p.GetOrders())

	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}

	listE := slices.Map[M, E](list, func(model M) E {
		return *model.ToEntity()
	})

	return listE, nil
}

func (r Repository[M, E, P]) AllBy(ctx context.Context, p Param) ([]E, error) {
	var m M
	var list []M

	q := r.D.GetDB(ctx).Table(m.TableName())
	ApplyCondition(q, p.GetCondition())
	ApplySelectColumns(q, p.GetSelectColumns())
	ApplyEagerLoad(q, p.GetPreload())
	ApplyScopes(q, p.GetScopes())
	ApplyOrder(q, p.GetOrders())

	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}

	listE := slices.Map[M, E](list, func(model M) E {
		return *model.ToEntity()
	})

	return listE, nil
}

func (r Repository[M, E, P]) AllByWithPaginate(ctx context.Context, p Param) (*page_list.PageList[E], error) {
	var m M
	var list []M
	var count int64

	q := r.D.GetDB(ctx).Table(m.TableName())
	ApplyCondition(q, p.GetCondition())
	ApplySelectColumns(q, p.GetSelectColumns())
	ApplyScopes(q, p.GetScopes())

	if err := q.Count(&count).Error; err != nil {
		return nil, err
	}

	ApplyEagerLoad(q, p.GetPreload())
	ApplyOrder(q, p.GetOrders())

	ApplyPaginate(q, p.GetPage(), p.GetLimit())

	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}

	listE := slices.Map[M, E](list, func(model M) E {
		return *model.ToEntity()
	})

	return page_list.NewPageList[E](listE, p.GetPage(), p.GetLimit(), count), nil
}

func (r Repository[M, E, P]) FirstBy(ctx context.Context, p Param) (*E, error) {
	var m M
	q := r.D.GetDB(ctx).Table(m.TableName())
	ApplyCondition(q, p.GetCondition())
	ApplyEagerLoad(q, p.GetPreload())
	ApplySelectColumns(q, p.GetSelectColumns())
	ApplyScopes(q, p.GetScopes())

	if err := q.First(&m).Error; err == nil {
		entity := m.ToEntity()
		return entity, nil
	} else {
		if r.D.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
}

func (r Repository[M, E, P]) FindByID(ctx context.Context, id P, p Param) (*E, error) {
	var m M
	q := r.D.GetDB(ctx).Table(m.TableName())
	ApplyCondition(q, Equal[P]("id", id))
	ApplyEagerLoad(q, p.GetPreload())
	ApplySelectColumns(q, p.GetSelectColumns())
	ApplyScopes(q, p.GetScopes())

	if err := q.First(&m).Error; err == nil {
		entity := m.ToEntity()
		return entity, nil
	} else {
		if r.D.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
}

func (r Repository[M, E, P]) Insert(ctx context.Context, e *E) error {
	var m M
	model := m.FromEntity(*e).(*M)
	if err := r.D.GetDB(ctx).Omit(clause.Associations).Create(&model).Error; err != nil {
		return err
	}
	*e = *(*model).ToEntity()
	return nil
}

func (r Repository[M, E, P]) InsertMany(ctx context.Context, es *[]E) error {
	var m M
	models := slices.Map[E, *M](*es, func(v E) *M {
		return m.FromEntity(v).(*M)
	})
	if err := r.D.GetDB(ctx).Omit(clause.Associations).Create(&models).Error; err != nil {
		return err
	}

	c := slices.Map[*M, E](models, func(v *M) E {
		return *(*v).ToEntity()
	})
	*es = c
	return nil
}

func (r Repository[M, E, P]) Update(ctx context.Context, e *E) error {
	var m M
	model := m.FromEntity(*e).(*M)
	if err := r.D.GetDB(ctx).Omit(clause.Associations).Updates(model).Error; err != nil {
		return err
	}
	*e = *(*model).ToEntity()
	return nil
}

func (r Repository[M, E, P]) Delete(ctx context.Context, e *E) error {
	var m M
	model := m.FromEntity(*e).(*M)
	if err := r.D.GetDB(ctx).Omit(clause.Associations).Delete(model).Error; err != nil {
		return err
	}
	return nil
}

func (r Repository[M, E, P]) Exists(ctx context.Context, id int) (bool, error) {
	var m M
	var exists bool
	err := r.D.GetDB(ctx).Table(m.TableName()).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r Repository[M, E, P]) ExistsBy(ctx context.Context, c Condition) (bool, error) {
	var m M
	var exists bool
	q := r.D.GetDB(ctx).Table(m.TableName()).Select("count(*) > 0")
	ApplyCondition(q, c)
	err := q.Find(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}
