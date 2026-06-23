package mod

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/scheduler"
	"github.com/bwmarrin/discordgo"
)

func createTempBanCommand() *discord.Command {
	return discord.NewCommand(
		"tempban",
		"⏳ | Banea temporalmente a alguien",
		"mod",
		tempBanHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "🛡️ | Usuario a banear",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "duracion_horas",
			Description: "🛡️ | Duración del ban en horas",
			Required:    true,
			MinValue:    func() *float64 { v := 1.0; return &v }(),
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "razon",
			Description: "🛡️ | Razón del ban",
			Required:    false,
		},
	).WithUserPermissions(discordgo.PermissionBanMembers).
		WithBotPermissions(discordgo.PermissionBanMembers)
}

func tempBanHandler(ctx *discord.CommandContext) error {
	user := ctx.GetUserOption("usuario")
	if user == nil {
		return ctx.ReplyEphemeral("❌ Debes especificar un usuario.")
	}

	horas := ctx.GetIntOption("duracion_horas")
	duracion := time.Duration(horas) * time.Hour

	reason := ctx.GetStringOption("razon")
	if reason == "" {
		reason = "Sin razón especificada"
	}

	// Ban user with 0 days of message deletion (can be adjusted)
	err := ctx.Session.GuildBanCreateWithReason(
		ctx.Interaction.GuildID,
		user.ID,
		reason,
		0,
	)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error al banear: %v", err))
	}

	// Register in scheduler
	err = scheduler.AddTempBan(ctx.Interaction.GuildID, user.ID, duracion)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("⚠️ El usuario fue baneado, pero hubo un error al programar su desbaneo: %v", err))
	}

	return ctx.Reply(fmt.Sprintf("🔨 **%s** ha sido baneado temporalmente por %d horas.\n**Razón:** %s", user.Username, horas, reason))
}
