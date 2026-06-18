package mod

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createNukeCommand() *discord.Command {
	return discord.NewCommand(
		"nuke",
		"Elimina todo el contenido de un canal clonándolo y borrando el original",
		"mod",
		nukeHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "canal",
			Description: "Canal a nuke (opcional, por defecto el actual)",
			Required:    false,
			ChannelTypes: []discordgo.ChannelType{
				discordgo.ChannelTypeGuildText,
				discordgo.ChannelTypeGuildNews,
			},
		},
	).WithUserPermissions(discordgo.PermissionManageChannels | discordgo.PermissionManageMessages).
		WithBotPermissions(discordgo.PermissionManageChannels)
}

func nukeHandler(ctx *discord.CommandContext) error {
	channelOpt := ctx.GetOption("canal")
	channelID := ctx.Interaction.ChannelID
	if channelOpt != nil {
		channelID = channelOpt.Value.(string)
	}

	channel, err := ctx.Session.Channel(channelID)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error obteniendo la información del canal: %v", err))
	}

	// Create new channel payload
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

	newChannel, err := ctx.Session.GuildChannelCreateComplex(ctx.Interaction.GuildID, createData)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error al intentar clonar el canal: %v", err))
	}

	// Delete old channel
	_, err = ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		// If we couldn't delete the old one, we might want to delete the new one or just warn
		return ctx.ReplyEphemeral(fmt.Sprintf("⚠️ El canal fue clonado, pero hubo un error al borrar el original: %v", err))
	}

	// Send success message to the new channel
	embed := &discordgo.MessageEmbed{
		Title:       "Canal Eliminado con Éxito",
		Description: fmt.Sprintf("El canal <#%s> ha sido eliminado y recreado con éxito. Todo su contenido ha sido eliminado.", newChannel.ID),
		Color:       0x00FF00, // Green
		Footer: &discordgo.MessageEmbedFooter{
			Text:    ctx.User().Username,
			IconURL: ctx.User().AvatarURL(""),
		},
	}

	_, _ = ctx.Session.ChannelMessageSendEmbed(newChannel.ID, embed)

	return ctx.ReplyEphemeral(fmt.Sprintf("✅ El canal <#%s> ha sido eliminado y recreado con éxito. Todo su contenido ha sido eliminado.", newChannel.ID))
}
