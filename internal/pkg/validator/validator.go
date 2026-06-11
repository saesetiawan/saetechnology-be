package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(i any) string
}

type validatorImpl struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	return &validatorImpl{
		validate: validator.New(),
	}
}

func (v *validatorImpl) Validate(i any) string {
	err := v.validate.Struct(i)
	message := ""
	if err == nil {
		return message
	}

	for _, e := range err.(validator.ValidationErrors) {
		field := e.Field()
		switch e.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s", field, e.Param())
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", field)
		default:
			message = fmt.Sprintf("%s is not valid", field)
		}
	}

	return message
}
