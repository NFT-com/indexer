package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator represents a validator.
type Validator struct {
	validate *validator.Validate
}

// New returns a new validator.
func New() *Validator {
	v := Validator{
		validate: validator.New(),
	}

	return &v
}

// Request validates the request. Returns error if request is invalid.
func (v *Validator) Request(request interface{}) error {
	return v.validate.Struct(request)
}
