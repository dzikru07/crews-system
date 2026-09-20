package config

import "os"

func Load() Config {
	return Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "crews-api-platform"),
			Environment: getEnv("APP_ENV", "development"),
			Version:     getEnv("APP_VERSION", "0.1.0"),
		},
		HTTP: HTTPConfig{
			Host: getEnv("HTTP_HOST", "0.0.0.0"),
			Port: getEnv("HTTP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", ""),
		},
		Security: SecurityConfig{
			JWTSecret: getEnv("JWT_SECRET", ""),
		},
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
