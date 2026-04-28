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

func GetJWTSecrets() [][]byte {
	current := os.Getenv("JWT_SECRET_CURRENT")
	old := os.Getenv("JWT_SECRET_OLD")

	var secrets [][]byte

	if current != "" {
		secrets = append(secrets, []byte(current))
	}

	if old != "" {
		secrets = append(secrets, []byte(old))
	}

	// if empty
	if len(secrets) == 0 {
		log.Fatal("no JWT secrets configured")
	}

	return secrets
}
