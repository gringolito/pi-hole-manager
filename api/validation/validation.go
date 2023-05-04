package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func newError(err validator.FieldError) *error {
	var message string
	if err.Tag() == "required" {
		message = fmt.Sprintf("Missing required field")
	} else {
		message = fmt.Sprintf("Invalid value '%s' for type '%s'", err.Value(), err.Tag())
	}

	return &error{
		Field:   err.Field(),
		Message: message,
	}
}

var validate = validator.New()

func Validate(obj interface{}) []error {
	err := validate.Struct(obj)
	if err != nil {
		errors := []error{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, *newError(err))
		}
		return errors
	}

	return nil
}
