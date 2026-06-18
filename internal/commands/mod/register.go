// Package mod provides moderation commands organized as subcommands under /mod
// Each command is in its own file for better organization
package mod

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// RegisterModCommands registers all moderation commands as /mod subcommands
func RegisterModCommands(client *discord.ExtendedClient) {
	// Create individual subcommands (each can be in its own file)
	commands := []*discord.Command{
		createWarnCommand(),
		createRemoveWarnCommand(),
		createWarningsCommand(),
		createKickCommand(),
		createBanCommand(),
		createMuteCommand(),
		createClearCommand(),
		createSoftbanCommand(),
		createAssignRoleCommand(),
		createRemoveRoleCommand(),
		createNukeCommand(),
		createLockdownCommand(),
		createTempBanCommand(),
	}

	// Build the /mod command group with all subcommands
	modGroup := client.CommandHandler.BuildCommandGroup(
		"mod",
		"Comandos de moderación",
		commands...,
	)

	// Register the command group
	client.CommandHandler.AddGlobalCommand(modGroup)
}
