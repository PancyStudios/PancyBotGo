package utils

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register util text commands
func Register() {
	messagecommands.RegisterCommand("ping", "Comando ping", "pan!ping", "Utils", pingCommand)
	messagecommands.RegisterCommand("botinfo", "Comando botinfo", "pan!botinfo", "Utils", botinfoCommand)
	messagecommands.RegisterCommand("suggest", "Comando suggest", "pan!suggest", "Utils", suggestCommand)
	messagecommands.RegisterCommand("confess", "Comando confess", "pan!confess", "Utils", confessCommand)
	messagecommands.RegisterCommand("screenshot", "Comando screenshot", "pan!screenshot", "Utils", screenshotCommand)
	messagecommands.RegisterCommand("status", "Comando status", "pan!status", "Utils", statusCommand)
	messagecommands.RegisterCommand("invite", "Comando invite", "pan!invite", "Utils", inviteCommand)
}
