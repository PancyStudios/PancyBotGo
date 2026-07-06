package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
)

func inventoryCommand(ctx *messagecommands.MessageContext, isGlobal bool) error {
	targetUser := ctx.Message.Author
	if len(ctx.Args) > 0 {
		parsedUserID := ctx.ParseUser(0)
		if parsedUserID != "" {
			member, err := ctx.Session.GuildMember(ctx.Message.GuildID, parsedUserID)
			if err == nil {
				targetUser = member.User
			}
		}
	}

	userID := targetUser.ID
	guildID := ctx.Message.GuildID

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

	embed := discord.NewEmbed().
		SetTitle(fmt.Sprintf("🎒 Inventario de %s", targetUser.Username)).
		SetColor(0x3498DB).
		SetThumbnail(targetUser.AvatarURL("")).
		SetDescription("💰 | Aquí están todos tus objetos coleccionados.").
		AddField("🌟 Objetos Globales", globalInvStr, false).
		AddField("🏠 Objetos Locales", localInvStr, false).
		Build()

	_, err := ctx.ReplyEmbed(embed)
	return err
}
