package mod

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register mod text commands
func Register() {
	messagecommands.RegisterCommand("ban", "Comando ban", "pan!ban", "Mod", banCommand)
	messagecommands.RegisterCommand("kick", "Comando kick", "pan!kick", "Mod", kickCommand)
	messagecommands.RegisterCommand("mute", "Comando mute", "pan!mute", "Mod", muteCommand)
	messagecommands.RegisterCommand("tempban", "Comando tempban", "pan!tempban", "Mod", tempbanCommand)
	messagecommands.RegisterCommand("softban", "Comando softban", "pan!softban", "Mod", softbanCommand)
	messagecommands.RegisterCommand("warn", "Comando warn", "pan!warn", "Mod", warnCommand)
	messagecommands.RegisterCommand("removewarn", "Comando removewarn", "pan!removewarn", "Mod", removewarnCommand)
	messagecommands.RegisterCommand("clear", "Comando clear", "pan!clear", "Mod", clearCommand)
	messagecommands.RegisterCommand("lockdown", "Comando lockdown", "pan!lockdown", "Mod", lockdownCommand)
	messagecommands.RegisterCommand("nuke", "Comando nuke", "pan!nuke", "Mod", nukeCommand)
	messagecommands.RegisterCommand("warns", "Comando warns", "pan!warns", "Mod", warningsCommand)
	messagecommands.RegisterCommand("assign-role", "Comando assign-role", "pan!assign-role", "Mod", assignRoleCommand)
	messagecommands.RegisterCommand("removerole", "Comando removerole", "pan!removerole", "Mod", removeRoleCommand)
}
