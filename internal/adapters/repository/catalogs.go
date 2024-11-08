package repository

import (
	"errors"

	"github.com/RomanshkVolkov/test-api/internal/core/domain"
)

func (database *DSNSource) GetKitchenByID(id uint) (domain.Kitchen, error) {
	kitchen := domain.Kitchen{}
	database.DB.Model(&domain.Kitchen{}).Where("id = ?", id).First(&kitchen)

	if kitchen.ID == 0 {
		return domain.Kitchen{}, nil
	}

	return kitchen, nil
}

func (database *DSNSource) GetShiftByID(id uint) (domain.Shift, error) {
	shift := domain.Shift{}
	database.DB.Model(&domain.Shift{}).Where("id = ?", id).First(&shift)

	if shift.ID == 0 {
		return domain.Shift{}, nil
	}

	return shift, nil
}

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

var MAPPED_AVAILABLE_CATALOGS = map[string]interface{}{
	"kitchen": domain.Kitchen{},
	"shift":   domain.Shift{},
}

func GetCatalogTable(name string) interface{} {
	return MAPPED_AVAILABLE_CATALOGS[name]
}
