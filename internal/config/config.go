package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	SECRET_KEY            = "your-secret-key"
	TOKEN_EXPIRATION_TIME = 1 * time.Hour
	TOKEN_ISSUER          = "trail-data-service"
	TOKEN_SUBJECT         = "user-auth"

	DUPLICATE_TRAIL_RADIUS_KM = 25.0
)

var readFile = os.ReadFile

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Port    string `yaml:"port"`
	GinMode string `yaml:"gin_mode"`
}

func NewConfig() (*Config, error) {
	var cfg *Config

	yamlFile, err := readFile("config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
