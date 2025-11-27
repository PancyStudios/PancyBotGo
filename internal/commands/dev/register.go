package dev

import "github.com/PancyStudios/PancyBotGo/pkg/discord"

// Register registers all dev commands as /dev subcommands (only in dev guild)
func Register(client *discord.ExtendedClient) {
	// Create individual subcommands
	codegenCmd := CreateCodeGenCommand()
	codelistCmd := CreateCodeListCommand()
	codedelCmd := CreateCodeDelCommand()
	evalCmd := CreateEvalCommand()

	// Build the /dev command group with all subcommands
	devGroup := client.CommandHandler.BuildCommandGroup(
		"dev",
		"Comandos de desarrollo",
		codegenCmd,
		codelistCmd,
		codedelCmd,
		evalCmd,
	)

	// Register the command group as dev-only command
	client.CommandHandler.AddDevCommand(devGroup)
}
