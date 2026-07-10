package economy

import (
	"fmt"
	"strconv"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func buyCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el ID del objeto a comprar.\nUso: `pan!buy <id> [cantidad]`")
		return err
	}

	itemID := ctx.Args[0]
	qty := int64(1)

	if len(ctx.Args) > 0 {
		q, err := strconv.ParseInt(ctx.Args[0], 10, 64)
		if err == nil && q > 0 {
			qty = q
		}
	}

	if qty <= 0 {
		_, err := ctx.ReplyError("Error", "❌ La cantidad debe ser mayor a 0.")
		return err
	}

	items, err := database.GetItems(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error al acceder a la tienda.")
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
		_, err = ctx.ReplyError("Error", "❌ No existe un objeto con ese ID en esta tienda.")
		return err
	}

	totalCost := selectedItem.Price * qty

	if selectedItem.IsGlobal {
		_, err = database.AddStars(ctx.Message.Author.ID, -totalCost, false)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ No tienes suficientes estrellas para comprar esto.")
			return err
		}

		profile, _ := database.GetGlobalProfile(ctx.Message.Author.ID)
		profile.Inventory[selectedItem.ID] += int(qty)
		database.GlobalEconomyDM.Set(bson.M{"_id": profile.UserID}, profile)

		_, err = ctx.ReplySuccess("Compra Exitosa", fmt.Sprintf("Has comprado **x%d %s** por 🌟 %d estrellas.", qty, selectedItem.Name, totalCost))
		return err
	} else {
		_, err = database.AddLocalBalance(ctx.Message.GuildID, ctx.Message.Author.ID, -totalCost, false)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ No tienes suficientes monedas locales para comprar esto.")
			return err
		}

		profile, _ := database.GetLocalProfile(ctx.Message.GuildID, ctx.Message.Author.ID)
		profile.Inventory[selectedItem.ID] += int(qty)
		database.LocalEconomyDM.Set(bson.M{"_id": profile.ID}, profile)

		_, err = ctx.ReplySuccess("Compra Exitosa", fmt.Sprintf("Has comprado **x%d %s** por 💵 %d monedas locales.", qty, selectedItem.Name, totalCost))
		return err
	}
}
