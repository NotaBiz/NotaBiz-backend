// Package utils provides utility functions for general use.
package utils

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// customErrorMessages defines custom error messages for specific validation tags.
var customErrorMessages = map[string]string{
	"required": "Field %s is required",
	"gte":      "Field %s must be greater than or equal to %s",
	"lte":      "Field %s must be less than or equal to %s",
	"gt":       "Field %s must be greater than %s",
	"lt":       "Field %s must be less than %s",
	"min":      "Field %s can't contain less than %s character",
	"max":      "Field %s can't contain more than %s character",
	"datetime": "Field %s must be a valid datetime",
	"email":    "Field %s must be a valid email address",
	"e164":     "Field %s must be a valid e164 phone number",
}

// ValidateStruct validates the fields of a given struct based on the validation tags.
//
// Parameters:
//   - input: The struct to be validated.
//
// Returns:
//   - A map of validation errors, where the key is the field name and the value is the error message.
//   - Returns nil if there are no validation errors.
func ValidateStruct(input interface{}) map[string]string {
	validate := validator.New()
	val := reflect.ValueOf(input)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Handle slice of structs
	if val.Kind() == reflect.Slice {
		errors := make(map[string]string)
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()
			if err := validate.Struct(item); err != nil {
				if _, ok := err.(*validator.InvalidValidationError); ok {
					continue
				}
				for _, ve := range err.(validator.ValidationErrors) {
					errors[fmt.Sprintf("[%d].%s", i, ve.Field())] = buildErrorMessage(ve)
				}
			}
		}
		if len(errors) > 0 {
			return errors
		}
		return nil
	}

	// Handle single struct
	if err := validate.Struct(input); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return map[string]string{"error": "Invalid input for validation"}
		}
		errors := make(map[string]string)
		for _, ve := range err.(validator.ValidationErrors) {
			errors[ve.Field()] = buildErrorMessage(ve)
		}
		return errors
	}

	return nil
}

// buildErrorMessage generates a custom error message for a validation error.
//
// Parameters:
//   - err: The validation error.
//
// Returns:
//   - A string containing the custom error message.
func buildErrorMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	if msg, exists := customErrorMessages[tag]; exists {
		return fmt.Sprintf(msg, field, param)
	}
	return fmt.Sprintf("Field %s is invalid", field)
}

// ValidateRequest validates the request body and binds it to the provided struct.
//
// Parameters:
//   - c: The Gin context.
//   - input: A pointer to the struct where the request data will be bound.
//
// Returns:
//   - A boolean indicating whether the validation was successful.
//   - Sends a 400 Bad Request response with error details if validation fails.
func ValidateRequest[T any](c *gin.Context, input *T) bool {
	if err := c.ShouldBind(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format : " + err.Error()})
		return false
	}
	errors := ValidateStruct(input)
	if errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return false
	}

	return true
}
