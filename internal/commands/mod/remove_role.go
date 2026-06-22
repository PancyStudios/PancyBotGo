package mod

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createRemoveRoleCommand() *discord.Command {
	return discord.NewCommand(
		"removerole",
		"✨ | Remueve un rol a un usuario",
		"mod",
		removeRoleHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "Usuario al que se le removerá el rol",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "rol",
			Description: "Rol que se removerá",
			Required:    true,
		},
	).WithUserPermissions(discordgo.PermissionManageRoles).
		WithBotPermissions(discordgo.PermissionManageRoles)
}

func removeRoleHandler(ctx *discord.CommandContext) error {
	user := ctx.GetUserOption("usuario")
	if user == nil {
		return ctx.ReplyEphemeral("❌ Debes especificar un usuario.")
	}

	roleID := ctx.GetOption("rol").Value.(string)

	err := ctx.Session.GuildMemberRoleRemove(ctx.Interaction.GuildID, user.ID, roleID)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error al remover el rol (¿El rol del bot es más bajo que el rol que intentas remover?): %v", err))
	}

	return ctx.Reply(fmt.Sprintf("✅ Rol removido correctamente a **%s**.", user.Username))
}
