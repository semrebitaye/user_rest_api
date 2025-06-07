package initializer

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariable() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
