package requests

import (
	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *RegisterRequest) Validate(v *validator.Validate) error {
	if err := v.Struct(r); err != nil {
		return err
	}

	return nil
}
