package security

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// RegisterSecurityCommands registers all security commands as /security subcommands
func RegisterSecurityCommands(client *discord.ExtendedClient) {
	commands := []*discord.Command{
		createAntibotsCommand(),
	}

	securityGroup := client.CommandHandler.BuildCommandGroup(
		"security",
		"Comandos de seguridad",
		commands...,
	)

	client.CommandHandler.AddGlobalCommand(securityGroup)
}
