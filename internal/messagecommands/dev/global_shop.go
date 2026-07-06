package dev

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/google/uuid"
)

func globalShopCommand(ctx *messagecommands.MessageContext) error {
	if !isDev(ctx.Message.Author.ID) {
		_, err := ctx.ReplyError("Acceso Denegado", "Este comando es solo para la desarrolladora.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!globalshop add <nombre> | <desc> | <precio> | <efecto> | <valor> | [emoji]`\nO `pan!globalshop remove <id>`")
		return err
	}

	action := strings.ToLower(ctx.Args[0])

	if action == "add" {
		fullArgs := strings.Join(ctx.Args[1:], " ")
		parts := strings.Split(fullArgs, "|")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		if len(parts) < 5 {
			_, err := ctx.ReplyError("Uso Incorrecto", "Faltan argumentos para añadir un objeto global.")
			return err
		}

		name := parts[0]
		desc := parts[1]
		price, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil || price <= 0 {
			_, err = ctx.ReplyError("Error", "Precio inválido.")
			return err
		}

		effect := parts[3]
		effectValue, _ := strconv.ParseFloat(parts[4], 64)

		emoji := "📦"
		if len(parts) > 5 && parts[5] != "" {
			emoji = parts[5]
		}

		item := models.Item{
			ID:          uuid.New().String()[:8],
			GuildID:     "",
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

		err = database.SaveItem(item)
		if err != nil {
			_, err = ctx.ReplyError("Error", "Error al guardar el objeto.")
			return err
		}

		_, err = ctx.ReplySuccess("Objeto Añadido", fmt.Sprintf("✅ Objeto global creado exitosamente.\n**Nombre:** %s\n**Precio:** %d\n**ID:** `%s`", name, price, item.ID))
		return err

	} else if action == "remove" {
		if len(ctx.Args) < 2 {
			_, err := ctx.ReplyError("Uso Incorrecto", "Falta el ID del objeto.")
			return err
		}

		id := ctx.Args[1]
		err := database.DeleteItem(id)
		if err != nil {
			_, err = ctx.ReplyError("Error", "Error al eliminar el objeto.")
			return err
		}

		_, err = ctx.ReplySuccess("Objeto Eliminado", fmt.Sprintf("✅ El objeto con ID `%s` fue eliminado.", id))
		return err
	}

	_, err := ctx.ReplyError("Uso Incorrecto", "Acción no reconocida.")
	return err
}
