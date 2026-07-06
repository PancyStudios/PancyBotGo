package mod

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func assignRoleCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionManageRoles) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para gestionar roles.")
		return err
	}

	if len(ctx.Args) < 2 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el usuario y el rol.\nUso: `pan!assign-role @usuario @rol`")
		return err
	}

	userID := ctx.ParseUser(0)
	if userID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un usuario válido.")
		return err
	}

	roleID := ctx.ParseRole(1)
	if roleID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un rol válido.")
		return err
	}

	err := ctx.Session.GuildMemberRoleAdd(ctx.Message.GuildID, userID, roleID)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error al asignar el rol (¿El rol del bot es más bajo que el rol que intentas asignar?): %v", err))
		return err
	}

	_, err = ctx.ReplySuccess("Rol Asignado", fmt.Sprintf("✅ Rol <@&%s> asignado correctamente a <@%s>.", roleID, userID))
	return err
}
