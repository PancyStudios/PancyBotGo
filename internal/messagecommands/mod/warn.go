package mod

import (
	"fmt"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func warnCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionModerateMembers) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para advertir miembros.")
		return err
	}

	userID := ctx.ParseUser(0)
	if userID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un usuario válido.\nUso: `pan!warn @usuario [razón]`")
		return err
	}

	if userID == ctx.Message.Author.ID {
		_, err := ctx.ReplyError("Error", "❌ No puedes advertirte a ti mismo.")
		return err
	}

	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error obteniendo información del servidor.")
		return err
	}

	if userID == guild.OwnerID {
		_, err := ctx.ReplyError("Error", "❌ No puedes advertir al dueño del servidor.")
		return err
	}

	targetMember, err := ctx.Session.GuildMember(ctx.Message.GuildID, userID)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ No se pudo obtener la información del miembro objetivo.")
		return err
	}

	if targetMember.User.Bot {
		_, err := ctx.ReplyError("Error", "❌ No puedes advertir a un bot.")
		return err
	}

	executorMember, err := ctx.Session.GuildMember(ctx.Message.GuildID, ctx.Message.Author.ID)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ No se pudo obtener tu información de miembro.")
		return err
	}

	executorPosition := getHighestRolePosition(guild, executorMember)
	targetPosition := getHighestRolePosition(guild, targetMember)

	if targetPosition >= executorPosition && ctx.Message.Author.ID != guild.OwnerID {
		_, err := ctx.ReplyError("Error", "❌ No puedes advertir a un usuario con un rol mayor o igual al tuyo.")
		return err
	}

	reason := "Razón no proporcionada"
	if len(ctx.Args) > 1 {
		reason = strings.Join(ctx.Args[1:], " ")
	}

	warnID := uuid.New().String()
	shortID := strings.ReplaceAll(warnID, "-", "")

	newWarn := models.Warn{
		Reason:    reason,
		Moderator: ctx.Message.Author.ID,
		ID:        shortID,
		Timestamp: time.Now().Unix(),
	}

	dm := database.GlobalWarnDM
	query := bson.M{"guildId": ctx.Message.GuildID, "userId": userID}
	doc, err := dm.Get(query)

	if err != nil {
		logger.Error(fmt.Sprintf("Error DB Warn: %v", err), "CMD-Warn")
		_, err = ctx.ReplyError("Error", "❌ Error al acceder a la base de datos.")
		return err
	}

	if doc == nil {
		newDoc := models.WarnsDocument{
			GuildID: ctx.Message.GuildID,
			UserID:  userID,
			Warns:   []models.Warn{newWarn},
		}
		_, err = dm.Set(query, newDoc)
	} else {
		doc.Warns = append(doc.Warns, newWarn)
		_, err = dm.Set(query, doc)
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Error guardando Warn: %v", err), "CMD-Warn")
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ No se pudo guardar la advertencia en la base de datos.\nError: `%v`", err))
		return err
	}

	embedSuccess := &discordgo.MessageEmbed{
		Title:       "✅ Usuario advertido con éxito",
		Description: fmt.Sprintf("El usuario **<@%s>** ha sido advertido correctamente.\n\n**Razón:** %s\n**ID de Advertencia:** `%s`", userID, reason, shortID),
		Color:       0x00FF00, // Green
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Solicitado por %s", ctx.Message.Author.String()),
			IconURL: ctx.Message.Author.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	ctx.ReplyEmbed(embedSuccess)

	embedDM := &discordgo.MessageEmbed{
		Title: "⚠️ - Has recibido una advertencia",
		Color: 0xFFFF00,
		Description: fmt.Sprintf(
			"⚒️ - **Servidor:** %s (%s)\n"+
				"🔨 - **Razón:** %s\n\n"+
				"🕒 - **Fecha:** <t:%d:F>",
			guild.Name, guild.ID, reason, time.Now().Unix(),
		),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "💫 - Developed by PancyStudios",
			IconURL: ctx.Session.State.User.AvatarURL(""),
		},
	}

	userChannel, err := ctx.Session.UserChannelCreate(userID)
	if err == nil {
		ctx.Session.ChannelMessageSendEmbed(userChannel.ID, embedDM)
	}

	return nil
}
