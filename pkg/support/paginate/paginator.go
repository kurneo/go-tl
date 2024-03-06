package paginate

type Paginator struct {
	Page       int   `json:"page"`
	Limit      int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Populate(page, limit int, total int64) *Paginator {
	return &Paginator{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: ResolveTotalPages(total, limit),
	}
}

func ResolveOffset(page, limit int) int {
	if page < 0 {
		page = 1
	}
	return (page - 1) * limit
}

func ResolveTotalPages(total int64, limit int) int {

	if total <= 0 {
		return 0
	}

	if limit <= 0 {
		return 0
	}

	totalPages := total / int64(limit)

	if total%int64(limit) > 0 {
		totalPages++
	}

	return int(totalPages)
}
