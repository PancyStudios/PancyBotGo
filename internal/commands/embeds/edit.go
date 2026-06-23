package embeds

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createEmbedEditCommand() *discord.Command {
	return &discord.Command{
		Name:        "edit",
		Description: "📝 | Carga un embed existente para editarlo",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "mensaje_id",
				Description: "📝 | ID del mensaje que contiene el embed (debe estar en el mismo canal)",
				Required:    true,
			},
		},
		Run: func(ctx *discord.CommandContext) error {
			messageID := ctx.GetStringOption("mensaje_id")

			msg, err := ctx.Session.ChannelMessage(ctx.Interaction.ChannelID, messageID)
			if err != nil {
				return ctx.ReplyEphemeral("❌ No pude encontrar el mensaje. Asegúrate de estar en el canal donde se envió y de que el ID sea correcto.")
			}

			if msg.Author.ID != ctx.Session.State.User.ID {
				return ctx.ReplyEphemeral("❌ Solo puedo editar mis propios mensajes.")
			}

			if len(msg.Embeds) == 0 {
				return ctx.ReplyEphemeral("❌ Este mensaje no contiene ningún embed.")
			}

			user := ctx.User()

			// Load the first embed into the builder state
			saveBuilderState(user.ID, msg.Embeds[0])

			return ctx.ReplyEphemeral("✅ Embed cargado con éxito. Ahora ejecuta `/embed create` para modificar sus propiedades y usar `/embed send` o editar manualmente con el ID.")
		},
	}
}
