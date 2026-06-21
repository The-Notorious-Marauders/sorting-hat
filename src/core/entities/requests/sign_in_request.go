package requests

import (
	"github.com/go-playground/validator/v10"
)

type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r *SignInRequest) Validate(v *validator.Validate) error {
	if err := v.Struct(r); err != nil {
		return err
	}

	return nil
}
