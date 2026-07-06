package economy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func adminShopCommand(ctx *messagecommands.MessageContext, isGlobal bool) error {
	if !ctx.HasPermission(discordgo.PermissionAdministrator) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para administrar la tienda local.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!adminshop add <nombre> | <desc> | <precio> | [emoji] | [efecto: NONE/EXPAND_BANK/GIVE_ROLE] | [valor_efecto] | [rol]`\nO `pan!adminshop delete <id>`")
		return err
	}

	action := strings.ToLower(ctx.Args[0])

	if action == "add" {
		fullArgs := strings.Join(ctx.Args[1:], " ")
		parts := strings.Split(fullArgs, "|")

		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		if len(parts) < 3 {
			_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!adminshop add <nombre> | <desc> | <precio> | [emoji] | [efecto] | [valor] | [rol]`\nSepara los argumentos con el carácter `|`.")
			return err
		}

		name := parts[0]
		desc := parts[1]

		price, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil || price <= 0 {
			_, err = ctx.ReplyError("Error", "❌ El precio debe ser un número mayor a 0.")
			return err
		}

		emoji := "📦"
		if len(parts) > 3 && parts[3] != "" {
			emoji = parts[3]
		}

		effect := "NONE"
		if len(parts) > 4 && parts[4] != "" {
			e := strings.ToUpper(parts[4])
			if e == "NONE" || e == "EXPAND_BANK" || e == "GIVE_ROLE" {
				effect = e
			}
		}

		effectValue := float64(0)
		if len(parts) > 5 && parts[5] != "" {
			val, err := strconv.ParseFloat(parts[5], 64)
			if err == nil {
				effectValue = val
			}
		}

		roleID := ""
		if len(parts) > 6 && parts[6] != "" {
			r := parts[6]
			r = strings.TrimPrefix(r, "<@&")
			r = strings.TrimSuffix(r, ">")
			roleID = r
		}

		itemType := models.ItemTypeCollectible
		if effect == "GIVE_ROLE" {
			itemType = models.ItemTypeRole
		} else if effect != "NONE" {
			itemType = models.ItemTypeConsumable
		}

		item := models.Item{
			ID:          uuid.New().String()[:8],
			GuildID:     ctx.Message.GuildID,
			Name:        name,
			Description: desc,
			Price:       price,
			SellPrice:   price / 2,
			Type:        itemType,
			Emoji:       emoji,
			Stock:       -1,
			Effect:      effect,
			EffectValue: effectValue,
			RoleID:      roleID,
		}

		err = database.SaveItem(item)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ Hubo un error al guardar el objeto en la tienda local.")
			return err
		}

		_, err = ctx.ReplySuccess("Objeto Creado", fmt.Sprintf("✅ Objeto local creado exitosamente.\n**Nombre:** %s\n**Precio:** %d\n**ID:** `%s`", name, price, item.ID))
		return err

	} else if action == "delete" {
		if len(ctx.Args) < 1 {
			_, err := ctx.ReplyError("Uso Incorrecto", "Debes proporcionar el ID del objeto a eliminar.\nUso: `pan!adminshop delete <id>`")
			return err
		}

		id := ctx.Args[0]

		items, err := database.GetItems(ctx.Message.GuildID)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ Error al buscar el catálogo.")
			return err
		}

		found := false
		for _, it := range items {
			if it.ID == id && it.GuildID == ctx.Message.GuildID {
				found = true
				break
			}
		}

		if !found {
			_, err = ctx.ReplyError("Error", "❌ No se encontró un objeto local con esa ID en este servidor.")
			return err
		}

		err = database.DeleteItem(id)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ Hubo un error al eliminar el objeto.")
			return err
		}

		_, err = ctx.ReplySuccess("Objeto Eliminado", fmt.Sprintf("✅ El objeto con ID `%s` fue eliminado de la tienda del servidor.", id))
		return err
	}

	_, err := ctx.ReplyError("Uso Incorrecto", "Acción no reconocida. Usa `add` o `delete`.")
	return err
}
