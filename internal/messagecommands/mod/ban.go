package mod

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func banCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionBanMembers) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para banear miembros.")
		return err
	}

	userID := ctx.ParseUser(0)
	if userID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un usuario.\nUso: `pan!ban @usuario [razón]`")
		return err
	}

	reason := "Sin razón especificada"
	if len(ctx.Args) > 1 {
		reason = strings.Join(ctx.Args[1:], " ")
	}

	err := ctx.Session.GuildBanCreateWithReason(ctx.Message.GuildID, userID, reason, 0)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("No se pudo banear al usuario: %v", err))
		return err
	}

	_, err = ctx.ReplySuccess("Usuario Baneado", fmt.Sprintf("🔨 **<@%s>** ha sido baneado.\n**Razón:** %s", userID, reason))
	return err
}
