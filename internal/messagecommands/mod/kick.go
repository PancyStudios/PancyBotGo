package mod

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func kickCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionKickMembers) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para expulsar miembros.")
		return err
	}

	userID := ctx.ParseUser(0)
	if userID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un usuario.\nUso: `pan!kick @usuario [razón]`")
		return err
	}

	reason := "Sin razón especificada"
	if len(ctx.Args) > 1 {
		reason = strings.Join(ctx.Args[1:], " ")
	}

	err := ctx.Session.GuildMemberDeleteWithReason(ctx.Message.GuildID, userID, reason)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("No se pudo expulsar al usuario: %v", err))
		return err
	}

	_, err = ctx.ReplySuccess("Usuario Expulsado", fmt.Sprintf("👢 **<@%s>** ha sido expulsado.\n**Razón:** %s", userID, reason))
	return err
}
