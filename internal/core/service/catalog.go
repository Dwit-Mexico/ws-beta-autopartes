package service

import (
	"github.com/RomanshkVolkov/test-api/internal/adapters/repository"
	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	schema "github.com/RomanshkVolkov/test-api/internal/core/domain/schemas"
)

func (server Server) GetKitchenByID(id uint) domain.APIResponse[domain.Kitchen, any] {
	repo := repository.GetDBConnection(server.Host)
	kitchen, err := repo.GetKitchenByID(id)
	if err != nil {
		return domain.APIResponse[domain.Kitchen, any]{
			Success: false,
			Message: domain.Message{
				En: "Error on get kitchen",
				Es: "Error al obtener cocina",
			},
			Error: err,
		}
	}

	if kitchen.ID == 0 {
		return repository.RecordNotFound[domain.Kitchen]()
	}

	return domain.APIResponse[domain.Kitchen, any]{
		Success: true,
		Message: domain.Message{
			En: "Kitchen data",
			Es: "Datos de cocina",
		},
		Data: kitchen,
	}
}
func (server Server) GetShiftByID(id uint) domain.APIResponse[domain.Shift, any] {
	repo := repository.GetDBConnection(server.Host)
	shift, err := repo.GetShiftByID(id)
	if err != nil {
		return domain.APIResponse[domain.Shift, any]{
			Success: false,
			Message: domain.Message{
				En: "Error on get kitchen",
				Es: "Error al obtener cocina",
			},
			Error: err,
		}
	}

	if shift.ID == 0 {
		return repository.RecordNotFound[domain.Shift]()
	}

	return domain.APIResponse[domain.Shift, any]{
		Success: true,
		Message: domain.Message{
			En: "Kitchen data",
			Es: "Datos de cocina",
		},
		Data: shift,
	}
}

func (server Server) UpdateGenericCatalog(request *domain.GenericCatalog, table interface{}) domain.APIResponse[domain.GenericCatalog, any] {
	fields := schema.GenericForm[domain.GenericCatalog]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)

	if len(failValidatedFields) > 0 {
		return domain.APIResponse[domain.GenericCatalog, any]{
			Success: false,
			Message: domain.Message{
				En: "Invalid fields",
				Es: "Campos inválidos",
			},
			SchemaError: failValidatedFields,
		}
	}

	repo := repository.GetDBConnection(server.Host)
	err := repo.UpdateGenericCatalog(fields.Data.ID, table, fields.Data.Name)
	if err != nil {
		return domain.APIResponse[domain.GenericCatalog, any]{
			Success: false,
			Message: domain.Message{
				En: "Error on update catalog",
				Es: "Error al actualizar catálogo",
			},
			Error: err,
		}
	}

	return domain.APIResponse[domain.GenericCatalog, any]{
		Success: true,
		Message: domain.Message{
			En: "Catalog updated",
			Es: "Catálogo actualizado",
		},
		Data: *request,
	}
}

func (server Server) DeleteGenericCatalog(id uint, table interface{}) domain.APIResponse[domain.GenericCatalog, any] {
	repo := repository.GetDBConnection(server.Host)
	err := repo.DeleteRecord(id, table)
	if err != nil {
		return domain.APIResponse[domain.GenericCatalog, any]{
			Success: false,
			Message: domain.Message{
				En: "Error on delete catalog",
				Es: "Error al eliminar catálogo",
			},
			Error: err,
		}
	}

	return domain.APIResponse[domain.GenericCatalog, any]{
		Success: true,
		Message: domain.Message{
			En: "Catalog deleted",
			Es: "Catálogo eliminado",
		},
		Data: domain.GenericCatalog{
			ID: id,
		},
	}
}
