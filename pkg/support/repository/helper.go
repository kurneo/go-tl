package repository

import (
	"github.com/kurneo/go-template/pkg/support/paginate"
	"gorm.io/gorm"
)

func ApplyCondition(query *gorm.DB, cds Condition) {
	query.Where(cds.GetQuery(), cds.GetValues()...)
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

func ApplyScopes(query *gorm.DB, s []func(*gorm.DB) *gorm.DB) {
	query.Scopes(s...)
}

func ApplyPaginate(query *gorm.DB, page, perPage *int) bool {
	if page != nil && perPage != nil {
		query.Offset(paginate.ResolveOffset(*page, *perPage)).Limit(*perPage)
		return true
	}
	return false
}

func GetPreload(vars ...interface{}) []Preload {
	var p []Preload
	if len(vars) > 0 && vars[0] != nil {
		temp := vars[0]
		switch temp.(type) {
		case []Preload:
			p = temp.([]Preload)
			break
		default:
			p = []Preload{temp.(Preload)}
		}
	}
	return p
}

func GetScopes(vars ...interface{}) []func(*gorm.DB) *gorm.DB {
	s := make([]func(*gorm.DB) *gorm.DB, 0)
	if len(vars) > 1 && vars[1] != nil {
		s = vars[1].([]func(*gorm.DB) *gorm.DB)
	}
	return s
}

func GetSelectColumns(vars ...interface{}) []string {
	s := []string{"*"}
	if len(vars) > 2 && vars[2] != nil {
		s = vars[2].([]string)
	}
	return s
}

func GetOrders(vars ...interface{}) map[string]string {
	o := make(map[string]string)
	if len(vars) > 4 && vars[3] != nil {
		o = vars[3].(map[string]string)
	}
	return o
}

func GetPage(vars ...interface{}) *int {
	var page int
	if len(vars) >= 5 && vars[4] != nil {
		page = vars[4].(int)
	}

	if page == 0 {
		page = 1
	}

	return &page
}

func GetPerPage(vars ...interface{}) *int {
	var val int
	if len(vars) >= 6 && vars[5] != nil {
		val = vars[5].(int)
	}

	if val == 0 {
		val = 10
	}

	return &val
}
