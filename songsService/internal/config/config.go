package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `env:"ENV" env-default:"local"`
	Database   Database
	HTTPServer HTTPServer
	Cache      Cache
}

type HTTPServer struct {
	Address     string        `env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `env:"HTTP_SERVER_TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"5s"`
}

type Database struct {
	Host     string `env:"DB_HOST" env-default:"localhost:5432"`
	Name     string `env:"DB_NAME" env-required:"true"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
}

type Cache struct {
	Address string `env:"CACHE_ADDRESS" env-default:"localhost:6379"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if err := godotenv.Load(configPath); err != nil {
		log.Fatalf("Error loading .env file from %s: %v", configPath, err)
	}

	var config Config

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalf("Error reading environment variables: %s", err.Error())
	}

	return &config
}
