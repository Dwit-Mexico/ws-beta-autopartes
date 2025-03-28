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
	// db.Exec("DROP TABLE IF EXISTS users_has_kitchens")
	// db.Exec("DROP TABLE IF EXISTS shifts")
	// db.Exec("DROP TABLE IF EXISTS kitchens")
	// db.Exec("DROP TABLE IF EXISTS users")
	// db.Exec("DROP TABLE IF EXISTS user_profiles")
	// db.Exec("DROP TABLE IF EXISTS permissions")
	// db.Exec("DROP TABLE IF EXISTS devs")
	// db.Exec("DROP TABLE IF EXISTS hosting_centers")
	// db.Exec("DROP TABLE IF EXISTS detail_document_tables")
	// db.Exec("DROP TABLE IF EXISTS document_tables")
	// db.Exec("DROP TABLE IF EXISTS detail_documents")
	// db.Exec("DROP TABLE IF EXISTS documents")

	SeedProfiles(db)
	SeedPermissions(db)
	SeedShifts(db)
	SeedUsers(db)
	SeedKitchens(db)
	SeedDevAuthorizedIPAddress(db)
	SeedDocumentsAndReports(db)
	MigrateProcedures(db)

	SeedExampleReports(db)

	// this function is used to create the tables of the document definition
	MigrateDocumentTables(db)

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

	rootProfile := domain.UserProfiles{}
	db.Model(&domain.UserProfiles{}).Where("slug = ?", "root").First(&rootProfile)

	users := []domain.User{
		{
			UserData: domain.UserData{
				Username: "dwitmx",
				Email:    "sistemas@dwitmexico.com",
				Name:     "Dwit México",
				IsActive: true,
			},
			ProfileID: rootProfile.ID,
			Password:  "password",
		},
		{
			UserData: domain.UserData{
				Username: "romanshkvolkov",
				Email:    "jose@guz-studio.dev",
				Name:     "Romanshk Volkov",
				IsActive: true,
			},
			ProfileID: rootProfile.ID,
			Password:  "password",
		},
	}

	for _, user := range users {
		hashedPassword, _ := HashPassword(user.Password)
		user.Password = hashedPassword

		db.Create(&user)
	}
}

func SeedProfiles(db *gorm.DB) {
	AutoMigrateTable(db, &domain.UserProfiles{})

	var currentRows int64
	db.Model(&domain.UserProfiles{}).Count(&currentRows)

	if currentRows > 0 {
		return
	}

	profiles := []*domain.UserProfiles{
		{
			Name: "Super Admin",
			Slug: "root",
		},
		{
			Name: "Administrador",
			Slug: "admin",
		},
		{
			Name: "Cliente",
			Slug: "customer",
		},
	}

	for _, profile := range profiles {
		db.Create(&profile)
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

func SeedKitchens(db *gorm.DB) {
	AutoMigrateTable(db, &domain.Kitchen{})
	AutoMigrateTable(db, &domain.UsersHasKitchens{})
}

func SeedShifts(db *gorm.DB) {
	AutoMigrateTable(db, &domain.Shift{})
}

func SeedPermissions(db *gorm.DB) {
	AutoMigrateTable(db, &domain.Permission{})
	AutoMigrateTable(db, &domain.ProfilesHasPermissions{})

}

func SeedDocumentsAndReports(db *gorm.DB) {
	AutoMigrateTable(db, &domain.Document{})
	AutoMigrateTable(db, &domain.DetailDocument{})
	AutoMigrateTable(db, &domain.Report{})
	AutoMigrateTable(db, &domain.ChartReport{})
	AutoMigrateTable(db, &domain.ChartLine{})

	var currentRows int64
	db.Model(&domain.Document{}).Count(&currentRows)

	if currentRows > 0 {
		return
	}

}

func MigrateDocumentTables(db *gorm.DB) {
	documents := []domain.Document{}

	currentRows := int64(0)
	db.Model(&domain.Document{}).Count(&currentRows)

	if currentRows == 0 {
		return
	}

	db.Find(&documents)

	for _, document := range documents {
		table := document.Table

		isExist := ExistTable(db, table)
		if isExist {
			continue
		}

		id := document.ID
		_, err := ExecuteProcedureSQLServer(db, "sp_CreateTableToDocument", id)
		if err != nil {
			fmt.Println("Error when creating table: %w", err)
		}

	}
}

func SeedExampleReports(db *gorm.DB) {
	currentRows := int64(0)
	db.Model(&domain.Report{}).Count(&currentRows)

	if currentRows > 0 {
		return
	}

	report := domain.Report{
		Name:            "Example Report",
		StoredProcedure: "sp_example_report",
	}

	chart := domain.ChartReport{
		Name:            "Example Chart",
		StoredProcedure: "sp_example_chart",
		XAxisKey:        "date",
	}

	db.Create(&report)

	chart.ReportID = 1
	db.Create(&chart)

	chartLine := domain.ChartLine{
		LineKey: "name_field_in_procedure",
		Stroke:  "#8884d8",
		Name:    "Example Line",
	}

	chartLine.ChartID = chart.ID

	db.Create(&chartLine)
}
