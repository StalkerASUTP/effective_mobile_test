package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Dsn string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}
	return &Config{
		Dsn: os.Getenv("DSN"),
	}
}
