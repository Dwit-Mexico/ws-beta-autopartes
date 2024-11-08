package repository

import "github.com/RomanshkVolkov/test-api/internal/core/domain"

func (database *DSNSource) GetCurrentHostingCenter() (domain.HostingCenter, error) {
	hotel := domain.HostingCenter{}
	database.DB.Model(&domain.HostingCenter{}).Where("id = ?", 1).First(&hotel)

	if hotel.ID == 0 {
		return domain.HostingCenter{}, nil
	}

	return hotel, nil
}

func (database *DSNSource) UpdateHostingCenter(data domain.HostingCenter) (domain.HostingCenter, error) {
	hotel := domain.HostingCenter{}
	result := database.DB.Model(&domain.HostingCenter{}).Where("id = ?", 1).Updates(data)

	if result.Error != nil {
		return domain.HostingCenter{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domain.HostingCenter{}, nil
	}

	database.DB.Model(&domain.HostingCenter{}).Where("id = ?", 1).First(&hotel)

	return hotel, nil
}
