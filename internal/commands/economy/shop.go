package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

func createShopCommand() *discord.Command {
	return discord.NewCommand(
		"view",
		"🛒 | Explora el mercado intergaláctico",
		"economy",
		shopHandler,
	)
}

func shopHandler(ctx *discord.CommandContext) error {
	items, err := database.GetItems(ctx.Interaction.GuildID)
	if err != nil {
		ctx.Reply("❌ " + "No se pudo cargar el catálogo de la tienda.")
		return err
	}

	if len(items) == 0 {
		ctx.Reply("🛒 **Tienda Vacía**\nNo hay objetos disponibles en la tienda en este momento.")
		return nil
	}

	var globalItems string
	var localItems string

	for _, item := range items {
		line := fmt.Sprintf("%s **%s** - %d (ID: `%s`)\n* %s\n\n", item.Emoji, item.Name, item.Price, item.ID, item.Description)
		if item.GuildID == "" {
			globalItems += line
		} else {
			localItems += line
		}
	}

	if globalItems == "" {
		globalItems = "No hay objetos estelares."
	}
	if localItems == "" {
		localItems = "El administrador del servidor no ha creado objetos locales."
	}

	embed := discord.NewEmbed().
		SetTitle("🛒 Mercado Intergaláctico").
		SetColor(0x9B59B6).
		SetDescription("💰 | Usa `/shop buy <id>` para comprar un objeto.").
		AddField("🌟 Mercado Global (Cuesta Estrellas)", globalItems, false).
		AddField("🏪 Tienda del Servidor (Cuesta Monedas)", localItems, false).
		Build()

	ctx.ReplyEmbed(embed)
	return nil
}
