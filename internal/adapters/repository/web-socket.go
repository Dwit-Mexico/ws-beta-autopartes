package repository

import "github.com/RomanshkVolkov/test-api/internal/core/domain"

func GetWebSocketWarehouses() ([]domain.WebSocketWarehouses, error) {
	database := domain.DBBetaAutopartes
	dns := GetDBConnection(database)
	db := dns.DB
	var warehouses []domain.WebSocketWarehouses
	err := db.Find(&warehouses).Error
	return warehouses, err
}
