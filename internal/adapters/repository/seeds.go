package repository

import (
	"fmt"
	"time"

	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	"gorm.io/gorm"
)

func PrintSeedAction(nameTable string, action string) {
	fmt.Println("Seeding table: " + nameTable + " " + action + " Success")
}

func AutoMigrateTable(db *gorm.DB, table interface{}) {
	fmt.Println("AutoMigrateTable")
	isInitialized := db.Migrator().HasTable(&table)
	if !isInitialized {
		db.AutoMigrate(table)
		PrintSeedAction("Shifts", "Create")
	}
}

func RunSeeds(db *gorm.DB) {
	startTimePoint := time.Now().UTC()
	fmt.Println("====================================================================================")
	fmt.Println("Operation run on database", db.Name())
	fmt.Println("Start operation RunSeeds Seeding tables...")

	SeedUsers(db)
	SeedDevAuthorizedIPAddress(db)
	latency := time.Since(startTimePoint)
	fmt.Println("RunSeeds end operation " + latency.String())
	fmt.Println("====================================================================================")
}

func SeedUsers(db *gorm.DB) {
	AutoMigrateTable(db, &domain.User{})

	var currentRows int64
	db.Model(&domain.User{}).Count(&currentRows)

	if currentRows > 0 {
		return
	}

	users := []domain.User{
		{
			UserData: domain.UserData{
				Username: "dwitmx",
				Email:    "sistemas@dwitmexico.com",
				Name:     "Dwit México",
				IsActive: true,
			},
			Password: "password",
		},
		{
			UserData: domain.UserData{
				Username: "romanshkvolkov",
				Email:    "jose@guz-studio.dev",
				Name:     "Romanshk Volkov",
				IsActive: true,
			},
			Password: "password",
		},
	}

	for _, user := range users {
		hashedPassword, _ := HashPassword(user.Password)
		user.Password = hashedPassword

		db.Create(&user)
	}
}

func SeedDevAuthorizedIPAddress(db *gorm.DB) {
	AutoMigrateTable(db, &domain.Dev{})

	var currentRows int64
	db.Model(&domain.Dev{}).Count(&currentRows)

	if currentRows > 0 {
		return
	}

	devs := []*domain.Dev{
		{
			IP:  "172.18.0.1",
			Tag: "docker local",
		},
	}

	for _, dev := range devs {
		db.Create(&dev)
	}
}
