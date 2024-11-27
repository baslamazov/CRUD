package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"environment" env-required:"true"`
	Database   `yaml:"database" env-required:"true"`
	HTTPServer `yaml:"http_server" env-required:"true"`
	GRPC       GRPCConfig    `yaml:"grpc"`
	TokenTTL   time.Duration `yaml:"token_ttl" env-required:"true"`
}
type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Port        string        `yaml:"port" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}
type Database struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Default().Print("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("нет пути")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Не существует конфига %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Не читается конфиг %s", err)
	}
	return &config
}
