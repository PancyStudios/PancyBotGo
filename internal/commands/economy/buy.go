package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
)

func createBuyCommand() *discord.Command {
	return discord.NewCommand(
		"buy",
		"🛍️ | Compra un objeto del mercado estelar o local",
		"economy",
		func(ctx *discord.CommandContext) error {
			return buyHandler(ctx)
		},
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "id",
			Description: "💰 | El ID del objeto que quieres comprar (revisa /shop)",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "cantidad",
			Description: "💰 | Cuántos quieres comprar",
			Required:    false,
		},
	)
}

func buyHandler(ctx *discord.CommandContext) error {
	itemID := ctx.GetStringOption("id")
	qty := int64(1)

	optQty := ctx.GetIntOption("cantidad")
	if optQty > 0 {
		qty = int64(optQty)
	}

	if qty <= 0 {
		ctx.Reply("❌ " + "La cantidad debe ser mayor a 0.")
		return nil
	}

	items, err := database.GetItems(ctx.Interaction.GuildID)
	if err != nil {
		ctx.Reply("❌ " + "Error al acceder a la tienda.")
		return err
	}

	var selectedItem *models.Item
	for _, it := range items {
		if it.ID == itemID {
			copyIt := it
			selectedItem = &copyIt
			break
		}
	}

	if selectedItem == nil {
		ctx.Reply("❌ " + "No existe un objeto con ese ID en esta tienda.")
		return nil
	}

	if err := database.PurchaseItem(ctx.Interaction.GuildID, ctx.User().ID, *selectedItem, qty); err != nil {
		ctx.Reply("❌ No se pudo completar la compra. Verifica tu saldo e inténtalo de nuevo.")
		return err
	}
	totalCost := selectedItem.Price * qty

	if selectedItem.IsGlobal {
		ctx.Reply(fmt.Sprintf("Has comprado **x%d %s** por 🌟 %d estrellas.", qty, selectedItem.Name, totalCost))
	} else {
		ctx.Reply(fmt.Sprintf("Has comprado **x%d %s** por 💵 %d monedas locales.", qty, selectedItem.Name, totalCost))
	}

	return nil
}
