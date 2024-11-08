package repository

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// var dsn_sql_server = GetEnv("DB_SQLSERVER_STRING_CONECTION")

var DSNList, _ = GetDSNList()
var DBSQLServer []DSNSource

func DBConnection() {
	for _, dsn := range DSNList {
		db, err := gorm.Open(sqlserver.Open(dsn.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err != nil {
			panic("failed to connect database")
		}

		DBSQLServer = append(DBSQLServer, DSNSource{
			Name: dsn.Name,
			DB:   db,
		})

		RunSeeds(db)
	}

}

func GetDBConnection(subdomain string) DSNSource {
	fmt.Println("Subdomain: ", subdomain)
	// pending replace to read subdomain and match with DBSQLServer list
	authorizedHost := MAPPED_AUTHORIZED_DOMAINS[subdomain]
	for _, db := range DBSQLServer {
		if db.Name == authorizedHost {
			fmt.Println("DB: ", db.Name)
			return db
		}
	}

	return DSNSource{}
}

type StoredProcedureParams struct {
	Procedure string
	Params    []interface{}
}

func ExecuteProcedureSQLServer(db *gorm.DB, procedure string, args ...interface{}) (*sql.Rows, error) {
	query := fmt.Sprintf("EXEC %s", procedure)
	for range args {
		query += " ?,"
	}

	query = query[:len(query)-1]

	result := db.Exec(query, args...)
	if result.Error != nil {
		return nil, result.Error
	}

	rows, err := result.Rows()
	if err != nil {
		return nil, err
	}

	return rows, nil

}
