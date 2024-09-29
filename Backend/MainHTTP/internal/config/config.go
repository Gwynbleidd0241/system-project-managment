package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Logs_path              string                 `yaml:"logs_path" env-required:"true"`
	SecretKey              string                 `yaml:"secret_key" env-required:"true"`
	AuthgRPCConfig         AuthgRPCConfig         `yaml:"auth_grpc"`
	TaskgRPCConfig         TaskgRPCConfig         `yaml:"task_grpc"`
	NotificationgRPCConfig NotificationgRPCConfig `yaml:"notification_grpc"`
	MainHTTPConfig         MainHTTPConfig         `yaml:"main_http"`
}

type AuthgRPCConfig struct {
	Port string `yaml:"port" env-required:"true"`
}

type TaskgRPCConfig struct {
	Port string `yaml:"port" env-required:"true"`
}

type NotificationgRPCConfig struct {
	Port string `yaml:"port" env-required:"true"`
}

type MainHTTPConfig struct {
	Port    string        `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" yaml-default:"5s"`
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
