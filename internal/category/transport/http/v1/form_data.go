package v1

type CategoryFormData struct {
	Name        string  `validate:"required" form:"name" json:"name"`
	Description *string `validate:"omitempty" form:"description" json:"description"`
	Status      int     `validate:"required,oneof=1 2" form:"status" json:"status"`
	IsDefault   *bool   `validate:"required" form:"is_default" json:"is_default"`
}

func (c CategoryFormData) GetName() string {
	return c.Name
}

func (c CategoryFormData) GetDescription() *string {
	return c.Description
}

func (c CategoryFormData) GetStatus() int {
	return c.Status
}

func (c CategoryFormData) GetIsDefault() bool {
	isDefault := true

	if c.IsDefault == nil || *(c.IsDefault) == false {
		isDefault = false
	}
	return isDefault
}
