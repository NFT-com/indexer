package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	v := Validator{
		validate: validator.New(),
	}

	return &v
}

func (v *Validator) Request(request interface{}) error {
	return v.validate.Struct(request)
}
