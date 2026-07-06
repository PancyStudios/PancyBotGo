package config

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register config text commands
func Register() {
	messagecommands.RegisterCommand("poj", "Comando poj", "pan!poj", "General", pojCommand)
	messagecommands.RegisterCommand("autorole", "Comando autorole", "pan!autorole", "General", autoroleCommand)
	messagecommands.RegisterCommand("welcome", "Comando welcome", "pan!welcome", "General", welcomeCommand)
	messagecommands.RegisterCommand("leave", "Comando leave", "pan!leave", "General", leaveCommand)
	messagecommands.RegisterCommand("channels", "Comando channels", "pan!channels", "General", channelsCommand)
	messagecommands.RegisterCommand("logs", "Comando logs", "pan!logs", "General", logsCommand)
}
