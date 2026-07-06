package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func useCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el ID del objeto a usar.\nUso: `pan!use <id>`")
		return err
	}

	itemID := ctx.Args[0]
	userID := ctx.Message.Author.ID
	guildID := ctx.Message.GuildID

	items, err := database.GetItems(guildID)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error al cargar los objetos.")
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
		_, err = ctx.ReplyError("Error", "❌ No existe ningún objeto con ese ID.")
		return err
	}

	if selectedItem.GuildID == "" {
		profile, _ := database.GetGlobalProfile(userID)
		qty := profile.Inventory[selectedItem.ID]
		if qty <= 0 {
			_, err = ctx.ReplyError("Error", "❌ No tienes ese objeto en tu inventario global.")
			return err
		}

		profile.Inventory[selectedItem.ID] -= 1
		if profile.Inventory[selectedItem.ID] == 0 {
			delete(profile.Inventory, selectedItem.ID)
		}

		efectoAplicado := false
		if selectedItem.Effect == "EXPAND_BANK" {
			profile.BankCapacity += int64(selectedItem.EffectValue)
			efectoAplicado = true
		} else if selectedItem.Effect == "STAR_TICKET" {
			profile.StarsWallet += int64(selectedItem.EffectValue)
			efectoAplicado = true
		}

		database.GlobalEconomyDM.Set(bson.M{"_id": profile.UserID}, profile)

		if efectoAplicado {
			_, err = ctx.ReplySuccess("Objeto Usado", fmt.Sprintf("✨ Has usado **%s**. ¡Tu capacidad estelar o balance ha mejorado!", selectedItem.Name))
			return err
		}

	} else {
		profile, _ := database.GetLocalProfile(guildID, userID)
		qty := profile.Inventory[selectedItem.ID]
		if qty <= 0 {
			_, err = ctx.ReplyError("Error", "❌ No tienes ese objeto en tu inventario local.")
			return err
		}

		profile.Inventory[selectedItem.ID] -= 1
		if profile.Inventory[selectedItem.ID] == 0 {
			delete(profile.Inventory, selectedItem.ID)
		}

		efectoAplicado := false
		if selectedItem.Effect == "EXPAND_BANK" {
			profile.BankCapacity += int64(selectedItem.EffectValue)
			efectoAplicado = true
		} else if selectedItem.Effect == "LOCAL_TICKET" {
			profile.Wallet += int64(selectedItem.EffectValue)
			efectoAplicado = true
		}

		database.LocalEconomyDM.Set(bson.M{"_id": profile.ID}, profile)

		if efectoAplicado {
			_, err = ctx.ReplySuccess("Objeto Usado", fmt.Sprintf("✨ Has usado **%s**. ¡Tu capacidad del banco local o cartera ha mejorado!", selectedItem.Name))
			return err
		}
	}

	_, err = ctx.ReplySuccess("Objeto Usado", fmt.Sprintf("✨ Has usado **%s**, pero no tiene un efecto especial programado.", selectedItem.Name))
	return err
}
