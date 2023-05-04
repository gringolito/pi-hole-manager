package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v Error) FromError(err validator.FieldError) Error {
	v.Field = err.Field()

	if err.Tag() == "required" {
		v.Message = fmt.Sprintf("Missing required field")
	} else {
		v.Message = fmt.Sprintf("Invalid value '%s' for type '%s'", err.Value(), err.Tag())
	}

	return v
}

var validate = validator.New()

func Validate(obj interface{}) []Error {
	err := validate.Struct(obj)
	if err != nil {
		errors := []Error{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, Error{}.FromError(err))
		}
		return errors
	}

	return nil
}
