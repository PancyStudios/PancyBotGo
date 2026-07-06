package levels

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/bwmarrin/discordgo"
)

func toggleCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionAdministrator) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso de Administrador.")
		return err
	}

	guildID := ctx.Message.GuildID
	guildData, err := database.GlobalGuildDM.Get(bson.M{"id": guildID})
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error al obtener la configuración del servidor: %v", err))
		return err
	}

	newState := !guildData.Levels.Enable
	guildData.Levels.Enable = newState

	_, err = database.GlobalGuildDM.Set(bson.M{"id": guildID}, guildData)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error al guardar la configuración: %v", err))
		return err
	}

	if newState {
		_, err = ctx.ReplySuccess("Niveles Activados", "✅ **Sistema de Niveles Activado.**\n\nLos usuarios ahora ganarán experiencia al chatear.")
		return err
	}

	_, err = ctx.ReplySuccess("Niveles Desactivados", "✅ **Sistema de Niveles Desactivado.**\n\nLos usuarios ya no ganarán experiencia al chatear. La experiencia acumulada no se perderá.")
	return err
}
