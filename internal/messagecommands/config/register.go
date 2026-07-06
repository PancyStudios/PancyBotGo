package config

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register config text commands
func Register() {
	messagecommands.RegisterCommand("poj", "Comando poj", "pan!poj", "Config", pojCommand)
	messagecommands.RegisterCommand("autorole", "Comando autorole", "pan!autorole", "Config", autoroleCommand)
	messagecommands.RegisterCommand("welcome", "Comando welcome", "pan!welcome", "Config", welcomeCommand)
	messagecommands.RegisterCommand("leave", "Comando leave", "pan!leave", "Config", leaveCommand)
	messagecommands.RegisterCommand("channels", "Comando channels", "pan!channels", "Config", channelsCommand)
	messagecommands.RegisterCommand("logs", "Comando logs", "pan!logs", "Config", logsCommand)
}
