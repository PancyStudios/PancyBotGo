package embeds

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

func createEmbedSendCommand() *discord.Command {
	return &discord.Command{
		Name:        "send",
		Description: "📝 | Envía el embed que estás creando actualmente",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionChannel,
				Name:         "canal",
				Description:  "📝 | Canal donde enviar el embed (opcional, por defecto el actual)",
				Required:     false,
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText, discordgo.ChannelTypeGuildNews},
			},
		},
		Run: func(ctx *discord.CommandContext) error {
			user := ctx.User()

			// Get embed state
			builderMutex.RLock()
			embedState, exists := builderStateMap[user.ID]
			builderMutex.RUnlock()

			if !exists {
				return ctx.ReplyEphemeral("❌ No tienes ningún embed en construcción. Usa `/embed create` primero.")
			}

			// Determine target channel
			targetChannelID := ctx.Interaction.ChannelID
			chOpt := ctx.GetChannelOption("canal")
			if chOpt != nil {
				targetChannelID = chOpt.ID
			}

			// Send embed
			_, err := ctx.Session.ChannelMessageSendEmbed(targetChannelID, embedState)
			if err != nil {
				logger.Error(fmt.Sprintf("Error sending embed: %v", err), "Embeds")
				return ctx.ReplyEphemeral("❌ No pude enviar el embed. Verifica que tengo permisos de escribir y enviar embeds en ese canal.")
			}

			return ctx.ReplyEphemeral(fmt.Sprintf("✅ Embed enviado exitosamente a <#%s>", targetChannelID))
		},
	}
}
