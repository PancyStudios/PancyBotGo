package utils

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register util text commands
func Register() {
	messagecommands.RegisterCommand("ping", pingCommand)
	messagecommands.RegisterCommand("botinfo", botinfoCommand)
	messagecommands.RegisterCommand("suggest", suggestCommand)
	messagecommands.RegisterCommand("confess", confessCommand)
	messagecommands.RegisterCommand("screenshot", screenshotCommand)
	messagecommands.RegisterCommand("status", statusCommand)
	messagecommands.RegisterCommand("invite", inviteCommand)
}
