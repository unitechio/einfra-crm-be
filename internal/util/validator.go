
package util

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// a single instance of the validator
var validate = validator.New()

// ValidateStruct validates a struct and returns a formatted error message if the validation fails.
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {

		// this check is only needed when your code could produce an invalid value for validation itself.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("invalid validation error: %w", err)
		}

		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			// You can customize the error messages here
			msg := fmt.Sprintf("Validation failed on field '%s' with tag '%s'", err.Field(), err.Tag())
			errorMessages = append(errorMessages, msg)
		}

		return fmt.Errorf(strings.Join(errorMessages, "; "))
	}

	return nil
}
