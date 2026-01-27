package config

import "os"

var (
	PostgresURL = getEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/authdb?sslmode=disable")
	JwtSecret   = []byte(getEnv("JWT_SECRET", "supersecret"))
	Issuer      = "auth-service"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
