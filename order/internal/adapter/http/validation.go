package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func generateValidationMessage(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is a required field", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", field, strings.ReplaceAll(e.Param(), " ", ", "))
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func writeValidationError(w http.ResponseWriter, err error) {
	var validationErrs validator.ValidationErrors

	if errors.As(err, &validationErrs) {
		var errs []FieldError

		for _, fe := range validationErrs {
			errs = append(errs, FieldError{
				Field:   strings.ToLower(fe.Field()),
				Message: generateValidationMessage(fe),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Code:             "VALIDATION_ERROR",
			Message:          "invalid request data",
			ValidationErrors: errs,
		})
		return
	}

	writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request data", nil)
}
