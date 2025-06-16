package config

import (
	"os"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	OpenAIAPIKey string
	AWSConfig    AWSConfig
	ServerConfig ServerConfig
}

// AWSConfig holds AWS configuration
type AWSConfig struct {
	Region          string
	S3Bucket        string
	SNSTopicARN     string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// New creates a new Config
func New() *Config {
	return &Config{
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
		AWSConfig: AWSConfig{
			Region:      getEnvOrDefault("AWS_REGION", "ap-northeast-2"),
			S3Bucket:    getEnvOrDefault("AWS_S3_BUCKET", "snp500-stocks"),
			SNSTopicARN: getEnvOrDefault("AWS_SNS_TOPIC_ARN", ""),
		},
		ServerConfig: ServerConfig{
			Port:         getEnvOrDefault("SERVER_PORT", "8080"),
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 