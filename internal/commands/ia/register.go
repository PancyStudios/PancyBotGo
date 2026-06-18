package ia

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// Register registers all AI commands as /ia subcommands
func Register(client *discord.ExtendedClient) {
	createImageCmd := createCreateImageCommand()
	getImageCmd := createGetImageCommand()

	iaGroup := client.CommandHandler.BuildCommandGroup(
		"ia",
		"Comandos de Inteligencia Artificial",
		createImageCmd,
		getImageCmd,
	)

	client.CommandHandler.AddGlobalCommand(iaGroup)
}
