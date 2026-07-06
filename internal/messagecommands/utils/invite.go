package utils

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func inviteCommand(ctx *messagecommands.MessageContext) error {
	clientID := ctx.Session.State.User.ID
	inviteURL := fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&permissions=8&scope=bot%%20applications.commands", clientID)

	embed := &discordgo.MessageEmbed{
		Title:       "¡Invítame a tu servidor!",
		Description: fmt.Sprintf("[Haz clic aquí para invitarme](%s)", inviteURL),
		Color:       0x00FF00,
	}

	_, err := ctx.ReplyEmbed(embed)
	return err
}
