package embeds

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

func createEmbedDeleteCommand() *discord.Command {
	return &discord.Command{
		Name:        "delete",
		Description: "📝 | Elimina tu progreso actual de creación de embed",
		Run: func(ctx *discord.CommandContext) error {
			user := ctx.User()
			clearBuilderState(user.ID)
			return ctx.ReplyEphemeral("🗑️ El embed en el que estabas trabajando ha sido descartado.")
		},
	}
}
