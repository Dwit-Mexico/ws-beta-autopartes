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
	SeedShifts(db)
	SeedUsers(db)
	SeedKitchens(db)
	SeedDevAuthorizedIPAddress(db)
	SeedHostingCenters(db)
	SeedDocumentsAndReports(db)
	MigrateProcedures(db)

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

	var currentRows int64
	db.Model(&domain.Permission{}).Count(&currentRows)

	if currentRows > 0 {
		return
	}

	// permissions := []*domain.Permission{
	// 	{
	// 		Name: "Dashboard",
	// 		Path: "/dashboard",
	// 	},
	// 	{
	// 		Name: "Configuración",
	// 		Path: "/dashboard/settings",
	// 	},
	// 	{
	// 		Name: "Usuarios",
	// 		Path: "/dashboard/settings/users",
	// 	},
	// 	{
	// 		Name: "Platillos",
	// 		Path: "/dashboard/dishes",
	// 	},
	// 	{
	// 		Name: "Administración",
	// 		Path: "/dashboard/managment",
	// 	},
	// }
}

func SeedHostingCenters(db *gorm.DB) {
	AutoMigrateTable(db, &domain.HostingCenter{})
	defaultHostingCenter := domain.HostingCenter{
		Name:        "Default",
		CompanyName: "Default",
	}

	var currentRows int64
	db.Model(&domain.HostingCenter{}).Count(&currentRows)

	if currentRows > 0 {
		return
	}
	db.Create(&defaultHostingCenter)
}

