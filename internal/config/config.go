package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Telegram     TelegramConfig     `yaml:"telegram"`
	AI           AIConfig           `yaml:"ai"`
	Database     DatabaseConfig     `yaml:"database"`
	Scheduler    SchedulerConfig    `yaml:"scheduler"`
	RateLimiting RateLimitingConfig `yaml:"rate_limiting"`
}

// TelegramConfig holds Telegram-related settings
type TelegramConfig struct {
	AppID         int    `yaml:"app_id"`
	AppHash       string `yaml:"app_hash"`
	BotToken      string `yaml:"bot_token"`
	AdminID       int64  `yaml:"admin_id"`
	DeviceModel   string `yaml:"device_model"`
	SystemVersion string `yaml:"system_version"`
	AppVersion    string `yaml:"app_version"`
}

// AIConfig holds AI provider settings
type AIConfig struct {
	ZhipuAPIKey string  `yaml:"zhipu_api_key"`
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"`
}

// SchedulerConfig holds scheduler settings
type SchedulerConfig struct {
	ReportTime              string `yaml:"report_time"`
	RawMessagesRetentionDays int    `yaml:"raw_messages_retention_days"`
}

// RateLimitingConfig holds rate limiting settings
type RateLimitingConfig struct {
	MinDelay                int `yaml:"min_delay"`
	MaxDelay                int `yaml:"max_delay"`
	TranscriptionsPerMinute int `yaml:"transcriptions_per_minute"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	cfg := &Config{
		// Set defaults
		Telegram: TelegramConfig{
			DeviceModel:   "Desktop",
			SystemVersion: "Windows 10",
			AppVersion:    "4.9.0",
		},
		AI: AIConfig{
			Model:       "glm-4",
			Temperature: 0.3,
			MaxTokens:   2000,
		},
		Database: DatabaseConfig{
			Host:    "localhost",
			Port:    5432,
			Name:    "telemonitor",
			User:    "telemonitor_user",
			SSLMode: "disable",
		},
		Scheduler: SchedulerConfig{
			ReportTime:              "08:00",
			RawMessagesRetentionDays: 7,
		},
		RateLimiting: RateLimitingConfig{
			MinDelay:                1000,
			MaxDelay:                5000,
			TranscriptionsPerMinute: 10,
		},
	}

	// Try to load from config.yaml
	if data, err := os.ReadFile("config.yaml"); err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
		}
	}

	// Override with environment variables
	if v := os.Getenv("TG_APP_ID"); v != "" {
		if appID, err := strconv.Atoi(v); err == nil {
			cfg.Telegram.AppID = appID
		}
	}
	if v := os.Getenv("TG_APP_HASH"); v != "" {
		cfg.Telegram.AppHash = v
	}
	if v := os.Getenv("TG_BOT_TOKEN"); v != "" {
		cfg.Telegram.BotToken = v
	}
	if v := os.Getenv("TG_ADMIN_ID"); v != "" {
		if adminID, err := strconv.ParseInt(v, 10, 64); err == nil {
			cfg.Telegram.AdminID = adminID
		}
	}

	if v := os.Getenv("ZHIPU_API_KEY"); v != "" {
		cfg.AI.ZhipuAPIKey = v
	}

	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Database.Port = port
		}
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.Name = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DB_SSL_MODE"); v != "" {
		cfg.Database.SSLMode = v
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Telegram validation
	if c.Telegram.AppID == 0 {
		return fmt.Errorf("telegram.app_id is required")
	}
	if c.Telegram.AppHash == "" || c.Telegram.AppHash == "your_app_hash_here" {
		return fmt.Errorf("telegram.app_hash is required")
	}
	if c.Telegram.BotToken == "" || c.Telegram.BotToken == "your_bot_token_here" {
		return fmt.Errorf("telegram.bot_token is required")
	}
	if c.Telegram.AdminID == 0 || c.Telegram.AdminID == 999999 {
		return fmt.Errorf("telegram.admin_id must be set to your actual Telegram user ID")
	}

	// AI validation
	if c.AI.ZhipuAPIKey == "" || c.AI.ZhipuAPIKey == "your_zhipu_api_key" {
		return fmt.Errorf("ai.zhipu_api_key is required")
	}

	// Database validation
	if c.Database.Password == "" {
		return fmt.Errorf("database.password is required")
	}

	return nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}
