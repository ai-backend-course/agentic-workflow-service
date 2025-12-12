package config

import "os"

type Config struct {
	ServiceName string
	Env         string
	HTTPAddr    string
	OtelEnabled bool
}

func Load() Config {
	return Config{
		ServiceName: getEnv("SERVICE_NAME", "agentic-workflow-service"),
		Env:         getEnv("APP_ENV", "dev"),
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		OtelEnabled: getEnv("OTEL_ENABLED", "true") == "true",
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
