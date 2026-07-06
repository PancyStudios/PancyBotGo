package mod

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func nukeCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionManageChannels | discordgo.PermissionManageMessages) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permisos de Administrar Canales y Mensajes.")
		return err
	}

	channelID := ctx.Message.ChannelID
	if len(ctx.Args) > 0 {
		parsedChannel := ctx.ParseChannel(0)
		if parsedChannel != "" {
			channelID = parsedChannel
		}
	}

	channel, err := ctx.Session.Channel(channelID)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error obteniendo la información del canal: %v", err))
		return err
	}

	createData := discordgo.GuildChannelCreateData{
		Name:                 channel.Name,
		Type:                 channel.Type,
		Topic:                channel.Topic,
		Bitrate:              channel.Bitrate,
		UserLimit:            channel.UserLimit,
		RateLimitPerUser:     channel.RateLimitPerUser,
		Position:             channel.Position,
		PermissionOverwrites: channel.PermissionOverwrites,
		ParentID:             channel.ParentID,
		NSFW:                 channel.NSFW,
	}

	newChannel, err := ctx.Session.GuildChannelCreateComplex(ctx.Message.GuildID, createData)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error al intentar clonar el canal: %v", err))
		return err
	}

	_, err = ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		_, err = ctx.ReplyError("Advertencia", fmt.Sprintf("⚠️ El canal fue clonado, pero hubo un error al borrar el original: %v", err))
		return err
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Canal Eliminado con Éxito",
		Description: fmt.Sprintf("El canal <#%s> ha sido eliminado y recreado con éxito. Todo su contenido ha sido eliminado.", newChannel.ID),
		Color:       0x00FF00,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    ctx.Message.Author.Username,
			IconURL: ctx.Message.Author.AvatarURL(""),
		},
	}

	_, _ = ctx.Session.ChannelMessageSendEmbed(newChannel.ID, embed)
	return nil
}
