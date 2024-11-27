package service

import (
	"github.com/RomanshkVolkov/test-api/internal/core/domain"
)

func SchemaFieldsError[T interface{}](schema map[string][]string) domain.APIResponse[T] {
	return domain.APIResponse[T]{
		Success: false,
		Message: domain.Message{
			En: "Check the red fields",
			Es: "Verifique los campos en rojo",
		},
		SchemaError: schema,
	}
}
