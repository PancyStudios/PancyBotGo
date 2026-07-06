package dev

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func CreateShopAddCommand() *discord.Command {
	return discord.NewCommand(
		"add",
		"Añade un objeto a la tienda global estelar",
		"dev",
		shopAddHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "nombre",
			Description: "💻 | Nombre del objeto global",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "descripcion",
			Description: "💻 | Descripción del objeto global",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "precio",
			Description: "💻 | Precio de compra (en estrellas)",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "efecto",
			Description: "💻 | Efecto (ej. EXPAND_BANK, STAR_TICKET)",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        "valor_efecto",
			Description: "💻 | Valor asociado al efecto (ej. 1000 para expandir banco en 1000)",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "emoji",
			Description: "💻 | Emoji que representa al objeto",
			Required:    false,
		},
	)
}

func CreateShopRemoveCommand() *discord.Command {
	return discord.NewCommand(
		"remove",
		"✨ | Elimina un objeto de la tienda global",
		"dev",
		shopRemoveHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "id",
			Description: "💻 | ID del objeto global a eliminar",
			Required:    true,
		},
	)
}

func shopAddHandler(ctx *discord.CommandContext) error {
	name := ctx.GetStringOption("nombre")
	desc := ctx.GetStringOption("descripcion")
	price := ctx.GetIntOption("precio")
	effect := ctx.GetStringOption("efecto")
	var effectValue float64
	if opt := ctx.GetOption("valor_efecto"); opt != nil {
		effectValue = opt.FloatValue()
	}

	emoji := "📦"
	if ctx.GetStringOption("emoji") != "" {
		emoji = ctx.GetStringOption("emoji")
	}

	if price <= 0 {
		return ctx.Reply("❌ El precio debe ser mayor a 0.")
	}

	item := models.Item{
		ID:          uuid.New().String()[:8],
		GuildID:     "", // Global item
		Name:        name,
		Description: desc,
		Price:       price,
		SellPrice:   price / 2,
		Type:        models.ItemTypeConsumable,
		Emoji:       emoji,
		Stock:       -1,
		Effect:      effect,
		EffectValue: effectValue,
	}

	err := database.SaveItem(item)
	if err != nil {
		return ctx.Reply("❌ Hubo un error al guardar el objeto global.")
	}

	return ctx.Reply(fmt.Sprintf("✅ Objeto global estelar creado exitosamente.\n**Nombre:** %s\n**Precio:** %d 🌟\n**ID:** `%s`\n**Efecto:** %s (%.2f)", name, price, item.ID, effect, effectValue))
}

func shopRemoveHandler(ctx *discord.CommandContext) error {
	id := ctx.GetStringOption("id")

	items, err := database.GetItems("") // Global items
	if err != nil {
		return ctx.Reply("❌ Error al buscar el catálogo global.")
	}

	found := false
	for _, it := range items {
		if it.ID == id {
			found = true
			break
		}
	}

	if !found {
		return ctx.Reply("❌ No se encontró un objeto global con esa ID.")
	}

	err = database.DeleteItem(id)
	if err != nil {
		return ctx.Reply("❌ Hubo un error al eliminar el objeto global.")
	}

	return ctx.Reply(fmt.Sprintf("✅ El objeto global con ID `%s` fue eliminado del mercado intergaláctico.", id))
}
