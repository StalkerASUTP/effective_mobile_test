package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config представляет основную конфигурацию приложения
// @Description Основная конфигурационная структура приложения
type Config struct {
	// DSN (Data Source Name) для подключения к базе данных
	// Пример: "host=localhost user=postgres password=secret dbname=efmob port=5432 sslmode=disable"
	Dsn string `json:"dsn"`
}

// LoadConfig загружает конфигурацию из .env файла и переменных окружения
// @Summary Загрузить конфигурацию
// @Description Загружает конфигурацию из .env файла (если существует) и переменных окружения
// @Tags Конфигурация
// @Produce json
// @Success {object} Config
// @Failure {string} string "Ошибка при загрузке .env файла"

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}
	return &Config{
		Dsn: os.Getenv("DSN"),
	}
}