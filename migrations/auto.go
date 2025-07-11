package main

import (
	"go/ef-mob-api/person"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AutoMigrate выполняет автоматические миграции базы данных
// @Title Database Migrations
// @Description Запускает автоматические миграции для всех моделей приложения
// @Tags Database
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Migrator().AutoMigrate(&person.Person{})

}
