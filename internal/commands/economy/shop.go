package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
)

func createShopCommand() *discord.Command {
	return discord.NewCommand(
		"view",
		"🛒 | Explora el mercado intergaláctico",
		"economy",
		func(ctx *discord.CommandContext) error {
			return shopHandler(ctx)
		},
	)
}

func shopHandler(ctx *discord.CommandContext) error {
	embed, components := ShopMenu()

	err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})

	if err != nil {
		ctx.Reply("❌ No se pudo mostrar el menú de la tienda.")
		return err
	}
	return nil
}

// ShopMenu returns the main shop menu components
func ShopMenu() (*discordgo.MessageEmbed, []discordgo.MessageComponent) {
	embed := discord.NewEmbed().
		SetTitle("🛒 Mercado").
		SetColor(0x9B59B6).
		SetDescription("💰 | Elige qué tienda deseas explorar.\n\n🌟 **Mercado Global:** Usa Estrellas.\n🏪 **Tienda Local:** Usa Monedas del Servidor.").
		Build()

	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "🌟 Tienda Global",
				Style:    discordgo.PrimaryButton,
				CustomID: "shop_nav_global_0",
			},
			discordgo.Button{
				Label:    "🏪 Tienda Local",
				Style:    discordgo.SecondaryButton,
				CustomID: "shop_nav_local_0",
			},
		},
	}
	return embed, []discordgo.MessageComponent{actionRow}
}

// RenderShopPage renders a specific page of the shop
func RenderShopPage(guildID, shopType string, page int) (*discordgo.MessageEmbed, []discordgo.MessageComponent, error) {
	items, err := database.GetItems(guildID)
	if err != nil {
		return nil, nil, err
	}

	var filteredItems []models.Item
	for _, item := range items {
		isGlobalItem := item.GuildID == "" || item.GuildID == "global"
		if shopType == "global" && isGlobalItem {
			filteredItems = append(filteredItems, item)
		} else if shopType == "local" && !isGlobalItem {
			filteredItems = append(filteredItems, item)
		}
	}

	if len(filteredItems) == 0 {
		var title string
		if shopType == "global" {
			title = "🌟 Mercado Global Vacío"
		} else {
			title = "🏪 Tienda Local Vacía"
		}
		
		embed := discord.NewEmbed().
			SetTitle(title).
			SetColor(0xE74C3C).
			SetDescription("No hay objetos disponibles en esta tienda en este momento.").
			Build()

		actionRow := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "🏠 Volver al Menú",
					Style:    discordgo.SecondaryButton,
					CustomID: "shop_nav_menu",
				},
			},
		}
		return embed, []discordgo.MessageComponent{actionRow}, nil
	}

	itemsPerPage := 5
	totalPages := len(filteredItems) / itemsPerPage
	if len(filteredItems)%itemsPerPage != 0 {
		totalPages++
	}
	
	if page < 0 {
		page = 0
	} else if page >= totalPages {
		page = totalPages - 1
	}

	startIdx := page * itemsPerPage
	endIdx := startIdx + itemsPerPage
	if endIdx > len(filteredItems) {
		endIdx = len(filteredItems)
	}

	var list string
	for _, item := range filteredItems[startIdx:endIdx] {
		var currency string
		if shopType == "global" {
			currency = "🌟"
		} else {
			currency = "🪙"
		}
		list += fmt.Sprintf("### %s %s\n> %s\n> **Precio:** `%d %s` | **ID:** `%s`\n\n", item.Emoji, item.Name, item.Description, item.Price, currency, item.ID)
	}

	var title string
	var color int
	if shopType == "global" {
		title = "🌌 Mercado Intergaláctico (Global)"
		color = 0x9b59b6 // Purple
	} else {
		title = "🏪 Tienda Local del Servidor"
		color = 0x2ecc71 // Emerald Green
	}

	embed := discord.NewEmbed().
		SetTitle(fmt.Sprintf("%s (Página %d/%d)", title, page+1, totalPages)).
		SetColor(color).
		SetDescription("Usa `/shop buy <id> [cantidad]` para comprar objetos. Usa `/inventory` para ver tus pertenencias.\n\n" + list).
		Build()

	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "◀️ Anterior",
				Style:    discordgo.PrimaryButton,
				CustomID: fmt.Sprintf("shop_nav_%s_%d", shopType, page-1),
				Disabled: page <= 0,
			},
			discordgo.Button{
				Label:    "🏠 Menú",
				Style:    discordgo.SecondaryButton,
				CustomID: "shop_nav_menu",
			},
			discordgo.Button{
				Label:    "Siguiente ▶️",
				Style:    discordgo.PrimaryButton,
				CustomID: fmt.Sprintf("shop_nav_%s_%d", shopType, page+1),
				Disabled: page >= totalPages-1,
			},
		},
	}

	return embed, []discordgo.MessageComponent{actionRow}, nil
}
