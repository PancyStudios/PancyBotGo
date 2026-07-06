package mod

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func removeRoleCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionManageRoles) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso para gestionar roles.")
		return err
	}

	if len(ctx.Args) < 2 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el usuario y el rol.\nUso: `pan!removerole @usuario @rol`")
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

	err := ctx.Session.GuildMemberRoleRemove(ctx.Message.GuildID, userID, roleID)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error al remover el rol (¿El rol del bot es más bajo que el rol que intentas remover?): %v", err))
		return err
	}

	_, err = ctx.ReplySuccess("Rol Removido", fmt.Sprintf("✅ Rol <@&%s> removido correctamente a <@%s>.", roleID, userID))
	return err
}
