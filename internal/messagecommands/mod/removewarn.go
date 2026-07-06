package mod

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func removewarnCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionModerateMembers) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para eliminar advertencias.")
		return err
	}

	userID := ctx.ParseUser(0)
	if userID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un usuario.\nUso: `pan!removewarn @usuario <ID_Advertencia>`")
		return err
	}

	if len(ctx.Args) < 2 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el ID de la advertencia.\nUso: `pan!removewarn @usuario <ID_Advertencia>`")
		return err
	}

	warnID := ctx.Args[1]

	dm := database.GlobalWarnDM
	query := bson.M{"guildId": ctx.Message.GuildID, "userId": userID}

	doc, err := dm.Get(query)
	if err != nil {
		logger.Error(fmt.Sprintf("Error DB RemoveWarn: %v", err), "CMD-RemoveWarn")
		_, err = ctx.ReplyError("Error", "❌ Error al consultar la base de datos.")
		return err
	}

	if doc == nil || len(doc.Warns) == 0 {
		_, err = ctx.ReplyError("Error", "❌ El usuario no tiene advertencias.")
		return err
	}

	found := false
	var updatedWarns []models.Warn
	var removedWarn models.Warn

	for _, warn := range doc.Warns {
		if warn.ID == warnID {
			removedWarn = warn
			found = true
		} else {
			updatedWarns = append(updatedWarns, warn)
		}
	}

	if !found {
		_, err = ctx.ReplyError("Error", "❌ No se encontró una advertencia con ese ID.")
		return err
	}

	doc.Warns = updatedWarns
	_, err = dm.Set(query, doc)
	if err != nil {
		logger.Error(fmt.Sprintf("Error guardando RemoveWarn: %v", err), "CMD-RemoveWarn")
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ No se pudo eliminar la advertencia.\nError: `%v`", err))
		return err
	}

	embedSuccess := &discordgo.MessageEmbed{
		Title:       "✅ Advertencia eliminada con éxito",
		Description: fmt.Sprintf("La advertencia de **<@%s>** ha sido eliminada.\n\n**Razón original:** %s\n**ID:** `%s`", userID, removedWarn.Reason, warnID),
		Color:       0x00FF00, // Green
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Solicitado por %s", ctx.Message.Author.String()),
			IconURL: ctx.Message.Author.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	ctx.ReplyEmbed(embedSuccess)

	guild, _ := ctx.Session.State.Guild(ctx.Message.GuildID)

	embedDM := &discordgo.MessageEmbed{
		Title: "ℹ - Advertencia eliminada",
		Color: 0x00FF00,
		Description: fmt.Sprintf(
			"⚒ - **Servidor:** %s (%s)\n"+
				"🗑 ️ - **Advertencia eliminada:** %s\n\n"+
				"🕒 - **Fecha:** <t:%d:F>",
			guild.Name, guild.ID, removedWarn.Reason, time.Now().Unix(),
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
