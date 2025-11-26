// filepath: /home/turbis/GolandProjects/PancyBotGo/internal/commands/mod/removewarn.go
package mod

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

// createRemoveWarnCommand creates the /mod removewarn subcommand
func createRemoveWarnCommand() *discord.Command {
	return discord.NewCommand(
		"removewarn",
		"Elimina una advertencia específica de un usuario",
		"mod",
		removeWarnHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "Usuario del cual eliminar la advertencia",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "id",
			Description:  "ID de la advertencia a eliminar",
			Required:     true,
			Autocomplete: true,
		},
	).WithUserPermissions(discordgo.PermissionModerateMembers).WithAutoComplete(removeWarnAutoComplete).RequiresDatabase()
}

// removeWarnHandler handles the /mod removewarn command
func removeWarnHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()

		// 1. Obtener argumentos
		targetUser := ctx.GetUserOption("usuario")
		warnID := ctx.GetStringOption("id")

		if targetUser == nil {
			ctx.ReplyEphemeral("❌ Debes especificar un usuario válido.")
			return
		}

		if warnID == "" {
			ctx.ReplyEphemeral("❌ Debes especificar el ID de la advertencia.")
			return
		}

		// 2. Feedback inicial
		embedProcess := &discordgo.MessageEmbed{
			Title:       "🗑️ Eliminando advertencia...",
			Description: fmt.Sprintf("Eliminando advertencia de **%s**...\n\nEspere un momento...", targetUser.String()),
			Color:       0xFFFF00, // Yellow
			Footer: &discordgo.MessageEmbedFooter{
				Text:    fmt.Sprintf("Solicitado por %s", ctx.User().String()),
				IconURL: ctx.User().AvatarURL(""),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		if err := ctx.ReplyEmbed(embedProcess); err != nil {
			logger.Error(fmt.Sprintf("Error enviando reply inicial: %v", err), "CMD-RemoveWarn")
			return
		}

		// 3. Consulta DB
		dm := database.GlobalWarnDM
		query := bson.M{"guildId": ctx.Interaction.GuildID, "userId": targetUser.ID}

		doc, err := dm.Get(query)
		if err != nil {
			logger.Error(fmt.Sprintf("Error DB RemoveWarn: %v", err), "CMD-RemoveWarn")
			ctx.EditReply("❌ Error al consultar la base de datos.")
			return
		}

		if doc == nil || len(doc.Warns) == 0 {
			ctx.EditReply("❌ El usuario no tiene advertencias.")
			return
		}

		// 4. Encontrar y eliminar la advertencia
		found := false
		var updatedWarns []models.Warn
		var removedWarn models.Warn

		for _, warn := range doc.Warns {
			if warn.ID == warnID {
				removedWarn = warn
				found = true
			} else {
				updatedWarns = append(updatedWarns, warn)
			}
		}

		if !found {
			ctx.EditReply("❌ No se encontró una advertencia con ese ID.")
			return
		}

		// 5. Actualizar DB
		doc.Warns = updatedWarns
		_, err = dm.Set(query, doc)
		if err != nil {
			logger.Error(fmt.Sprintf("Error guardando RemoveWarn: %v", err), "CMD-RemoveWarn")
			embedError := &discordgo.MessageEmbed{
				Title:       "❌ Error al eliminar advertencia",
				Description: fmt.Sprintf("No se pudo eliminar la advertencia.\nError: `%v`", err),
				Color:       0xFF0000,
			}
			ctx.EditReplyEmbed(embedError)
			return
		}

		// 6. Embed de Éxito
		embedSuccess := &discordgo.MessageEmbed{
			Title:       "✅ Advertencia eliminada con éxito",
			Description: fmt.Sprintf("La advertencia de **%s** ha sido eliminada.\n\n**Razón original:** %s\n**ID:** `%s`", targetUser.String(), removedWarn.Reason, warnID),
			Color:       0x00FF00, // Green
			Footer: &discordgo.MessageEmbedFooter{
				Text:    fmt.Sprintf("Solicitado por %s", ctx.User().String()),
				IconURL: ctx.User().AvatarURL(""),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		ctx.EditReplyEmbed(embedSuccess)

		// 7. Enviar MD al usuario
		embedDM := &discordgo.MessageEmbed{
			Title: "ℹ - Advertencia eliminada",
			Color: 0x00FF00,
			Description: fmt.Sprintf(
				"⚒ - **Servidor:** %s (%s)\n"+
					"🗑 ️ - **Advertencia eliminada:** %s\n\n"+
					"🕒 - **Fecha:** <t:%d:F>",
				ctx.Guild().Name, ctx.Interaction.GuildID, removedWarn.Reason, time.Now().Unix(),
			),
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "💫 - Developed by PancyStudios",
				IconURL: ctx.Client.Session.State.User.AvatarURL(""),
			},
		}

		userChannel, err := ctx.Session.UserChannelCreate(targetUser.ID)
		if err == nil {
			_, _ = ctx.Session.ChannelMessageSendEmbed(userChannel.ID, embedDM)
		} else {
			// Notificar fallo de MD
			msg, _ := ctx.Session.ChannelMessageSend(ctx.Interaction.ChannelID, fmt.Sprintf("ℹ️ No se pudo enviar un mensaje directo a **%s**.", targetUser.String()))
			go func() {
				time.Sleep(5 * time.Second)
				err := ctx.Session.ChannelMessageDelete(ctx.Interaction.ChannelID, msg.ID)
				if err != nil {
					return
				}
			}()
		}
	}()

	return nil
}

// removeWarnAutoComplete handles autocomplete for the removewarn command
func removeWarnAutoComplete(ctx *discord.CommandContext) {
	go func() {
		defer errors.RecoverMiddleware()()

		targetUser := ctx.GetUserOption("usuario")
		if targetUser == nil {
			return
		}

		// Consulta DB
		dm := database.GlobalWarnDM
		query := bson.M{"guildId": ctx.Interaction.GuildID, "userId": targetUser.ID}

		doc, err := dm.Get(query)
		if err != nil || doc == nil || len(doc.Warns) == 0 {
			return
		}

		choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, 25)
		for i, warn := range doc.Warns {
			if i >= 25 {
				break
			}
			name := fmt.Sprintf("ID: %s - Razón: %s", warn.ID, warn.Reason)
			if len(name) > 100 {
				name = name[:97] + "..."
			}
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  name,
				Value: warn.ID,
			})
		}

		ctx.SendAutoCompleteChoices(choices)
	}()
}
