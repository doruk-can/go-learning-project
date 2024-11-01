package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using the go-playground/validator library
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		var validationErrors []string
		errorMessages := map[string]string{
			"required": "Field '%s' is required",
			"min":      "Field '%s' must be at least %s",
			"max":      "Field '%s' must be at most %s",
			"uuid4":    "Field '%s' must be a valid UUID v4",
		}

		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			if msg, exists := errorMessages[err.Tag()]; exists {
				validationErrors = append(validationErrors, fmt.Sprintf(msg, fieldName, err.Param()))
			} else {
				validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' is invalid", fieldName))
			}
		}
		return fmt.Errorf(strings.Join(validationErrors, "; "))
	}
	return nil
}
