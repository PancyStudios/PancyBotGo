package mod

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createLockdownCommand() *discord.Command {
	return discord.NewCommand(
		"lockdown",
		"Bloquea o desbloquea el canal actual para que nadie pueda escribir",
		"mod",
		lockdownHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "estado",
			Description: "True para bloquear, False para desbloquear",
			Required:    true,
		},
	).WithUserPermissions(discordgo.PermissionManageChannels).
		WithBotPermissions(discordgo.PermissionManageChannels)
}

func lockdownHandler(ctx *discord.CommandContext) error {
	estado := ctx.GetOption("estado").Value.(bool)
	channelID := ctx.Interaction.ChannelID

	channel, err := ctx.Session.Channel(channelID)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error obteniendo la información del canal: %v", err))
	}

	everyoneRoleID := ctx.Interaction.GuildID

	// Find the existing permission overwrite for @everyone
	var everyoneOverwrite *discordgo.PermissionOverwrite
	for _, ow := range channel.PermissionOverwrites {
		if ow.ID == everyoneRoleID {
			everyoneOverwrite = ow
			break
		}
	}

	if everyoneOverwrite == nil {
		everyoneOverwrite = &discordgo.PermissionOverwrite{
			ID:    everyoneRoleID,
			Type:  discordgo.PermissionOverwriteTypeRole,
			Allow: 0,
			Deny:  0,
		}
	}

	if estado {
		// Lockdown: Deny SendMessages
		everyoneOverwrite.Deny |= discordgo.PermissionSendMessages
		everyoneOverwrite.Allow &= ^int64(discordgo.PermissionSendMessages)
	} else {
		// Unlock: Clear SendMessages from Deny
		everyoneOverwrite.Deny &= ^int64(discordgo.PermissionSendMessages)
		// We can also allow it explicitly, but usually clearing Deny is enough
	}

	err = ctx.Session.ChannelPermissionSet(
		channelID,
		everyoneOverwrite.ID,
		everyoneOverwrite.Type,
		everyoneOverwrite.Allow,
		everyoneOverwrite.Deny,
	)

	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error cambiando permisos: %v", err))
	}

	if estado {
		return ctx.Reply("🔒 **Canal Bloqueado.** Nadie (sin permisos) puede escribir ahora.")
	}
	return ctx.Reply("🔓 **Canal Desbloqueado.** Todos pueden volver a escribir.")
}
