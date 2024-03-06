package entities

import "time"

const (
	StatusPublish = 10
	StatusDraft   = 11
)

type Category struct {
	ID          int
	Name        string
	Description *string
	Status      int
	IsDefault   bool
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func (c Category) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":          c.ID,
		"name":        c.Name,
		"description": c.Description,
		"status":      c.Status,
		"is_default":  c.IsDefault,
		"created_at":  c.CreatedAt,
		"updated_at":  c.UpdatedAt,
	}
}

func (c Category) IsPublic() bool {
	return c.Status == StatusPublish
}
