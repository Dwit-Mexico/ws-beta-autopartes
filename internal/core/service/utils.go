package service

import (
	"github.com/RomanshkVolkov/test-api/internal/core/domain"
)

func FieldError[T interface{}](schema map[string][]string) domain.APIResponse[T, any] {
	return domain.APIResponse[T, any]{
		Success: false,
		Message: domain.Message{
			En: "Validation error",
			Es: "Error de validación",
		},
		SchemaError: schema,
	}
}
