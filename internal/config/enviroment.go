package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	Port                      string
	AppName                   string
	AppEnv                    string
	DBPGHost                  string
	DBPGPort                  string
	DBPGUsername              string
	DBPGPassword              string
	DBPGName                  string
	DBPGSSLMode               string
	DBPGSchema                string
	PWDefaultAdmin            string
	JWTSecretKey              string
	JWTRefreshTokenSecretKey  string
	JWTExpiredHours           int
	JWTRfreshTokenExpiredDays int
	FormatDateTime            string
	UserMaxDevice             int
	ApmUrl                    string
	ApmNameService            string
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback
	}

	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Failed convert value to int used default value, error: %v", err)
		return fallback
	}

	return valueInt
}

func LoadEnvironment() *Environment {
	err := godotenv.Load()
	if err != nil {
		log.Println("Info file .env is not found, used default environment variable")
	}

	return &Environment{
		Port:                      getEnv("PORT", "8000"),
		AppName:                   getEnv("APP_NAME", "rizat-service-blog"),
		AppEnv:                    getEnv("APP_ENV", "development"),
		DBPGHost:                  getEnv("DB_PG_HOST", "localhost"),
		DBPGPort:                  getEnv("DB_PG_PORT", "5432"),
		DBPGUsername:              getEnv("DB_PG_USERNAME", "postgres"),
		DBPGPassword:              getEnv("DB_PG_PASSWORD", ""),
		DBPGName:                  getEnv("DB_PG_NAME", "my_db"),
		DBPGSSLMode:               getEnv("DB_PG_SSLMODE", "disable"),
		DBPGSchema:                getEnv("DB_PG_SCHEMA", "public"),
		PWDefaultAdmin:            getEnv("PW_DEFAULT_ADMIN", "admin"),
		JWTSecretKey:              getEnv("JWT_SECRET_KEY", ""),
		JWTRefreshTokenSecretKey:  getEnv("JWT_RERESH_TOKEN_SECRET_KEY", ""),
		JWTExpiredHours:           getEnvInt("JWT_EXPIRED_HOURS", 1),
		JWTRfreshTokenExpiredDays: getEnvInt("JWT_REFRESH_TOKEN_EXPIRED_DAYS", 7),
		FormatDateTime:            getEnv("FORMAT_DATE_TIME", "2006-01-02 15:04:05"),
		UserMaxDevice:             getEnvInt("USER_MAX_DEVICE", 1),
		ApmUrl:                    getEnv("APM_URL", "localhost:3000"),
		ApmNameService:            getEnv("APM_NAME_SERVICE", "template-go-gin"),
	}
}
