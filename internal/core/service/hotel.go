package service

import (
	"github.com/RomanshkVolkov/test-api/internal/adapters/repository"
	"github.com/RomanshkVolkov/test-api/internal/core/domain"
)

func (server Server) GetCurrentHostingCenter() domain.APIResponse[domain.HostingCenter] {
	repo := repository.GetDBConnection(server.Host)
	hotel, err := repo.GetCurrentHostingCenter()

	if err != nil {
		return repository.RecordNotFound[domain.HostingCenter]()
	}

	return domain.APIResponse[domain.HostingCenter]{
		Success: true,
		Message: domain.Message{
			En: "Current hosting center",
			Es: "Centro de alojamiento actual",
		},
		Data: hotel,
	}
}

func (server Server) UpdateHostingCenter(request *domain.HostingCenter) domain.APIResponse[domain.HostingCenter] {
	repo := repository.GetDBConnection(server.Host)
	hotel, err := repo.UpdateHostingCenter(domain.HostingCenter{
		ID:          request.ID,
		Name:        request.Name,
		CompanyName: request.CompanyName,
	})

	if hotel.ID == 0 {
		return repository.RecordNotFound[domain.HostingCenter]()
	}

	if err != nil {
		return repository.HandleDatabaseError[domain.HostingCenter](err, domain.Message{En: "Error on update hosting center", Es: "Error al actualizar centro de alojamiento"})
	}

	return domain.APIResponse[domain.HostingCenter]{
		Success: true,
		Message: domain.Message{
			En: "Hosting center updated",
			Es: "Centro de alojamiento actualizado",
		},
		Data: hotel,
	}

}
