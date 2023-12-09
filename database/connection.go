package database

import (
	"boilerplate/backend/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Connection *gorm.DB

func Connect() {

	//connect to database
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	var err error
	Connection, err = gorm.Open(mysql.Open(dbUser+":"+dbPass+"@tcp(localhost:3306)/"+dbName+"?charset=utf8mb4"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	//table creation
	Connection.AutoMigrate(&models.User{})
}
