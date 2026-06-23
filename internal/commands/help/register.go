package help

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// Register registers all help commands
func Register(client *discord.ExtendedClient) {
	cmdsCmd := createCmdsCommand()

	helpGroup := client.CommandHandler.BuildCommandGroup(
		"help",
		"📚 | Centro de ayuda e información del bot",
		cmdsCmd,
	)

	client.CommandHandler.AddGlobalCommand(helpGroup)
}
