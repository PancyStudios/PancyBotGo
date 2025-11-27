package discord

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// BlacklistMiddleware verifica si el usuario o guild está en la blacklist
func (c *ExtendedClient) BlacklistMiddleware(ctx *CommandContext) error {
	userID := ctx.User().ID
	guildID := ctx.Interaction.GuildID

	// Verificar usuario blacklisted
	isUserBlacklisted, userEntry := database.IsUserBlacklisted(userID)
	if isUserBlacklisted {
		embed := &discordgo.MessageEmbed{
			Title:       "🚫 Acceso Denegado",
			Description: "Tu cuenta ha sido añadida a la blacklist y no puedes usar este bot.",
			Color:       0xFF0000,
			Timestamp:   time.Now().Format(time.RFC3339),
		}

		if userEntry != nil && userEntry.Reason != "" {
			embed.Fields = []*discordgo.MessageEmbedField{
				{
					Name:  "Razón",
					Value: userEntry.Reason,
				},
			}
		}

		ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})

		logger.Warn(fmt.Sprintf("Usuario blacklisted intentó usar comando: %s", userID), "BlacklistMiddleware")
		return fmt.Errorf("user is blacklisted")
	}

	// Verificar guild blacklisted
	if guildID != "" {
		isGuildBlacklisted, guildEntry := database.IsGuildBlacklisted(guildID)
		if isGuildBlacklisted {
			embed := &discordgo.MessageEmbed{
				Title:       "🚫 Servidor en Blacklist",
				Description: "Este servidor ha sido añadido a la blacklist. El bot se retirará automáticamente.",
				Color:       0xFF0000,
				Timestamp:   time.Now().Format(time.RFC3339),
			}

			if guildEntry != nil && guildEntry.Reason != "" {
				embed.Fields = []*discordgo.MessageEmbedField{
					{
						Name:  "Razón",
						Value: guildEntry.Reason,
					},
				}
			}

			ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
					Flags:  discordgo.MessageFlagsEphemeral,
				},
			})

			logger.Warn(fmt.Sprintf("Servidor blacklisted detectado: %s. Saliendo...", guildID), "BlacklistMiddleware")

			// Salir del servidor después de un pequeño delay
			go func() {
				time.Sleep(2 * time.Second)
				if err := ctx.Session.GuildLeave(guildID); err != nil {
					logger.Error(fmt.Sprintf("Error saliendo del servidor blacklisted %s: %v", guildID, err), "BlacklistMiddleware")
				}
			}()

			return fmt.Errorf("guild is blacklisted")
		}
	}

	return nil
}
