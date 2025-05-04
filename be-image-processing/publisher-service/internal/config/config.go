package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	ServiceName      string         `json:"serviceName"`
	ServicePort      string         `json:"servicePort"`
	GinMode          string         `json:"ginMode"`
	Environment      Environment    `json:"environment"`
	DatabaseConfig   DatabaseConfig `json:"databaseConfig"`
	CorsAllowOrigins []string       `json:"corsAllowOrigins"`
	RabbitMQConfig   RabbitMQConfig `json:"rabbitMQConfig"`
}

const logTagConifg = "[Init Config]"

var config *Config

func Init() {
	conf := Config{
		ServiceName: os.Getenv("SERVICE_NAME"),
		ServicePort: os.Getenv("SERVICE_PORT"),
		GinMode:     os.Getenv("GIN_MODE"),
		DatabaseConfig: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
		RabbitMQConfig: RabbitMQConfig{
			RabbitMQUrl: os.Getenv("RABBITMQ_URL"),
		},
	}

	if conf.ServiceName == "" {
		slog.Error(fmt.Sprintf("%s service name should not be empty", logTagConifg))
	}

	if conf.ServicePort == "" {
		slog.Error(fmt.Sprintf("%s service port should not be empty", logTagConifg))
	}

	if conf.ServicePort == "" {
		slog.Error(fmt.Sprintf("%s service port should not be empty", logTagConifg))
	}

	if conf.GinMode != "debug" && conf.GinMode != "release" {
		log.Fatalf("%s gin mode must be debug or release, found: %s", logTagConifg, conf.GinMode)
	}

	if conf.RabbitMQConfig.RabbitMQUrl == "" {
		log.Fatalf("%s rabbitMQ url cannot be empty", logTagConifg)
	}

	envString := os.Getenv("ENVIRONMENT")
	if envString != "dev" && envString != "prod" {
		slog.Error(fmt.Sprintf("%s environment must be eiher dev or prod, found: %s", logTagConifg, envString))
	}

	conf.Environment = Environment(envString)

	corsOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	conf.CorsAllowOrigins = strings.Split(corsOrigins, "|")
	config = &conf
}

func Get() (conf *Config) {
	conf = config
	return
}
