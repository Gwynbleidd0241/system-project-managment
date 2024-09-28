package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env          string `yaml:"env" yaml-default:"local"`
	Storage_path string `yaml:"storage_path" env-required:"true"`
	Logs_path    string `yaml:"logs_path" env-required:"true"`
}

func LoadConfig() (*Config, error) {

	errorStatement := "config.go"

	configPath := "./config/config.yaml"

	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("%s: %s", errorStatement, err)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, fmt.Errorf("%s: %s", errorStatement, err)
	}

	return &config, nil
}
