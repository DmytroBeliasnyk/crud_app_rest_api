package dto

import "github.com/go-playground/validator/v10"

type SignUpDTO struct {
	Name     string `json:"name" validate:"required,gte=2"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,gte=3"`
	Password string `json:"password" validate:"required,gte=8"`
}

type SignInDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (su *SignUpDTO) Validate() error {
	return validate.Struct(su)
}
