package utils

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func createSuggestCommand() *discord.Command {
	return &discord.Command{
		Name:        "suggest",
		Description: "🧰 | Envía una sugerencia para el servidor",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "sugerencia",
				Description: "🧰 | Tu sugerencia",
				Required:    true,
			},
		},
		Run: func(ctx *discord.CommandContext) error {
			sugerencia := ctx.GetStringOption("sugerencia")

			guildData, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Interaction.GuildID})
			if err != nil || guildData == nil || guildData.Configuration.SubData.SuggestChannel == "" {
				return ctx.ReplyEphemeral("❌ El sistema de sugerencias no está configurado en este servidor.")
			}

			suggestChannel := guildData.Configuration.SubData.SuggestChannel

			user := ctx.User()

			embed := &discordgo.MessageEmbed{
				Title:       "💡 Nueva Sugerencia",
				Description: sugerencia,
				Color:       0xF1C40F, // Yellow
				Author: &discordgo.MessageEmbedAuthor{
					Name:    user.String(),
					IconURL: user.AvatarURL(""),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			msg, err := ctx.Session.ChannelMessageSendEmbed(suggestChannel, embed)
			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando sugerencia: %v", err), "Suggest")
				return ctx.ReplyEphemeral("❌ No pude enviar la sugerencia. Verifica que tengo permisos en el canal configurado.")
			}

			// Add reactions
			ctx.Session.MessageReactionAdd(suggestChannel, msg.ID, "✅")
			ctx.Session.MessageReactionAdd(suggestChannel, msg.ID, "❌")

			return ctx.ReplyEphemeral(fmt.Sprintf("✅ Tu sugerencia ha sido enviada a <#%s>", suggestChannel))
		},
	}
}
