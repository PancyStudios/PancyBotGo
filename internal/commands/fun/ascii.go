package fun

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"github.com/common-nighthawk/go-figure"
)

func createAsciiCommand() *discord.Command {
	return discord.NewCommand(
		"ascii",
		"🔤 | Muestra un texto ASCII",
		"fun",
		asciiHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "texto",
			Description: "Texto a convertir a ASCII",
			Required:    true,
		},
	)
}

func asciiHandler(ctx *discord.CommandContext) error {
	texto := ctx.GetStringOption("texto")

			myFigure := figure.NewFigure(texto, "", true)
			asciiStr := myFigure.String()

			// Discord limits embed descriptions to 4096 chars, we wrap it in a code block
			if len(asciiStr) > 4000 {
				asciiStr = asciiStr[:4000]
			}
			
			desc := fmt.Sprintf("```\n%s\n```", asciiStr)

			embed := &discordgo.MessageEmbed{
				Title:       "ASCII",
				Color:       0x5865F2,
				Description: desc,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    ctx.User().Username,
					IconURL: ctx.User().AvatarURL(""),
				},
			}

			err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando ascii: %v", err), "Fun")
			}
			return nil
}
