package v1

type LoginFormData struct {
	Email    string `validate:"required,email" form:"email" json:"email"`
	Password string `validate:"required" form:"password" json:"password"`
}
