package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func createUseCommand() *discord.Command {
	return discord.NewCommand(
		"use",
		"✨ | Usa un objeto mágico de tu inventario",
		"economy",
		useHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "id",
			Description: "ID del objeto a usar",
			Required:    true,
		},
	)
}

func useHandler(ctx *discord.CommandContext) error {
	itemID := ctx.GetStringOption("id")
	userID := ctx.Interaction.Member.User.ID
	guildID := ctx.Interaction.GuildID

	items, err := database.GetItems(guildID)
	if err != nil {
		ctx.Reply("❌ Error al cargar los objetos.")
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
		ctx.Reply("❌ No existe ningún objeto con ese ID.")
		return nil
	}

	if selectedItem.GuildID == "" {
		// Global Item
		profile, _ := database.GetGlobalProfile(userID)
		qty := profile.Inventory[selectedItem.ID]
		if qty <= 0 {
			ctx.Reply("❌ No tienes ese objeto en tu inventario global.")
			return nil
		}
		
		// Use it
		profile.Inventory[selectedItem.ID] -= 1
		if profile.Inventory[selectedItem.ID] == 0 {
			delete(profile.Inventory, selectedItem.ID)
		}

		// Apply Global Effect
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
			ctx.Reply(fmt.Sprintf("✨ Has usado **%s**. ¡Tu capacidad estelar o balance ha mejorado!", selectedItem.Name))
			return nil
		}
		
	} else {
		// Local Item
		profile, _ := database.GetLocalProfile(guildID, userID)
		qty := profile.Inventory[selectedItem.ID]
		if qty <= 0 {
			ctx.Reply("❌ No tienes ese objeto en tu inventario local.")
			return nil
		}

		// Use it
		profile.Inventory[selectedItem.ID] -= 1
		if profile.Inventory[selectedItem.ID] == 0 {
			delete(profile.Inventory, selectedItem.ID)
		}

		if selectedItem.Type == models.ItemTypeRole && selectedItem.RoleID != "" {
			err = ctx.Session.GuildMemberRoleAdd(guildID, userID, selectedItem.RoleID)
			if err != nil {
				ctx.Reply("❌ No pude darte el rol asociado a este objeto. Revisa mis permisos.")
				return nil
			}
			ctx.Reply(fmt.Sprintf("✨ Has usado **%s** y obtuviste un rol nuevo.", selectedItem.Name))
			database.LocalEconomyDM.Set(bson.M{"_id": profile.ID}, profile)
			return nil
		}

		database.LocalEconomyDM.Set(bson.M{"_id": profile.ID}, profile)
	}

	ctx.Reply(fmt.Sprintf("✨ Has usado **%s**. ¡El efecto ha sido aplicado!", selectedItem.Name))
	return nil
}
