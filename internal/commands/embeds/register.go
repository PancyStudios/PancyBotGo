package embeds

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// RegisterEmbedsCommands registers the /embed command group
func RegisterEmbedsCommands(client *discord.ExtendedClient) {
	createCmd := createEmbedCreateCommand()
	sendCmd := createEmbedSendCommand()
	deleteCmd := createEmbedDeleteCommand()
	editCmd := createEmbedEditCommand()

	embedGroup := client.CommandHandler.BuildCommandGroup(
		"embed",
		"Constructor interactivo de Embeds",
		createCmd,
		sendCmd,
		deleteCmd,
		editCmd,
	)

	// Solo para usuarios con permisos de ManageMessages (opcional, configurarlo desde discord)
	// En discordgo se configuran los default permissions

	defaultPerms := int64(8192) // Manage Messages
	embedGroup.DefaultMemberPermissions = &defaultPerms

	client.CommandHandler.AddGlobalCommand(embedGroup)
}
