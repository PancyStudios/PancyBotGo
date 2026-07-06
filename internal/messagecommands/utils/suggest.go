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

func suggestCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes escribir una sugerencia.\nUso: `pan!suggest <sugerencia>`")
		return err
	}

	sugerencia := strings.Join(ctx.Args, " ")

	guildData, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Message.GuildID})
	if err != nil || guildData == nil || guildData.Configuration.SubData.SuggestChannel == "" {
		_, err := ctx.ReplyError("Error", "❌ El sistema de sugerencias no está configurado en este servidor.")
		return err
	}

	suggestChannel := guildData.Configuration.SubData.SuggestChannel
	user := ctx.Message.Author

	embed := &discordgo.MessageEmbed{
		Title:       "💡 Nueva Sugerencia",
		Description: sugerencia,
		Color:       0xF1C40F, // Yellow
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.String(),
			IconURL: user.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	msg, err := ctx.Session.ChannelMessageSendEmbed(suggestChannel, embed)
	if err != nil {
		logger.Error(fmt.Sprintf("Error enviando sugerencia: %v", err), "Suggest")
		_, err = ctx.ReplyError("Error", "❌ No pude enviar la sugerencia. Verifica que tengo permisos en el canal configurado.")
		return err
	}

	// Add reactions
	ctx.Session.MessageReactionAdd(suggestChannel, msg.ID, "✅")
	ctx.Session.MessageReactionAdd(suggestChannel, msg.ID, "❌")

	// Eliminar el mensaje del usuario para limpiar el chat
	ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)

	reply, err := ctx.ReplySuccess("Sugerencia Enviada", fmt.Sprintf("✅ Tu sugerencia ha sido enviada a <#%s>", suggestChannel))
	if err == nil {
		// Borrar la respuesta despues de 5 segundos
		go func() {
			time.Sleep(5 * time.Second)
			ctx.Session.ChannelMessageDelete(reply.ChannelID, reply.ID)
		}()
	}

	return nil
}
