package config

import "sync"

// GlobalBotConfig holds the global bot configuration fetched from MongoDB or MQTT
type GlobalBotConfig struct {
	MaintenanceMode  bool
	DisabledCommands []string
	Mu               sync.RWMutex
}

var (
	GlobalConfig *GlobalBotConfig
	onceConfig   sync.Once
)

// GetBotConfig returns the singleton instance of GlobalBotConfig
func GetBotConfig() *GlobalBotConfig {
	onceConfig.Do(func() {
		GlobalConfig = &GlobalBotConfig{
			MaintenanceMode:  false,
			DisabledCommands: make([]string, 0),
		}
	})
	return GlobalConfig
}

// Update updates the configuration in a thread-safe manner
func (c *GlobalBotConfig) Update(maintenance bool, disabled []string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.MaintenanceMode = maintenance
	c.DisabledCommands = disabled
}

// IsMaintenance checks if maintenance mode is active
func (c *GlobalBotConfig) IsMaintenance() bool {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.MaintenanceMode
}

// IsCommandDisabled checks if a specific command is disabled
func (c *GlobalBotConfig) IsCommandDisabled(cmdName string) bool {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	for _, cmd := range c.DisabledCommands {
		if cmd == cmdName {
			return true
		}
	}
	return false
}
