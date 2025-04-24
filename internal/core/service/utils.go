package service

import (
	"github.com/RomanshkVolkov/ws-beta-autopartes/internal/core/domain"
)

func SchemaFieldsError[T any](schema map[string][]string) domain.APIResponse[T] {
	return domain.APIResponse[T]{
		Success: false,
		Message: domain.Message{
			En: "Check the red fields",
			Es: "Verifique los campos en rojo",
		},
		SchemaError: schema,
	}
}
