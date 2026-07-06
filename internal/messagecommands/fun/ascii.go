package fun

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
	"github.com/common-nighthawk/go-figure"
)

func asciiCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes proporcionar un texto.\nUso: `pan!ascii <texto>`")
		return err
	}

	texto := strings.Join(ctx.Args, " ")

	myFigure := figure.NewFigure(texto, "", true)
	asciiStr := myFigure.String()

	if len(asciiStr) > 4000 {
		asciiStr = asciiStr[:4000]
	}

	desc := fmt.Sprintf("```\n%s\n```", asciiStr)

	embed := &discordgo.MessageEmbed{
		Title:       "ASCII",
		Color:       0x5865F2,
		Description: desc,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    ctx.Message.Author.Username,
			IconURL: ctx.Message.Author.AvatarURL(""),
		},
	}

	_, err := ctx.ReplyEmbed(embed)
	return err
}
