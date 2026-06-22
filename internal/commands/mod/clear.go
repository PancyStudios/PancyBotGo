package mod

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createClearCommand() *discord.Command {
	return discord.NewCommand(
		"clear",
		"🧹 | Elimina mensajes en un canal (funciona con mensajes de cualquier antigüedad)",
		"mod",
		clearHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "cantidad",
			Description: "Cantidad de mensajes a eliminar (hasta 99999)",
			Required:    true,
			MinValue:    func() *float64 { v := 1.0; return &v }(),
			MaxValue:    99999,
		},
	).WithUserPermissions(discordgo.PermissionManageMessages).
		WithBotPermissions(discordgo.PermissionManageMessages)
}

func clearHandler(ctx *discord.CommandContext) error {
	cantidad := int(ctx.GetIntOption("cantidad"))
	channelID := ctx.Interaction.ChannelID

	// Defer reply because deleting messages can take time
	err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return err
	}

	eliminadosTotal := 0
	var lastId string
	
	// 14 days ago for bulk delete limit
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

		// Bulk delete new messages
		if len(newMessages) > 0 {
			err = ctx.Session.ChannelMessagesBulkDelete(channelID, newMessages)
			if err == nil {
				eliminadosTotal += len(newMessages)
			}
		}

		// Single delete old messages
		for _, msg := range oldMessages {
			err = ctx.Session.ChannelMessageDelete(channelID, msg.ID)
			if err == nil {
				eliminadosTotal++
			}
			time.Sleep(100 * time.Millisecond) // avoid rate limits
			
			if eliminadosTotal >= cantidad {
				break
			}
		}

		if eliminadosTotal >= cantidad {
			break
		}
	}

	content := fmt.Sprintf("✅ Se eliminaron **%d** mensajes del canal.", eliminadosTotal)
	_, err = ctx.Session.InteractionResponseEdit(ctx.Interaction.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	return err
}
