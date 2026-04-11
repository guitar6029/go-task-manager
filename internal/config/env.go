package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load(".env")

	env := os.Getenv("APP_ENV")

	switch env {
	case "local":
		log.Println("Running in LOCAL mode")
	case "dev":
		log.Println("Running in DEV mode")
	default:
		log.Println("Running with system env variables")
	}
}
