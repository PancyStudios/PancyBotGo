package config

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register config text commands
func Register() {
	messagecommands.RegisterCommand("poj", pojCommand)
	messagecommands.RegisterCommand("autorole", autoroleCommand)
	messagecommands.RegisterCommand("welcome", welcomeCommand)
	messagecommands.RegisterCommand("leave", leaveCommand)
	messagecommands.RegisterCommand("channels", channelsCommand)
	messagecommands.RegisterCommand("logs", logsCommand)
}
