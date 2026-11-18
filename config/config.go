package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	AppEnv         string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	PresensiAPIKey string
}

var AppConfig *Config

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment variable sistem.")
	}

	cfg := &Config{
		AppPort:        getEnv("APP_PORT", "3000"),
		AppEnv:         getEnv("APP_ENV", "development"),
		DBHost:         getEnv("DB_HOST", "127.0.0.1"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "digitalpiket"),
		DBPassword:     getEnv("DB_PASSWORD", "piketsmkpenus#"),
		DBName:         getEnv("DB_NAME", "digitalpiket"),
		JWTSecret:      getEnv("JWT_SECRET", "epiket_smkpluspenus_super_secret_jwt_2026"),
		PresensiAPIKey: getEnv("PRESENSI_API_KEY", "penus2026secret"),
	}

	AppConfig = cfg
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