func SeedDocumentsAndReports(db *gorm.DB) {
	AutoMigrateTable(db, &domain.Document{})
	AutoMigrateTable(db, &domain.DetailDocument{})
	AutoMigrateTable(db, &domain.Reports{})
	AutoMigrateTable(db, &domain.ChartReports{})

	var currentRows int64
	db.Model(&domain.Document{}).Count(&currentRows)

	if currentRows > 0 {
		return
	}

	documents := []*domain.DocumentWithDetails{
		{
			Document: domain.Document{
				Name:  "Alimentos",
				Table: "foods",
			},
			Details: []domain.DetailDocument{
				{
					Field:       "code",
					TypeField:   "NVARCHAR(20)",
					DocumentKey: "CLAVE",
				},
				{
					Field:       "description",
					TypeField:   "NVARCHAR(300)",
					DocumentKey: "DESCRIPCIÓN",
				},
				{
					Field:       "category",
					TypeField:   "NVARCHAR(100)",
					DocumentKey: "CATEGORÍA",
				},
				{
					Field:       "purchase_unit",
					TypeField:   "NVARCHAR(20)",
					DocumentKey: "UNIDAD DE COMPRA",
				},
			},
		},
		{
			Document: domain.Document{
				Name:  "Platillos",
				Table: "dishes",
			},
			Details: []domain.DetailDocument{
				{
					Field:       "code",
					TypeField:   "NVARCHAR(20)",
					DocumentKey: "CLAVE",
				},
				{
					Field:       "description",
					TypeField:   "NVARCHAR(300)",
					DocumentKey: "DESCRIPCIÓN",
				},
				{
					Field:       "category",
					TypeField:   "NVARCHAR(100)",
					DocumentKey: "CATEGORIA",
				},
				{
					Field:       "unit",
					TypeField:   "NVARCHAR(20)",
					DocumentKey: "UNIDAD",
				},
			},
		},
		{
			Document: domain.Document{
				Name:  "Costos",
				Table: "costs",
			},
			Details: []domain.DetailDocument{
				// CLAVE	DESCRIPCIÓN	COSTO	Fecha
				{
					Field:       "code",
					TypeField:   "NVARCHAR(20)",
					DocumentKey: "CLAVE",
				}, {
					Field:       "description",
					TypeField:   "NVARCHAR(300)",
					DocumentKey: "DESCRIPCIÓN",
				},
				{
					Field:       "cost",
					TypeField:   "DECIMAL(10,2)",
					DocumentKey: "COSTO",
				},
				{
					Field:       "date",
					TypeField:   "DATETIME",
					DocumentKey: "FECHA",
				},
			},
		},
		{
			Document: domain.Document{
				Name:  "Menús",
				Table: "menus",
			},
			Details: []domain.DetailDocument{
				// CLAVE	DESCRIPCIÓN	TURNO	COCINA	DESDE	HASTA
				{
					Field:       "code",
					TypeField:   "NVARCHAR(20)",
					DocumentKey: "CLAVE",
				},
				{
					Field:       "description",
					TypeField:   "NVARCHAR(300)",
					DocumentKey: "DESCRIPCIÓN",
				},
				{
					Field:       "shift",
					TypeField:   "NVARCHAR(100)",
					DocumentKey: "TURNO",
				},
				{
					Field:       "kitchen",
					TypeField:   "NVARCHAR(100)",
					DocumentKey: "COCINA",
				},
				{
					Field:       "date_from",
					TypeField:   "DATETIME",
					DocumentKey: "DESDE",
				},
				{
					Field:       "date_to",
					TypeField:   "DATETIME",
					DocumentKey: "HASTA",
				},
			},
		},
		{
			Document: domain.Document{
				Name:  "Mermas",
				Table: "foos_waste",
			},
			Details: []domain.DetailDocument{
				// CLAVE	DESCRIPCIÓN	CANTIDAD	UNIDAD	DÍA	TURNO	COCINA
				{
					Field:       "code",
					TypeField:   "NVARCHAR(20)",
					DocumentKey: "CLAVE",
				},
				{
					Field:       "description",
					TypeField:   "NVARCHAR(300)",
					DocumentKey: "DESCRIPCIÓN",
				},
				{
					Field:       "quantity",
					TypeField:   "DECIMAL(10,2)",
					DocumentKey: "CANTIDAD",
				},
				{
					Field:       "unit",
					TypeField:   "NVARCHAR(20)",
					DocumentKey: "UNIDAD",
				},
				{
					Field:       "day",
					TypeField:   "DATETIME",
					DocumentKey: "DÍA",
				},
				{
					Field:       "shift",
					TypeField:   "NVARCHAR(100)",
					DocumentKey: "TURNO",
				},
			},
		},
		{
			Document: domain.Document{
				Name:  "Consumo",
				Table: "consumption",
			},
			Details: []domain.DetailDocument{
				// CLAVE	DESCRIPCIÓN	CANTIDAD	UNIDAD	TURNO	COCINA
				{
					Field:       "code",
					TypeField:   "NVARCHAR(20)",
					DocumentKey: "CLAVE",
				},
				{
					Field:       "description",
					TypeField:   "NVARCHAR(300)",
					DocumentKey: "DESCRIPCIÓN",
				},
				{
					Field:       "quantity",
					TypeField:   "DECIMAL(10,2)",
					DocumentKey: "CANTIDAD",
				},
				{
					Field:       "unit",
					TypeField:   "NVARCHAR(20)",
					DocumentKey: "UNIDAD",
				},
				{
					Field:       "shift",
					TypeField:   "NVARCHAR(100)",
					DocumentKey: "TURNO",
				},
				{
					Field:       "kitchen",
					TypeField:   "NVARCHAR(100)",
					DocumentKey: "COCINA",
				},
			},
		},
		{
			Document: domain.Document{
				Name:  "Asistencia",
				Table: "attendance",
			},
			Details: []domain.DetailDocument{
				// COCINA TURNO PORCENTAJE
				{
					Field:       "kitchen",
					TypeField:   "NVARCHAR(100)",
					DocumentKey: "COCINA",
				},
				{
					Field:       "shift",
					TypeField:   "NVARCHAR(100)",
					DocumentKey: "TURNO",
				},
				{
					Field:       "percentage",
					TypeField:   "DECIMAL(10,2)",
					DocumentKey: "PORCENTAJE",
				},
			},
		},
	}

	for _, document := range documents {
		db.Create(&document.Document)
		for _, detail := range document.Details {
			detail.DocumentID = document.Document.ID
			detail.DocumentKey = RemoveSpaces(RemoveAccents(detail.DocumentKey))
			db.Create(&detail)
		}
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
