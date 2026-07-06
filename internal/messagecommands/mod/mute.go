package mod

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func muteCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionModerateMembers) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para silenciar miembros.")
		return err
	}

	userID := ctx.ParseUser(0)
	if userID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un usuario.\nUso: `pan!mute @usuario <minutos> [razón]`")
		return err
	}

	if len(ctx.Args) < 2 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar la duración en minutos.\nUso: `pan!mute @usuario <minutos> [razón]`")
		return err
	}

	duration, err := strconv.Atoi(ctx.Args[1])
	if err != nil || duration < 1 || duration > 40320 {
		_, err := ctx.ReplyError("Uso Incorrecto", "La duración debe ser un número entre 1 y 40320 minutos (28 días).")
		return err
	}

	reason := "Sin razón especificada"
	if len(ctx.Args) > 2 {
		reason = strings.Join(ctx.Args[2:], " ")
	}

	timeoutUntil := time.Now().Add(time.Duration(duration) * time.Minute)

	err = ctx.Session.GuildMemberTimeout(
		ctx.Message.GuildID,
		userID,
		&timeoutUntil,
	)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error al silenciar: %v", err))
		return err
	}

	_, err = ctx.ReplySuccess("Usuario Silenciado", fmt.Sprintf("🔇 **<@%s>** ha sido silenciado por %d minutos.\n**Razón:** %s", userID, duration, reason))
	return err
}
