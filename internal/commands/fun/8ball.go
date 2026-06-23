package fun

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

func create8BallCommand() *discord.Command {
	return discord.NewCommand(
		"8ball",
		"🎱 | Pregúntale algo a la bola mágica",
		"fun",
		eightBallHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "pregunta",
			Description: "🎉 | Pregunta que quieres hacerle al bot",
			Required:    true,
		},
	)
}

func eightBallHandler(ctx *discord.CommandContext) error {
	pregunta := ctx.GetStringOption("pregunta")

			respuestas := []string{
				"Sí.",
				"No.",
				"Tal vez.",
				"Probablemente.",
				"Probablemente no.",
				"No sé.",
				"¿Tú qué piensas?",
				"Es cierto.",
				"Definitivamente.",
				"No cuentes con ello.",
			}

			rand.Seed(time.Now().UnixNano())
			idx := rand.Intn(len(respuestas))
			randomResp := respuestas[idx]

			embed := &discordgo.MessageEmbed{
				Title: "🎱 8Ball",
				Color: 0x5865F2, // Blurple
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Pregunta",
						Value: pregunta,
					},
					{
						Name:  "Respuesta",
						Value: randomResp,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    ctx.User().Username,
					IconURL: ctx.User().AvatarURL(""),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando 8ball: %v", err), "Fun")
			}
			return nil
}
