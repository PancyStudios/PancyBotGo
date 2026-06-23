package utils

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createInviteCommand() *discord.Command {
	return &discord.Command{
		Name:        "invite",
		Description: "🧰 | Obten el enlace de invitacion del bot",
		Run: func(ctx *discord.CommandContext) error {
			clientID := ctx.Session.State.User.ID
			inviteURL := fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&permissions=8&scope=bot%%20applications.commands", clientID)

			embed := &discordgo.MessageEmbed{
				Title:       "¡Invítame a tu servidor!",
				Description: fmt.Sprintf("[Haz clic aquí para invitarme](%s)", inviteURL),
				Color:       0x00FF00, // Green
			}

			return ctx.ReplyEmbed(embed)
		},
	}
}
