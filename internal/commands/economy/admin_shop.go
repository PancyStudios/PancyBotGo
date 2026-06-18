package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func createAdminShopCommand() *discord.Command {
	cmd := discord.NewCommand(
		"shopadmin",
		"🛠️ (Admins) Añade o elimina objetos de la tienda local",
		"economy",
		adminShopHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "Añadir un nuevo objeto local",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "nombre",
					Description: "Nombre del objeto",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "descripcion",
					Description: "Descripción del objeto",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "precio",
					Description: "Precio de compra (en monedas locales)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "emoji",
					Description: "Emoji que representa al objeto",
					Required:    false,
				},
			},
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "delete",
			Description: "Eliminar un objeto local",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "ID del objeto a eliminar",
					Required:    true,
				},
			},
		},
	).WithUserPermissions(discordgo.PermissionAdministrator)
	return cmd
}

func adminShopHandler(ctx *discord.CommandContext) error {
	subcommand := ctx.Interaction.ApplicationCommandData().Options[0]

	if subcommand.Name == "add" {
		name := ""
		desc := ""
		price := int64(0)
		emoji := "📦"

		for _, opt := range subcommand.Options {
			switch opt.Name {
			case "nombre":
				name = opt.StringValue()
			case "descripcion":
				desc = opt.StringValue()
			case "precio":
				price = opt.IntValue()
			case "emoji":
				emoji = opt.StringValue()
			}
		}

		if price <= 0 {
			ctx.Reply("❌ " + "El precio debe ser mayor a 0.")
			return nil
		}

		item := models.Item{
			ID:          uuid.New().String()[:8],
			GuildID:     ctx.Interaction.GuildID,
			Name:        name,
			Description: desc,
			Price:       price,
			SellPrice:   price / 2,
			Type:        models.ItemTypeCollectible,
			Emoji:       emoji,
			Stock:       -1,
		}

		err := database.SaveItem(item)
		if err != nil {
			ctx.Reply("❌ " + "Hubo un error al guardar el objeto en la tienda local.")
			return err
		}

		ctx.Reply(fmt.Sprintf("✅ Objeto local creado exitosamente.\n**Nombre:** %s\n**Precio:** %d\n**ID:** `%s`", name, price, item.ID))

	} else if subcommand.Name == "delete" {
		id := subcommand.Options[0].StringValue()

		items, err := database.GetItems(ctx.Interaction.GuildID)
		if err != nil {
			ctx.Reply("❌ " + "Error al buscar el catálogo.")
			return err
		}

		found := false
		for _, it := range items {
			if it.ID == id && it.GuildID == ctx.Interaction.GuildID {
				found = true
				break
			}
		}

		if !found {
			ctx.Reply("❌ " + "No se encontró un objeto local con esa ID en este servidor.")
			return nil
		}

		err = database.DeleteItem(id)
		if err != nil {
			ctx.Reply("❌ " + "Hubo un error al eliminar el objeto.")
			return err
		}

		ctx.Reply(fmt.Sprintf("✅ El objeto con ID `%s` fue eliminado de la tienda del servidor.", id))
	}
	
	return nil
}
