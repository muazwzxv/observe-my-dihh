package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	RedisHost  string
	RedisPort  string
	KafkaHost  string
	KafkaPort  string

	OTELServiceName      string
	OTELExporterEndpoint string
	OTELExporterProtocol string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "otel_api_demo_user"),
		DBPassword: getEnv("DB_PASSWORD", "otel_api_demo_pass"),
		DBName:     getEnv("DB_NAME", "otel_api_demo_db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		KafkaHost:  getEnv("KAFKA_HOST", "localhost"),
		KafkaPort:  getEnv("KAFKA_PORT", "9092"),

		OTELServiceName:      getEnv("OTEL_SERVICE_NAME", "otel_api_demo"),
		OTELExporterEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
		OTELExporterProtocol: getEnv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc"),

		ReadTimeout:  parseDuration(getEnv("READ_TIMEOUT", "5s")),
		WriteTimeout: parseDuration(getEnv("WRITE_TIMEOUT", "10s")),
		IdleTimeout:  parseDuration(getEnv("IDLE_TIMEOUT", "0")),
	}, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func (c *Config) GetKafkaAddr() string {
	return fmt.Sprintf("%s:%s", c.KafkaHost, c.KafkaPort)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}
