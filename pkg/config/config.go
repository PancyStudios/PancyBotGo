// Package config provides configuration management for the bot.
// It loads environment variables and makes them available throughout the application.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration values for the bot
type Config struct {
	// Discord
	BotToken    string
	DevGuildID  string
	DeveloperID string

	// MongoDB
	MongoDBURL string
	DBName     string

	// Redis cache (optional)
	RedisURL    string
	RedisPrefix string
	CacheTTL    time.Duration

	// MQTT
	MQTTHost     string
	MQTTPort     string
	MQTTUser     string
	MQTTPassword string

	// Web Server
	Port          string
	AdminUser     string
	AdminPassword string

	// Environment
	Environment string

	// Webhooks
	ErrorWebhook      string
	LogsWebhook       string
	LogsWebServerHook string
	GuildsWebhook     string

	// Lavalink
	LinkServer   string
	LinkPassword string

	//Craiyon
	CraiyonToken string
}

var (
	Version   = "Dev-Local"
	BuildTime = "Hoy"
)

// cfg holds the global configuration instance
var (
	cfg     *Config
	cfgErr  error
	cfgOnce sync.Once
)

// resetForTesting resets the configuration for testing purposes.
// This function should only be called from test code.
func resetForTesting() {
	cfg = nil
	cfgErr = nil
	cfgOnce = sync.Once{}
}

// loadConfig performs the actual configuration loading
func loadConfig() {
	// Load .env file if it exists (ignoring error if it doesn't)
	_ = godotenv.Load()

	cfg = &Config{
		// Discord
		BotToken:    getEnvAny([]string{"BOT_TOKEN", "botToken"}, ""),
		DevGuildID:  getEnvAny([]string{"DEV_GUILD_ID", "devGuildId"}, ""),
		DeveloperID: getEnv("DEVELOPER_ID", ""),

		// MongoDB
		MongoDBURL: getEnvAny([]string{"MONGODB_URL", "mongodbUrl"}, "mongodb://localhost:27017"),
		DBName:     getEnvAny([]string{"DB_NAME", "dbName"}, "PancyBot"),

		// Redis
		RedisURL:    strings.TrimSpace(getEnv("REDIS_URL", "")),
		RedisPrefix: strings.TrimSpace(getEnv("REDIS_PREFIX", "")),
		CacheTTL:    5 * time.Minute,

		// MQTT
		MQTTHost:     getEnv("MQTT_Host", "localhost"),
		MQTTPort:     getEnv("MQTT_Port", "1883"),
		MQTTUser:     getEnv("MQTT_User", ""),
		MQTTPassword: getEnv("MQTT_Password", ""),

		// Web Server
		Port:          getEnv("PORT", "3000"),
		AdminUser:     getEnv("ADMIN_USER", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),

		// Environment
		Environment: strings.ToLower(getEnvAny([]string{"BOT_ENV", "ENVIRONMENT", "enviroment"}, "dev")),

		// Webhooks
		ErrorWebhook:      getEnv("errorWebhook", ""),
		LogsWebhook:       getEnv("logsWebhook", ""),
		LogsWebServerHook: getEnv("logsWebServerWebhook", ""),
		GuildsWebhook:     getEnv("guildsWebhook", ""),

		// Lavalink
		LinkServer:   getEnv("linkserver", "localhost"),
		LinkPassword: getEnv("linkpassword", ""),

		// Craiyon
		CraiyonToken: getEnv("craiyonToken", ""),
	}

	if cfg.RedisPrefix == "" {
		cfg.RedisPrefix = "pancybot:" + cfg.Environment
	}
	if rawTTL := strings.TrimSpace(getEnv("CACHE_TTL", "5m")); rawTTL != "" {
		ttl, err := time.ParseDuration(rawTTL)
		if err != nil || ttl <= 0 {
			cfgErr = fmt.Errorf("CACHE_TTL must be a positive Go duration such as 5m or 1h")
			return
		}
		cfg.CacheTTL = ttl
	}

	if strings.TrimSpace(cfg.BotToken) == "" {
		cfgErr = fmt.Errorf("BOT_TOKEN (or legacy botToken) is required")
		return
	}
	if strings.TrimSpace(cfg.DeveloperID) == "" {
		cfgErr = fmt.Errorf("DEVELOPER_ID is required")
		return
	}
	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port < 1 || port > 65535 {
		cfgErr = fmt.Errorf("PORT must be a number between 1 and 65535")
	}
	if cfg.IsProd() && strings.TrimSpace(cfg.AdminPassword) == "" {
		cfgErr = fmt.Errorf("ADMIN_PASSWORD is required in production")
		return
	}
	if cfg.IsProd() && (strings.TrimSpace(cfg.MQTTUser) == "" || strings.TrimSpace(cfg.MQTTPassword) == "") {
		cfgErr = fmt.Errorf("MQTT_User and MQTT_Password are required in production")
	}
}

// Load initializes the configuration from environment variables
func Load() (*Config, error) {
	cfgOnce.Do(loadConfig)
	return cfg, cfgErr
}

// Get returns the current configuration
func Get() *Config {
	// Use sync.Once to ensure thread-safe initialization if Load wasn't called
	cfgOnce.Do(loadConfig)
	return cfg
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAny(keys []string, defaultValue string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return defaultValue
}

// IsProd returns true if the environment is production
func (c *Config) IsProd() bool {
	return c.Environment == "prod"
}
