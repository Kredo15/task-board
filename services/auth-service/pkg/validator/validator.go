package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var ErrInvalidInputData = errors.New("invalid input data")

type Validator struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	return Validator{
		validate: &validator.Validate{},
	}
}

func (v *Validator) Struct(data any) error {
	return v.validate.Struct(data)
}
