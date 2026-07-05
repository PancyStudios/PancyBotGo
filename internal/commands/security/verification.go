package security

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

// createVerificationCommand creates the /security verification command
func createVerificationCommand() *discord.Command {
	return discord.NewCommand(
		"verification",
		"Gestiona el sistema de verificación del servidor",
		"security",
		verificationHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "panel",
			Description: "Envía el panel interactivo de verificación al canal configurado",
		},
	)
}

// verificationHandler handles the /security verification command
func verificationHandler(ctx *discord.CommandContext) error {
	// Require Administrator permissions
	if ctx.Interaction.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		return ctx.ReplyEmbed(discord.NewEmbed().
			SetColor(0xFF0000).
			SetTitle("❌ Acceso Denegado").
			SetDescription("Necesitas permisos de Administrador para usar este comando.").
			Build())
	}

	// Fetch guild data
	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"_id": ctx.Interaction.GuildID})
	if err != nil || guildDoc == nil {
		return ctx.ReplyEmbed(discord.NewEmbed().
			SetColor(0xFF0000).
			SetDescription("❌ Error al cargar la configuración del servidor.").
			Build())
	}

	// Check if verification is enabled and configured
	if !guildDoc.Protection.Verification.Enable {
		return ctx.ReplyEmbed(discord.NewEmbed().
			SetColor(0xFF0000).
			SetTitle("❌ Sistema Inactivo").
			SetDescription("El sistema de verificación está desactivado. Habilítalo en el **Dashboard**.").
			Build())
	}

	if guildDoc.Protection.Verification.Channel == "" || guildDoc.Protection.Verification.Role == "" {
		return ctx.ReplyEmbed(discord.NewEmbed().
			SetColor(0xFF0000).
			SetTitle("❌ Configuración Incompleta").
			SetDescription("Asegúrate de configurar un Canal y un Rol de verificación en el **Dashboard**.").
			Build())
	}

	channelID := guildDoc.Protection.Verification.Channel

	// Create embed
	embed := discord.NewEmbed().
		SetTitle("🔐 Verificación Requerida").
		SetDescription("Para acceder al resto del servidor y canales, debes verificarte.\n\nHaz clic en el botón de abajo para confirmar que no eres un bot y aceptar las reglas del servidor.").
		SetColor(0x2ecc71).
		Build()

	// Create Action Row with button
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "✅ Verificarme",
					Style:    discordgo.SuccessButton,
					CustomID: "btn_verify_user",
				},
			},
		},
	}

	// Send message to the configured channel
	_, err = ctx.Session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})

	if err != nil {
		return ctx.ReplyEmbed(discord.NewEmbed().
			SetColor(0xFF0000).
			SetDescription(fmt.Sprintf("❌ Error al enviar el panel al canal <#%s>. Revisa mis permisos.", channelID)).
			Build())
	}

	return ctx.ReplyEmbed(discord.NewEmbed().
		SetColor(0x00FF00).
		SetDescription(fmt.Sprintf("✅ Panel de verificación enviado correctamente a <#%s>.", channelID)).
		Build())
}
