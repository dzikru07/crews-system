package config

import (
	"os"
	"time"
)

func Load() Config {
	return Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "crews-api-platform"),
			Environment: getEnv("APP_ENV", "development"),
			Version:     getEnv("APP_VERSION", "0.1.0"),
		},

		HTTP: HTTPConfig{
			Host: getEnv("HTTP_HOST", "0.0.0.0"),
			Port: getEnv("HTTP_PORT", "18080"),
		},

		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},

		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", ""),
		},

		JWT: JWTConfig{
			Secret: getEnv(
				"JWT_SECRET",
				"",
			),

			Issuer: getEnv(
				"JWT_ISSUER",
				"crews-api-platform",
			),

			AccessTTL: getDurationEnv(
				"JWT_ACCESS_TTL",
				time.Hour,
			),
		},
	}
}

func getEnv(
	key string,
	fallback string,
) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func getDurationEnv(
	key string,
	fallback time.Duration,
) time.Duration {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)

	if err != nil {
		return fallback
	}

	return duration
}
