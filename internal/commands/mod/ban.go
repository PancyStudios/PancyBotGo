// Package mod - /mod ban command
package mod

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

// createBanCommand creates the /mod ban subcommand
func createBanCommand() *discord.Command {
	return discord.NewCommand(
		"ban",
		"🔨 | Banea a un usuario del servidor",
		"mod",
		banHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "🛡️ | Usuario a banear",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "razon",
			Description: "🛡️ | Razón del ban",
			Required:    false,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "dias",
			Description: "🛡️ | Días de mensajes a eliminar (0-7)",
			Required:    false,
			MinValue:    func() *float64 { v := 0.0; return &v }(),
			MaxValue:    7,
		},
	).WithUserPermissions(discordgo.PermissionBanMembers).
		WithBotPermissions(discordgo.PermissionBanMembers)
}

// banHandler handles the /mod ban command
func banHandler(ctx *discord.CommandContext) error {
	user := ctx.GetUserOption("usuario")
	if user == nil {
		return ctx.ReplyEphemeral("❌ Debes especificar un usuario.")
	}

	reason := ctx.GetStringOption("razon")
	if reason == "" {
		reason = "Sin razón especificada"
	}

	days := int(ctx.GetIntOption("dias"))

	// Perform the ban
	err := ctx.Session.GuildBanCreateWithReason(
		ctx.Interaction.GuildID,
		user.ID,
		reason,
		days,
	)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error al banear: %v", err))
	}

	return ctx.Reply(fmt.Sprintf("🔨 **%s** ha sido baneado.\n**Razón:** %s", user.Username, reason))
}
