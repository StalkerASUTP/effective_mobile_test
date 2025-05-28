package db

import (
	"go/ef-mob-api/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Db представляет обёртку над подключением к базе данных через GORM
// @Description Основная структура для работы с базой данных
type Db struct {
	*gorm.DB // GORM подключение к базе данных
}

// NewDb создает новое подключение к PostgreSQL

func NewDb(conf *configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(conf.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{db}
}
