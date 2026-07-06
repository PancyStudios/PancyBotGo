package mod

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register mod text commands
func Register() {
	messagecommands.RegisterCommand("ban", "Comando ban", "pan!ban", "General", banCommand)
	messagecommands.RegisterCommand("kick", "Comando kick", "pan!kick", "General", kickCommand)
	messagecommands.RegisterCommand("mute", "Comando mute", "pan!mute", "General", muteCommand)
	messagecommands.RegisterCommand("tempban", "Comando tempban", "pan!tempban", "General", tempbanCommand)
	messagecommands.RegisterCommand("softban", "Comando softban", "pan!softban", "General", softbanCommand)
	messagecommands.RegisterCommand("warn", "Comando warn", "pan!warn", "General", warnCommand)
	messagecommands.RegisterCommand("removewarn", "Comando removewarn", "pan!removewarn", "General", removewarnCommand)
	messagecommands.RegisterCommand("clear", "Comando clear", "pan!clear", "General", clearCommand)
	messagecommands.RegisterCommand("lockdown", "Comando lockdown", "pan!lockdown", "General", lockdownCommand)
	messagecommands.RegisterCommand("nuke", "Comando nuke", "pan!nuke", "General", nukeCommand)
	messagecommands.RegisterCommand("warns", "Comando warns", "pan!warns", "General", warningsCommand)
	messagecommands.RegisterCommand("assign-role", "Comando assign-role", "pan!assign-role", "General", assignRoleCommand)
	messagecommands.RegisterCommand("removerole", "Comando removerole", "pan!removerole", "General", removeRoleCommand)
}
