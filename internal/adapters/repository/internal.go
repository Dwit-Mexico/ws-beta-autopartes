package repository

import "github.com/RomanshkVolkov/test-api/internal/core/domain"

func (database *DSNSource) PermissionsSynchronization(routes []domain.WebCatalogRoutes) error {
	for _, route := range routes {
		permission := domain.Permission{}

		database.DB.Where("path = ?", route.Path).First(&permission)

		if permission.ID == 0 {
			err := database.DB.Create(&domain.Permission{
				Name: route.Name,
				Path: route.Path,
			}).Error

			if err != nil {
				return err
			}
		} else {
			err := database.DB.Model(&domain.Permission{}).Where("id = ?", permission.ID).Update("name", route.Name).Error
			if err != nil {
				return err
			}
		}

	}
	return nil
}
