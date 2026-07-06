package security

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func verificationCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionAdministrator) {
		_, err := ctx.ReplyError("Acceso Denegado", "Necesitas permisos de Administrador para usar este comando.")
		return err
	}

	if len(ctx.Args) == 0 || strings.ToLower(ctx.Args[0]) != "panel" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!verification panel`")
		return err
	}

	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Message.GuildID})
	if err != nil || guildDoc == nil {
		_, err = ctx.ReplyError("Error", "❌ Error al cargar la configuración del servidor.")
		return err
	}

	if !guildDoc.Protection.Verification.Enable {
		_, err = ctx.ReplyError("Sistema Inactivo", "El sistema de verificación está desactivado. Habilítalo en el **Dashboard**.")
		return err
	}

	if guildDoc.Protection.Verification.Channel == "" || guildDoc.Protection.Verification.Role == "" {
		_, err = ctx.ReplyError("Configuración Incompleta", "Asegúrate de configurar un Canal y un Rol de verificación en el **Dashboard**.")
		return err
	}

	channelID := guildDoc.Protection.Verification.Channel

	embed := &discordgo.MessageEmbed{
		Title:       "🔐 Verificación Requerida",
		Description: "Para acceder al resto del servidor y canales, debes verificarte.\n\nHaz clic en el botón de abajo para confirmar que no eres un bot y aceptar las reglas del servidor.",
		Color:       0x2ecc71,
	}

	var verifyButton discordgo.MessageComponent
	if guildDoc.Protection.Verification.Type == "web" {
		verifyButton = discordgo.Button{
			Label: "🌐 Verificar en la Web",
			Style: discordgo.LinkButton,
			URL:   fmt.Sprintf("https://pancybot.miau.media/verify/%s", guildDoc.ID),
		}
	} else {
		verifyButton = discordgo.Button{
			Label:    "✅ Verificarme",
			Style:    discordgo.SuccessButton,
			CustomID: "btn_verify_user",
		}
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{verifyButton},
		},
	}

	_, err = ctx.Session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})

	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ No se pudo enviar el panel al canal configurado.")
		return err
	}

	_, err = ctx.ReplySuccess("Panel Enviado", fmt.Sprintf("✅ Panel de verificación enviado correctamente a <#%s>.", channelID))
	return err
}
