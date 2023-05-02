package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *ValidationError) FromError(err validator.FieldError) {
	v.Field = err.Field()

	if err.Tag() == "required" {
		v.Message = fmt.Sprintf("Missing required field")
	} else {
		v.Message = fmt.Sprintf("Invalid value '%s' for type '%s'", err.Value(), err.Tag())
	}
}

var validate = validator.New()

func Validate(obj interface{}) []ValidationError {
	err := validate.Struct(obj)
	if err != nil {
		errors := []ValidationError{}
		for _, err := range err.(validator.ValidationErrors) {
			validationError := ValidationError{}
			validationError.FromError(err)
			errors = append(errors, validationError)
		}
		return errors
	}

	return nil
}
