package db_repository

import (
	"github.com/kurneo/go-template/pkg/support/helper"
	"github.com/kurneo/go-template/pkg/support/slices"
	"gorm.io/gorm"
)

func ApplyCondition(query *gorm.DB, c Condition) {
	if c != nil {
		query.Where(c.GetQuery(), c.GetValues()...)
	}
}

func ApplyOrder(query *gorm.DB, orders map[string]string) {
	for f, d := range orders {
		query.Order(f + " " + d)
	}
}

func ApplyEagerLoad(query *gorm.DB, preloads []Preload) {
	if preloads == nil || len(preloads) == 0 {
		return
	}
	for _, p := range preloads {
		cd := p.GetCondition()
		if *cd != nil {
			query.Preload(p.GetRelation(), func(tx *gorm.DB) *gorm.DB {
				return tx.Select(p.GetSelectColumns()).Where((*cd).GetQuery(), (*cd).GetValues()...)
			})
		} else {
			query.Preload(p.GetRelation(), func(tx *gorm.DB) *gorm.DB {
				return tx.Select(p.GetSelectColumns())
			})
		}
	}
}

func ApplySelectColumns(query *gorm.DB, columns []string) {
	if columns != nil && len(columns) > 0 {
		query.Select(columns)
	}
}

func ApplyScopes(query *gorm.DB, s []Scope) {
	query.Scopes(slices.Map[Scope, func(db *gorm.DB) *gorm.DB](s, func(v Scope) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return v(db)
		}
	})...)
}

func ApplyPaginate(query *gorm.DB, page, limit int) {
	query.Offset(helper.ResolveOffset(page, limit)).Limit(limit)
}
