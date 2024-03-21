package page_list

import "github.com/kurneo/go-template/pkg/support/helper"

type Paginate struct {
	Page       int   `json:"page"`
	Limit      int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PageList[T any] struct {
	Paginate *Paginate `json:"paginate"`
	List     []T       `json:"list"`
}

func populate(page, limit int, total int64) *Paginate {
	return &Paginate{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: helper.ResolveTotalPages(total, limit),
	}
}

func NewPageList[T any](list []T, page, limit int, total int64) *PageList[T] {
	return &PageList[T]{
		List:     list,
		Paginate: populate(page, limit, total),
	}
}
