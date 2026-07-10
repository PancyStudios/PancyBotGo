package utils

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// RegisterModCommands registers all moderation commands as /utils subcommands
func RegisterUtilsCommands(client *discord.ExtendedClient) {
	// Create individual subcommands (each can be in its own file)
	pingCmd := createPingCommand()
	statusCmd := createStatusCommand()
	botinfoCmd := createBotinfoCommand()
	inviteCmd := createInviteCommand()
	screenshotCmd := createScreenshotCommand()
	suggestCmd := createSuggestCommand()
	confessCmd := createConfessCommand()
	avatarCmd := createAvatarCommand()

	// Build the /utils command group with all subcommands
	modGroup := client.CommandHandler.BuildCommandGroup(
		"utils",
		"Comandos de utilidad",
		pingCmd,
		statusCmd,
		botinfoCmd,
		inviteCmd,
		screenshotCmd,
		suggestCmd,
		confessCmd,
		avatarCmd,
	)

	// Register the command group
	client.CommandHandler.AddGlobalCommand(modGroup)
}
