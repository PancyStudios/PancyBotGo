package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
)

func createInventoryCommand() *discord.Command {
	return discord.NewCommand(
		"inventory",
		"🎒 | Revisa tu inventario de objetos",
		"economy",
		inventoryHandler,
	)
}

func inventoryHandler(ctx *discord.CommandContext) error {
	userID := ctx.Interaction.Member.User.ID
	guildID := ctx.Interaction.GuildID

	globalProfile, _ := database.GetGlobalProfile(userID)
	localProfile, _ := database.GetLocalProfile(guildID, userID)

	items, _ := database.GetItems(guildID)
	itemMap := make(map[string]models.Item)
	for _, it := range items {
		itemMap[it.ID] = it
	}

	var globalInvStr string
	if globalProfile != nil {
		for itemID, qty := range globalProfile.Inventory {
			if qty > 0 {
				item, exists := itemMap[itemID]
				if exists {
					globalInvStr += fmt.Sprintf("%s **%s** x%d\n", item.Emoji, item.Name, qty)
				}
			}
		}
	}

	var localInvStr string
	if localProfile != nil {
		for itemID, qty := range localProfile.Inventory {
			if qty > 0 {
				item, exists := itemMap[itemID]
				if exists {
					localInvStr += fmt.Sprintf("%s **%s** x%d\n", item.Emoji, item.Name, qty)
				}
			}
		}
	}

	if globalInvStr == "" {
		globalInvStr = "No tienes objetos estelares."
	}
	if localInvStr == "" {
		localInvStr = "No tienes objetos del servidor."
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🎒 Inventario de %s", ctx.Interaction.Member.User.Username),
		Color:       0x3498DB,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: ctx.Interaction.Member.User.AvatarURL("")},
		Description: "Aquí están todos tus objetos coleccionados.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "🌟 Objetos Globales",
				Value:  globalInvStr,
				Inline: false,
			},
			{
				Name:   "🏠 Objetos Locales",
				Value:  localInvStr,
				Inline: false,
			},
		},
	}

	ctx.ReplyEmbed(embed)
	return nil
}
