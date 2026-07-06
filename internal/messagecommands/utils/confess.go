package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func confessCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		// Lo mandamos como un DM temporal para el uso
		ctx.ReplyError("Uso Incorrecto", "Debes escribir tu confesión secreta.\nUso: `pan!confess <confesion>`")
		return nil
	}

	confesion := strings.Join(ctx.Args, " ")

	guildData, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Message.GuildID})
	if err != nil || guildData == nil || guildData.Configuration.SubData.ConfessionChannel == "" {
		ctx.ReplyError("Error", "❌ El sistema de confesiones no está configurado en este servidor.")
		return nil
	}

	confessChannel := guildData.Configuration.SubData.ConfessionChannel

	// ELIMINAR EL MENSAJE INMEDIATAMENTE PARA MANTENER EL ANONIMATO
	err = ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("No se pudo borrar mensaje de confesion: %v", err), "Confess")
		ctx.ReplyError("Advertencia", "⚠️ Tu confesión fue enviada, pero no pude borrar tu mensaje original. Verifica que tenga permisos de `Administrar Mensajes`.")
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🕵️ Nueva Confesión Anónima",
		Description: confesion,
		Color:       0x9B59B6, // Purple
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(confessChannel, embed)
	if err != nil {
		logger.Error(fmt.Sprintf("Error enviando confesion: %v", err), "Confess")
		// Mandar un mensaje directo al usuario informando del error
		dmChannel, err := ctx.Session.UserChannelCreate(ctx.Message.Author.ID)
		if err == nil {
			ctx.Session.ChannelMessageSend(dmChannel.ID, "❌ No pude enviar la confesión. Verifica que tengo permisos en el canal configurado.")
		}
		return nil
	}

	// Mandar un DM confirmando el exito
	dmChannel, err := ctx.Session.UserChannelCreate(ctx.Message.Author.ID)
	if err == nil {
		ctx.Session.ChannelMessageSend(dmChannel.ID, "✅ Tu confesión ha sido enviada de forma completamente anónima 🤫.")
	}

	return nil
}
