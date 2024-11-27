package repository

import "github.com/RomanshkVolkov/test-api/internal/core/domain"

func (database *DSNSource) PermissionsSynchronization(routes []domain.WebCatalogRoutes) error {
	// truncate the table
	tx := database.DB.Begin()

	for _, route := range routes {
		permission := domain.Permission{}

		tx.Where("path = ?", route.Path).First(&permission)

		if permission.ID == 0 {
			err := tx.Create(&domain.Permission{
				Name:   route.Name,
				Path:   route.Path,
				Status: true,
			}).Error
			if err != nil {
				return err
			}

		} else {
			err := tx.Model(&domain.Permission{}).Where("id = ?", permission.ID).Update("name", route.Name).Error
			if err != nil {
				return err
			}
		}

	}
	tx.Commit()
	return nil
}
