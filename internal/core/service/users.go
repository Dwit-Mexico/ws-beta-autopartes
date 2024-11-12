package service

import (
	"github.com/RomanshkVolkov/test-api/internal/adapters/repository"
	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	schema "github.com/RomanshkVolkov/test-api/internal/core/domain/schemas"
)

func (server Server) GetAllUsers() domain.APIResponse[[]domain.UserTableCRUD] {
	repo := repository.GetDBConnection(server.Host)
	users, err := repo.GetAllUsers()
	if err != nil {
		return domain.APIResponse[[]domain.UserTableCRUD]{
			Success: false,
			Message: domain.Message{
				En: "Error on get users",
				Es: "Error al obtener usuarios",
			},
			Error: err,
		}
	}

	return domain.APIResponse[[]domain.UserTableCRUD]{
		Success: true,
		Message: domain.Message{
			En: "Users list",
			Es: "Lista de usuarios",
		},
		Data: users,
	}
}

func (server Server) GetUserByID(id uint) domain.APIResponse[domain.EditableUser] {
	repo := repository.GetDBConnection(server.Host)
	user, err := repo.GetUserByID(id)
	if err != nil {
		return domain.APIResponse[domain.EditableUser]{
			Success: false,
			Message: domain.Message{
				En: "Error on get user",
				Es: "Error al obtener usuario",
			},
			Error: err,
		}
	}

	if user.ID == 0 {
		return repository.RecordNotFound[domain.EditableUser]()
	}

	return domain.APIResponse[domain.EditableUser]{
		Success: true,
		Message: domain.Message{
			En: "User data",
			Es: "Datos de usuario",
		},
		Data: user,
	}
}

func (server Server) CreateUser(request *domain.CreateUserRequest) domain.APIResponse[domain.User] {
	fields := schema.GenericForm[domain.CreateUserRequest]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)
	if len(failValidatedFields) > 0 {
		return SchemaFieldsError[domain.User](failValidatedFields)
	}

	repo := repository.GetDBConnection(server.Host)
	user, err := repo.CreateUser(fields.Data)
	if err != nil {
		return domain.APIResponse[domain.User]{
			Success: false,
			Message: domain.Message{
				En: "Error on create user",
				Es: "Error al crear usuario",
			},
			Error: err,
		}
	}

	return domain.APIResponse[domain.User]{
		Success: true,
		Message: domain.Message{
			En: "User created",
			Es: "Usuario creado",
		},
		Data: user,
	}
}

func (server Server) UpdateUser(request *domain.EditableUser) domain.APIResponse[domain.User] {
	fields := schema.GenericForm[domain.EditableUser]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)

	if len(failValidatedFields) > 0 {
		return SchemaFieldsError[domain.User](failValidatedFields)
	}

	repo := repository.GetDBConnection(server.Host)
	user, err := repo.UpdateUser(fields.Data)
	if err != nil {
		return domain.APIResponse[domain.User]{
			Success: false,
			Message: domain.Message{
				En: "Error on update user",
				Es: "Error al actualizar usuario",
			},
			Error: err,
		}
	}

	return domain.APIResponse[domain.User]{
		Success: true,
		Message: domain.Message{
			En: "User updated",
			Es: "Usuario actualizado",
		},
		Data: user,
	}
}

func (server Server) DeleteUser(id uint) domain.APIResponse[domain.User] {
	repo := repository.GetDBConnection(server.Host)
	err := repo.DeleteRecord(id, domain.User{})
	if err != nil {
		return domain.APIResponse[domain.User]{
			Success: false,
			Message: domain.Message{
				En: "Error on delete user",
				Es: "Error al eliminar usuario",
			},
			Error: err,
		}
	}

	return domain.APIResponse[domain.User]{
		Success: true,
		Message: domain.Message{
			En: "User deleted",
			Es: "Usuario eliminado",
		},
	}
}

func (server Server) GetUsersProfiles() domain.APIResponse[[]domain.UserProfiles] {
	repo := repository.GetDBConnection(server.Host)
	profiles, err := repo.GetUsersProfiles()
	if err != nil {
		return domain.APIResponse[[]domain.UserProfiles]{
			Success: false,
			Message: domain.Message{
				En: "Error on get users profiles",
				Es: "Error al obtener perfiles de usuarios",
			},
			Error: err,
		}
	}

	return domain.APIResponse[[]domain.UserProfiles]{
		Success: true,
		Message: domain.Message{
			En: "Users profiles",
			Es: "Perfiles de usuarios",
		},
		Data: profiles,
	}
}

func (server Server) GetKitchens() domain.APIResponse[[]domain.Kitchen] {
	repo := repository.GetDBConnection(server.Host)
	kitchens, err := repo.GetKitchens()
	if err != nil {
		return domain.APIResponse[[]domain.Kitchen]{
			Success: false,
			Message: domain.Message{
				En: "Error on get kitchens",
				Es: "Error al obtener cocinas",
			},
			Error: err,
		}
	}

	return domain.APIResponse[[]domain.Kitchen]{
		Success: true,
		Message: domain.Message{
			En: "Kitchens list",
			Es: "Lista de cocinas",
		},
		Data: kitchens,
	}
}

func (server Server) CreateKitchen(request *domain.GenericCatalog) domain.APIResponse[domain.Kitchen] {
	fields := schema.GenericForm[domain.GenericCatalog]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)

	if len(failValidatedFields) > 0 {
		return SchemaFieldsError[domain.Kitchen](failValidatedFields)
	}

	repo := repository.GetDBConnection(server.Host)
	createdKitchen, err := repo.CreateKitchen(fields.Data)
	if err != nil {
		return domain.APIResponse[domain.Kitchen]{
			Success: false,
			Message: domain.Message{
				En: "Error on create kitchen",
				Es: "Error al crear cocina",
			},
			Error: err,
		}
	}

	return domain.APIResponse[domain.Kitchen]{
		Success: true,
		Message: domain.Message{
			En: "Kitchen created",
			Es: "Cocina creada",
		},
		Data: createdKitchen,
	}
}

func (server Server) GetShifts() domain.APIResponse[[]domain.Shift] {
	repo := repository.GetDBConnection(server.Host)
	shifts, err := repo.GetShifts()
	if err != nil {
		return domain.APIResponse[[]domain.Shift]{
			Success: false,
			Message: domain.Message{
				En: "Error on get shifts",
				Es: "Error al obtener turnos",
			},
			Error: err,
		}
	}

	return domain.APIResponse[[]domain.Shift]{
		Success: true,
		Message: domain.Message{
			En: "Shifts list",
			Es: "Lista de turnos",
		},
		Data: shifts,
	}
}

func (server Server) CreateShift(request *domain.GenericCatalog) domain.APIResponse[domain.Shift] {
	fields := schema.GenericForm[domain.GenericCatalog]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)

	if len(failValidatedFields) > 0 {
		return SchemaFieldsError[domain.Shift](failValidatedFields)
	}

	repo := repository.GetDBConnection(server.Host)
	createdShift, err := repo.CreateShift(fields.Data)
	if err != nil {
		return domain.APIResponse[domain.Shift]{
			Success: false,
			Message: domain.Message{
				En: "Error on create shift",
				Es: "Error al crear turno",
			},
			Error: err,
		}
	}

	return domain.APIResponse[domain.Shift]{
		Success: true,
		Message: domain.Message{
			En: "Shift created",
			Es: "Turno creado",
		},
		Data: createdShift,
	}
}
