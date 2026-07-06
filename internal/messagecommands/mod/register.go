package mod

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register mod text commands
func Register() {
	messagecommands.RegisterCommand("ban", banCommand)
	messagecommands.RegisterCommand("kick", kickCommand)
	messagecommands.RegisterCommand("mute", muteCommand)
	messagecommands.RegisterCommand("tempban", tempbanCommand)
	messagecommands.RegisterCommand("softban", softbanCommand)
	messagecommands.RegisterCommand("warn", warnCommand)
	messagecommands.RegisterCommand("removewarn", removewarnCommand)
	messagecommands.RegisterCommand("clear", clearCommand)
	messagecommands.RegisterCommand("lockdown", lockdownCommand)
	messagecommands.RegisterCommand("nuke", nukeCommand)
	messagecommands.RegisterCommand("warns", warningsCommand)
	messagecommands.RegisterCommand("assign-role", assignRoleCommand)
	messagecommands.RegisterCommand("removerole", removeRoleCommand)
}
