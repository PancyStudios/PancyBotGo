package mod

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func lockdownCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionManageChannels) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para gestionar canales.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Especifica `on` para bloquear o `off` para desbloquear.\nUso: `pan!lockdown <on/off>`")
		return err
	}

	estadoArg := strings.ToLower(ctx.Args[0])
	estado := false
	if estadoArg == "on" || estadoArg == "true" {
		estado = true
	} else if estadoArg == "off" || estadoArg == "false" {
		estado = false
	} else {
		_, err := ctx.ReplyError("Uso Incorrecto", "El estado debe ser `on` o `off`.")
		return err
	}

	channelID := ctx.Message.ChannelID

	channel, err := ctx.Session.Channel(channelID)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error obteniendo la información del canal: %v", err))
		return err
	}

	everyoneRoleID := ctx.Message.GuildID

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
		everyoneOverwrite.Deny |= discordgo.PermissionSendMessages
		everyoneOverwrite.Allow &= ^int64(discordgo.PermissionSendMessages)
	} else {
		everyoneOverwrite.Deny &= ^int64(discordgo.PermissionSendMessages)
	}

	err = ctx.Session.ChannelPermissionSet(
		channelID,
		everyoneOverwrite.ID,
		everyoneOverwrite.Type,
		everyoneOverwrite.Allow,
		everyoneOverwrite.Deny,
	)

	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error cambiando permisos: %v", err))
		return err
	}

	if estado {
		_, err = ctx.ReplySuccess("Lockdown Activado", "🔒 **Canal Bloqueado.** Nadie (sin permisos) puede escribir ahora.")
		return err
	}
	_, err = ctx.ReplySuccess("Lockdown Desactivado", "🔓 **Canal Desbloqueado.** Todos pueden volver a escribir.")
	return err
}
