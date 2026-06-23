package fun

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

func createPPTCommand() *discord.Command {
	return discord.NewCommand(
		"ppt",
		"✂️ | Juega Piedra, Papel o Tijera",
		"fun",
		pptHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "move",
			Description: "🎉 | Elige una opción",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{Name: "Piedra", Value: "piedra"},
				{Name: "Papel", Value: "papel"},
				{Name: "Tijera", Value: "tijera"},
			},
		},
	)
}

func pptHandler(ctx *discord.CommandContext) error {
	action := ctx.GetStringOption("move")

			moves := map[string]int{"piedra": 0, "papel": 1, "tijera": 2}
			moveVals := []string{"piedra", "papel", "tijera"}

			if _, ok := moves[action]; !ok {
				action = "piedra"
			}

			rand.Seed(time.Now().UnixNano())
			machineInput := moveVals[rand.Intn(3)]

			i := moves[action]
			j := moves[machineInput]
			
			// determine winner: 0 = tie, 1 = user wins, 2 = machine wins
			winner := ((i - j) % 3 + 3) % 3

			// member color approximation
			color := 0x5865F2

			embed := &discordgo.MessageEmbed{
				Title: "Piedra, papel o tijera",
				Color: color,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   fmt.Sprintf("%s eligió", ctx.User().Username),
						Value:  strings.Title(action),
						Inline: true,
					},
					{
						Name:   "Computadora eligió",
						Value:  strings.Title(machineInput),
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    ctx.User().Username,
					IconURL: ctx.User().AvatarURL(""),
				},
			}

			if winner == 0 {
				embed.Description = "¡Vaya, hubo un empate!"
			} else if winner == 1 {
				embed.Description = "¡Has ganado, felicidades!"
			} else if winner == 2 {
				embed.Description = "¡La computadora ha ganado, suerte para la próxima!"
			}

			err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando ppt: %v", err), "Fun")
			}
			return nil
}
