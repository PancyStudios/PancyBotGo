package utils

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register util text commands
func Register() {
	messagecommands.RegisterCommand("ping", "Comando ping", "pan!ping", "General", pingCommand)
	messagecommands.RegisterCommand("botinfo", "Comando botinfo", "pan!botinfo", "General", botinfoCommand)
	messagecommands.RegisterCommand("suggest", "Comando suggest", "pan!suggest", "General", suggestCommand)
	messagecommands.RegisterCommand("confess", "Comando confess", "pan!confess", "General", confessCommand)
	messagecommands.RegisterCommand("screenshot", "Comando screenshot", "pan!screenshot", "General", screenshotCommand)
	messagecommands.RegisterCommand("status", "Comando status", "pan!status", "General", statusCommand)
	messagecommands.RegisterCommand("invite", "Comando invite", "pan!invite", "General", inviteCommand)
}
