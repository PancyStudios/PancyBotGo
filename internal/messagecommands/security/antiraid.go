package security

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func antiraidCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionAdministrator) {
		_, err := ctx.ReplyError("Acceso Denegado", "Necesitas permisos de Administrador para usar este comando.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!antiraid <toggle|age|limits|action> [args]`")
		return err
	}

	subcommand := strings.ToLower(ctx.Args[0])

	guildData, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Message.GuildID})
	if err != nil || guildData == nil {
		_, err = ctx.ReplyError("Error", "❌ Ocurrió un error al cargar la configuración del servidor.")
		return err
	}

	antiRaid := &guildData.Protection.AntiRaid
	if antiRaid.Action == "" {
		antiRaid.Action = "kick"
	}

	var response string

	switch subcommand {
	case "toggle":
		antiRaid.Enable = !antiRaid.Enable
		status := "desactivado"
		if antiRaid.Enable {
			status = "**ACTIVADO**"
		}
		response = fmt.Sprintf("🚨 El sistema Anti-Raid está ahora %s.", status)

	case "age":
		_, err = ctx.ReplyError("No Soportado", "La restricción por edad de cuenta ya no está soportada.")
		return err

	case "limits":
		if len(ctx.Args) < 2 {
			_, err = ctx.ReplyError("Uso Incorrecto", "Uso: `pan!antiraid limits <uniones>`")
			return err
		}
		uniones, _ := strconv.Atoi(ctx.Args[1])
		antiRaid.Amount = uniones

		if uniones > 0 {
			response = fmt.Sprintf("📊 Se detectará un raid si **%d usuarios** se unen de forma masiva.", uniones)
		} else {
			response = "📊 Detección automática de raid **desactivada**."
		}

	case "action":
		if len(ctx.Args) < 2 {
			_, err = ctx.ReplyError("Uso Incorrecto", "Uso: `pan!antiraid action <kick|ban>`")
			return err
		}
		tipo := strings.ToLower(ctx.Args[1])
		if tipo != "kick" && tipo != "ban" {
			_, err = ctx.ReplyError("Uso Incorrecto", "Acción inválida. Usa `kick` o `ban`.")
			return err
		}
		antiRaid.Action = tipo
		response = fmt.Sprintf("⚡ Acción configurada: Los asaltantes recibirán **%s**.", tipo)

	default:
		_, err = ctx.ReplyError("Uso Incorrecto", "Subcomando no reconocido. Usa `toggle`, `age`, `limits`, o `action`.")
		return err
	}

	_, err = database.GlobalGuildDM.Set(bson.M{"id": ctx.Message.GuildID}, guildData)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ No se pudo guardar la configuración.")
		return err
	}

	_, err = ctx.ReplySuccess("Configuración Anti-Raid", response)
	return err
}
