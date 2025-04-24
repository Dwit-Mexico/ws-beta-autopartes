package repository

import (
	"errors"
)

func (database *DSNSource) UpdateGenericCatalog(id uint, table interface{}, name string) error {
	result := database.DB.Model(table).Where("id = ?", id).Update("name", CapitalizeAll(name))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}

func (database *DSNSource) DeleteRecord(id uint, table interface{}) error {
	var item interface{}
	result := database.DB.Model(table).Where("id = ?", id).Delete(&item)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}
