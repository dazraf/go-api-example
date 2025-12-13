package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Environment string   `yaml:"environment"`
	Server      Server   `yaml:"server"`
	Database    Database `yaml:"database"`
	Logging     Logging  `yaml:"logging"`
}

// Server holds server configuration
type Server struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

// Database holds database configuration
type Database struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Logging holds logging configuration
type Logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	// Set defaults
	cfg := &Config{
		Environment: "development",
		Server: Server{
			Address: ":8080",
			Port:    8080,
		},
		Database: Database{
			Type: "memory",
		},
		Logging: Logging{
			Level:  "info",
			Format: "json",
		},
	}

	// Load from config file
	configFile := getConfigFile()
	if configFile != "" {
		if err := loadFromFile(cfg, configFile); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	// Override with environment variables
	loadFromEnv(cfg)

	return cfg, nil
}

// getConfigFile returns the config file path based on environment
func getConfigFile() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	configFile := fmt.Sprintf("configs/config.%s.yaml", env)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Fall back to default config
		configFile = "configs/config.yaml"
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			return ""
		}
	}

	return configFile
}

// loadFromFile loads configuration from YAML file
func loadFromFile(cfg *Config, filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) {
	if env := os.Getenv("GO_ENV"); env != "" {
		cfg.Environment = env
	}
	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		cfg.Server.Address = addr
	}
	if dbType := os.Getenv("DB_TYPE"); dbType != "" {
		cfg.Database.Type = dbType
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.Database.Host = dbHost
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.Logging.Level = logLevel
	}
}
