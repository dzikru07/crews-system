package config

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func (c Config) Validate() error {
	if err := c.validateApp(); err != nil {
		return err
	}

	if err := c.validateHTTP(); err != nil {
		return err
	}

	if err := c.validateDatabase(); err != nil {
		return err
	}

	if err := c.validateRedis(); err != nil {
		return err
	}

	if err := c.validateSecurity(); err != nil {
		return err
	}

	return nil
}

func (c Config) validateApp() error {
	if strings.TrimSpace(c.App.Name) == "" {
		return errors.New("APP_NAME is required")
	}

	switch c.App.Environment {
	case "development", "staging", "production":
		return nil
	default:
		return fmt.Errorf(
			"APP_ENV must be one of: development, staging, production; got %q",
			c.App.Environment,
		)
	}
}

func (c Config) validateHTTP() error {
	if strings.TrimSpace(c.HTTP.Host) == "" {
		return errors.New("HTTP_HOST is required")
	}

	port, err := strconv.Atoi(c.HTTP.Port)
	if err != nil {
		return fmt.Errorf("HTTP_PORT must be a number: %w", err)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("HTTP_PORT must be between 1 and 65535")
	}

	return nil
}

func (c Config) validateDatabase() error {
	if c.Database.URL == "" {
		if c.App.Environment == "production" {
			return errors.New("DATABASE_URL is required in production")
		}

		return nil
	}

	parsed, err := url.Parse(c.Database.URL)
	if err != nil {
		return fmt.Errorf("invalid DATABASE_URL: %w", err)
	}

	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return fmt.Errorf(
			"DATABASE_URL must use postgres or postgresql scheme",
		)
	}

	return nil
}

func (c Config) validateRedis() error {
	if c.Redis.URL == "" {
		if c.App.Environment == "production" {
			return errors.New("REDIS_URL is required in production")
		}

		return nil
	}

	parsed, err := url.Parse(c.Redis.URL)
	if err != nil {
		return fmt.Errorf("invalid REDIS_URL: %w", err)
	}

	if parsed.Scheme != "redis" && parsed.Scheme != "rediss" {
		return fmt.Errorf(
			"REDIS_URL must use redis or rediss scheme",
		)
	}

	return nil
}

func (c Config) validateSecurity() error {
	if c.App.Environment == "production" {
		if strings.TrimSpace(c.Security.JWTSecret) == "" {
			return errors.New("JWT_SECRET is required in production")
		}

		if len(c.Security.JWTSecret) < 32 {
			return errors.New(
				"JWT_SECRET must be at least 32 characters in production",
			)
		}
	}

	return nil
}
