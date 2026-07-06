package mod

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/scheduler"
	"github.com/bwmarrin/discordgo"
)

func tempbanCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionBanMembers) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para banear miembros.")
		return err
	}

	userID := ctx.ParseUser(0)
	if userID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un usuario.\nUso: `pan!tempban @usuario <horas> [razón]`")
		return err
	}

	if len(ctx.Args) < 2 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar la duración en horas.\nUso: `pan!tempban @usuario <horas> [razón]`")
		return err
	}

	horas, err := strconv.Atoi(ctx.Args[1])
	if err != nil || horas < 1 {
		_, err := ctx.ReplyError("Uso Incorrecto", "La duración debe ser un número entero mayor a 0.")
		return err
	}

	duracion := time.Duration(horas) * time.Hour

	reason := "Sin razón especificada"
	if len(ctx.Args) > 2 {
		reason = strings.Join(ctx.Args[2:], " ")
	}

	err = ctx.Session.GuildBanCreateWithReason(
		ctx.Message.GuildID,
		userID,
		reason,
		0,
	)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error al banear: %v", err))
		return err
	}

	err = scheduler.AddTempBan(ctx.Message.GuildID, userID, duracion)
	if err != nil {
		_, err = ctx.ReplyError("Advertencia", fmt.Sprintf("⚠️ El usuario fue baneado, pero hubo un error al programar su desbaneo: %v", err))
		return err
	}

	_, err = ctx.ReplySuccess("Usuario Baneado Temporalmente", fmt.Sprintf("⏳ **<@%s>** ha sido baneado temporalmente por %d horas.\n**Razón:** %s", userID, horas, reason))
	return err
}
