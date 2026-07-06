package help

import (
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func cmdsCommand(ctx *messagecommands.MessageContext) error {
	commands := messagecommands.GetRegisteredCommands()

	if len(commands) == 0 {
		_, err := ctx.ReplyError("Error", "❌ No hay comandos registrados actualmente.")
		return err
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📚 Lista de Comandos de PancyBot",
		Description: "Aquí tienes todos los comandos de prefijo disponibles.",
		Color:       0x3498DB,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ctx.Session.State.User.AvatarURL("128"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "💠 Comandos Disponibles",
				Value: "```\n" + strings.Join(commands, ", ") + "\n```",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Usa el prefijo de tu servidor o menciona al bot antes de cada comando.",
		},
	}

	_, err := ctx.ReplyEmbed(embed)
	return err
}
