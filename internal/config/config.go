package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTP       *HTTPConfig
	WorkerPool *WorkerPoolConfig
}

type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
type WorkerPoolConfig struct {
	Workers   int
	QueueSize int
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}

	return cfg
}

func Load() (*Config, error) {
	cfg := Config{
		HTTP:       &HTTPConfig{},
		WorkerPool: &WorkerPoolConfig{},
	}

	var err error

	// HTTP Config
	if cfg.HTTP.Port, err = getIntEnv("API_PORT", 8080); err != nil {
		return nil, fmt.Errorf("http port: %w", err)
	}
	if cfg.HTTP.ReadTimeout, err = getDurationEnv("HTTP_READ_TIMEOUT", 15*time.Second); err != nil {
		return nil, fmt.Errorf("read timeout: %w", err)
	}
	if cfg.HTTP.WriteTimeout, err = getDurationEnv("HTTP_WRITE_TIMEOUT", 15*time.Second); err != nil {
		return nil, fmt.Errorf("write timeout: %w", err)
	}

	// WorkerPool Config
	if cfg.WorkerPool.Workers, err = getIntEnv("WORKERS", 4); err != nil {
		return nil, fmt.Errorf("workers: %w", err)
	}
	if cfg.WorkerPool.QueueSize, err = getIntEnv("QUEUE_SIZE", 64); err != nil {
		return nil, fmt.Errorf("queue_size: %w", err)
	}

	return &cfg, nil
}

func getStringEnv(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("env %s is required", name)
	}
	return val, nil
}

func getIntEnv(name string, defaultValue int) (int, error) {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(val)
}

func getDurationEnv(name string, defaultValue time.Duration) (time.Duration, error) {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue, nil
	}
	return time.ParseDuration(val)
}

func (c *Config) validate() error {
	if c.HTTP.Port <= 0 || c.HTTP.Port > 65535 {
		return fmt.Errorf("http port out of range")
	}
	if c.HTTP.ReadTimeout < 0 {
		return fmt.Errorf("read timeout must be non-negative")
	}
	if c.HTTP.WriteTimeout < 0 {
		return fmt.Errorf("write timeout must be non-negative")
	}

	if c.WorkerPool.Workers < 1 {
		return fmt.Errorf("the number of workers must be greater than 0")
	}
	if c.WorkerPool.QueueSize < 1 {
		return fmt.Errorf("the queue size must be greater than 0")
	}
	return nil
}
