package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath    string
	Port      string
	UploadDir string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		DBPath:    getEnv("DB_PATH", "./store.db"),
		Port:      getEnv("SERVER_PORT", "8080"),
		UploadDir: getEnv("UPLOAD_DIR", "./uploads"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("переменная %s не найдена, используем дефолт: %s", key, fallback)
	return fallback
}
