package mod

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createAssignRoleCommand() *discord.Command {
	return discord.NewCommand(
		"assign-role",
		"Asigna un rol a un usuario",
		"mod",
		assignRoleHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "Usuario al que se le asignará el rol",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "rol",
			Description: "Rol que se asignará",
			Required:    true,
		},
	).WithUserPermissions(discordgo.PermissionManageRoles).
		WithBotPermissions(discordgo.PermissionManageRoles)
}

func assignRoleHandler(ctx *discord.CommandContext) error {
	user := ctx.GetUserOption("usuario")
	if user == nil {
		return ctx.ReplyEphemeral("❌ Debes especificar un usuario.")
	}

	roleID := ctx.GetOption("rol").Value.(string)

	err := ctx.Session.GuildMemberRoleAdd(ctx.Interaction.GuildID, user.ID, roleID)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error al asignar el rol (¿El rol del bot es más bajo que el rol que intentas asignar?): %v", err))
	}

	return ctx.Reply(fmt.Sprintf("✅ Rol asignado correctamente a **%s**.", user.Username))
}
