package mod

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func softbanCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionBanMembers) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para banear miembros.")
		return err
	}

	userID := ctx.ParseUser(0)
	if userID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un usuario.\nUso: `pan!softban @usuario [razón]`")
		return err
	}

	reason := "Sin razón especificada (Softban)"
	if len(ctx.Args) > 1 {
		reason = strings.Join(ctx.Args[1:], " ") + " (Softban)"
	}

	// 7 days of messages deleted
	err := ctx.Session.GuildBanCreateWithReason(ctx.Message.GuildID, userID, reason, 7)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error al banear: %v", err))
		return err
	}

	// Immediately unban
	err = ctx.Session.GuildBanDelete(ctx.Message.GuildID, userID)
	if err != nil {
		_, err = ctx.ReplyError("Advertencia", fmt.Sprintf("⚠️ El usuario fue baneado, pero hubo un error al desbanearlo: %v", err))
		return err
	}

	_, err = ctx.ReplySuccess("Softban Aplicado", fmt.Sprintf("🧹 **<@%s>** ha sido softbaneado (sus mensajes recientes fueron eliminados).\n**Razón:** %s", userID, reason))
	return err
}
