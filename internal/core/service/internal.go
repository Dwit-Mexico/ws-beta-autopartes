package service

import (
	"github.com/RomanshkVolkov/test-api/internal/adapters/repository"
	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	schema "github.com/RomanshkVolkov/test-api/internal/core/domain/schemas"
)

func (server Server) PermissionsSynchronization(request *domain.WebPages) (domain.APIResponse[any], error) {
	fields := schema.GenericForm[domain.WebPages]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)
	if len(failValidatedFields) > 0 {
		return SchemaFieldsError[any](failValidatedFields), nil
	}

	repo := repository.GetDBConnection(server.Host)
	err := repo.PermissionsSynchronization(fields.Data.Routes)
	if err != nil {
		return domain.APIResponse[any]{
			Success: false,
			Message: domain.Message{
				En: "Error on permissions synchronization",
				Es: "Error en la sincronización de permisos",
			},
			Error: err,
		}, nil
	}

	return domain.APIResponse[any]{
		Success: true,
		Message: domain.Message{
			En: "Permissions synchronized",
			Es: "Permisos sincronizados",
		},
	}, nil

}
