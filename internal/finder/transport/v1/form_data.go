package v1

type CreateDirFormData struct {
	Name string `json:"name" validate:"required" form:"name"`
	Path string `json:"path" form:"path"`
}

func (c CreateDirFormData) GetName() string {
	return c.Name
}

func (c CreateDirFormData) GetPath() string {
	return c.Path
}

type RenameFormData struct {
	Items map[string]string `form:"items" validate:"required"`
}

func (r RenameFormData) GetItems() map[string]string {
	return r.Items
}

type DeleteFormData struct {
	Items []string `form:"items" json:"items" validate:"required"`
}

func (d DeleteFormData) GetItems() []string {
	return d.Items
}

type CopyFormData struct {
	Path  string   `form:"path" json:"path" validate:"required"`
	Items []string `form:"items" json:"items" validate:"required"`
}

func (c CopyFormData) GetPath() string {
	return c.Path
}

func (c CopyFormData) GetItems() []string {
	return c.Items
}

type CutFormData struct {
	CopyFormData
}
