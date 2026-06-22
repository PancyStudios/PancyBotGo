package mod

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createSoftbanCommand() *discord.Command {
	return discord.NewCommand(
		"softban",
		"💨 | Banea temporalmente a un usuario para borrar sus mensajes y lo desbanea inmediatamente",
		"mod",
		softbanHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "Usuario a hacer softban",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "razon",
			Description: "Razón del softban",
			Required:    false,
		},
	).WithUserPermissions(discordgo.PermissionBanMembers).
		WithBotPermissions(discordgo.PermissionBanMembers)
}

func softbanHandler(ctx *discord.CommandContext) error {
	user := ctx.GetUserOption("usuario")
	if user == nil {
		return ctx.ReplyEphemeral("❌ Debes especificar un usuario.")
	}

	// Prevent banning self
	if user.ID == ctx.User().ID {
		return ctx.ReplyEphemeral("❌ No puedes banearte a ti mismo.")
	}

	// Prevent banning owner
	guild, err := ctx.Session.State.Guild(ctx.Interaction.GuildID)
	if err == nil && user.ID == guild.OwnerID {
		return ctx.ReplyEphemeral("❌ No puedes banear al dueño del servidor.")
	}

	reason := ctx.GetStringOption("razon")
	if reason == "" {
		reason = "Sin razón especificada"
	}

	// Ban user with 7 days of message deletion
	err = ctx.Session.GuildBanCreateWithReason(
		ctx.Interaction.GuildID,
		user.ID,
		reason,
		7,
	)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error al banear: %v", err))
	}

	// Immediately unban the user
	err = ctx.Session.GuildBanDelete(ctx.Interaction.GuildID, user.ID)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("⚠️ El usuario fue baneado pero ocurrió un error al desbanearlo: %v", err))
	}

	return ctx.Reply(fmt.Sprintf("♻️ **%s** ha sido softbaneado.\n**Razón:** %s", user.Username, reason))
}
