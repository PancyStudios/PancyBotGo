package security

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func antibotsCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionManageServer) {
		_, err := ctx.ReplyError("Acceso Denegado", "Necesitas permisos de Administrador para usar este comando.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!antibots <all|only_nv|only_v|disabled>`")
		return err
	}

	option := strings.ToLower(ctx.Args[0])
	validOptions := map[string]bool{"all": true, "only_nv": true, "only_v": true, "disabled": true}

	if !validOptions[option] {
		_, err := ctx.ReplyError("Opción Inválida", "Las opciones válidas son: `all`, `only_nv`, `only_v`, `disabled`.")
		return err
	}

	guildData, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Message.GuildID})
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error al obtener los datos del servidor: %v", err))
		return err
	}

	if option == "disabled" {
		guildData.Protection.Antibots.Enable = false
		guildData.Protection.Antibots.Type = ""
	} else {
		guildData.Protection.Antibots.Enable = true
		guildData.Protection.Antibots.Type = option
	}

	_, err = database.GlobalGuildDM.Set(bson.M{"id": ctx.Message.GuildID}, guildData)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error al guardar en la base de datos: %v", err))
		return err
	}

	status := "activada"
	if option == "disabled" {
		status = "desactivada"
	}

	_, err = ctx.ReplySuccess("Anti-Bots", fmt.Sprintf("🛡️ La protección Anti-Bots ha sido **%s** (Modo: %s).", status, option))
	return err
}
