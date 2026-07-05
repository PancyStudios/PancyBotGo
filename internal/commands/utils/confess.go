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

func createConfessCommand() *discord.Command {
	return &discord.Command{
		Name:        "confess",
		Description: "🧰 | Envía una confesión anónima al servidor",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "confesion",
				Description: "🧰 | Tu confesión secreta",
				Required:    true,
			},
		},
		Run: func(ctx *discord.CommandContext) error {
			confesion := ctx.GetStringOption("confesion")

			guildData, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Interaction.GuildID})
			if err != nil || guildData == nil || guildData.Configuration.SubData.ConfessionChannel == "" {
				return ctx.ReplyEphemeral("❌ El sistema de confesiones no está configurado en este servidor.")
			}

			confessChannel := guildData.Configuration.SubData.ConfessionChannel

			embed := &discordgo.MessageEmbed{
				Title:       "🕵️ Nueva Confesión Anónima",
				Description: confesion,
				Color:       0x9B59B6, // Purple
				Timestamp:   time.Now().Format(time.RFC3339),
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(confessChannel, embed)
			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando confesion: %v", err), "Confess")
				return ctx.ReplyEphemeral("❌ No pude enviar la confesión. Verifica que tengo permisos en el canal configurado.")
			}

			return ctx.ReplyEphemeral("✅ Tu confesión ha sido enviada de forma completamente anónima 🤫.")
		},
	}
}
