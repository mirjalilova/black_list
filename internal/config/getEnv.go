package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	API_GATEWAY    string
	AUTH_PORT      string
	BLACKLIST_PORT string

	DB_HOST     string
	DB_PORT     int
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	LOG_PATH string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.API_GATEWAY = cast.ToString(coalesce("API_GATEWAY", ":8070"))
	config.AUTH_PORT = cast.ToString(coalesce("AUTH_PORT", ":8070"))
	config.BLACKLIST_PORT = cast.ToString(coalesce("BLACKLIST_PORT", ":8090"))

	config.DB_HOST = cast.ToString(coalesce("DB_HOST", "localhost"))
	config.DB_PORT = cast.ToInt(coalesce("DB_PORT", 5432))
	config.DB_USER = cast.ToString(coalesce("DB_USER", "n10"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "12345"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", "n10"))

	config.LOG_PATH = cast.ToString(coalesce("LOG_PATH", "logs/info.log"))

	return config
}

func coalesce(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
