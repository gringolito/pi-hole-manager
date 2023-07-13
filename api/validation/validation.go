package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type error struct {
	Field  string      `json:"field"`
	Reason string      `json:"reason"`
	Value  interface{} `json:"value"`
}

func newError(err validator.FieldError) error {
	var reason string
	if err.Tag() == "required" {
		reason = fmt.Sprintf("The %s field is required.", err.Field())
	} else {
		reason = fmt.Sprintf("The %s field must be of type %s.", err.Field(), err.Tag())
	}

	return error{
		Field:  err.Field(),
		Reason: reason,
		Value:  err.Value(),
	}
}

var validate = validator.New()

func Validate(obj interface{}) []error {
	err := validate.Struct(obj)
	if err != nil {
		numErrors := len(err.(validator.ValidationErrors))
		errors := make([]error, 0, numErrors)
		for _, fieldError := range err.(validator.ValidationErrors) {
			errors = append(errors, newError(fieldError))
		}
		return errors
	}

	return nil
}
