package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn_ms = GetEnv("DB_MYSQL_STRING_CONECTION")
var DBMS *gorm.DB

func DBConnectionMySQL() {
	db, err := gorm.Open(mysql.Open(dsn_postgres), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	DBMS = db
}
