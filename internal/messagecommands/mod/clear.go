package mod

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func clearCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionManageMessages) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para gestionar mensajes.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar la cantidad de mensajes a eliminar.\nUso: `pan!clear <cantidad>`")
		return err
	}

	cantidad, err := strconv.Atoi(ctx.Args[0])
	if err != nil || cantidad <= 0 {
		_, err = ctx.ReplyError("Error", "❌ La cantidad debe ser un número mayor a 0.")
		return err
	}
	if cantidad > 99999 {
		cantidad = 99999
	}

	channelID := ctx.Message.ChannelID

	eliminadosTotal := 0
	var lastId string

	catorceDias := time.Now().Add(-14 * 24 * time.Hour)

	for eliminadosTotal < cantidad {
		limit := cantidad - eliminadosTotal
		if limit > 100 {
			limit = 100
		}

		messages, err := ctx.Session.ChannelMessages(channelID, limit, lastId, "", "")
		if err != nil || len(messages) == 0 {
			break
		}

		lastId = messages[len(messages)-1].ID

		var newMessages []string
		var oldMessages []*discordgo.Message

		for _, msg := range messages {
			msgTime, _ := discordgo.SnowflakeTimestamp(msg.ID)
			if msgTime.After(catorceDias) {
				newMessages = append(newMessages, msg.ID)
			} else {
				oldMessages = append(oldMessages, msg)
			}
		}

		if len(newMessages) > 0 {
			err = ctx.Session.ChannelMessagesBulkDelete(channelID, newMessages)
			if err == nil {
				eliminadosTotal += len(newMessages)
			}
		}

		for _, msg := range oldMessages {
			err = ctx.Session.ChannelMessageDelete(channelID, msg.ID)
			if err == nil {
				eliminadosTotal++
			}
			time.Sleep(100 * time.Millisecond)

			if eliminadosTotal >= cantidad {
				break
			}
		}

		if eliminadosTotal >= cantidad {
			break
		}
	}

	msg, err := ctx.ReplySuccess("Limpieza Completa", fmt.Sprintf("✅ Se eliminaron %d mensajes.", eliminadosTotal))
	if err == nil {
		go func() {
			time.Sleep(5 * time.Second)
			ctx.Session.ChannelMessageDelete(msg.ChannelID, msg.ID)
		}()
	}
	return err
}
