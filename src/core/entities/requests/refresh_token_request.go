package requests

import (
	"github.com/go-playground/validator/v10"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func (r *RefreshTokenRequest) Validate(v *validator.Validate) error {
	if err := v.Struct(r); err != nil {
		return err
	}

	return nil
}
