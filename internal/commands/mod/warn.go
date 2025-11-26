// Package mod - /mod warn command
package mod

import (
	"fmt"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// createWarnCommand creates the /mod warn subcommand
func createWarnCommand() *discord.Command {
	return discord.NewCommand(
		"warn",
		"Advierte a un usuario",
		"mod",
		warnHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "Usuario a advertir",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "razon",
			Description: "Razón de la advertencia",
			Required:    false,
		},
	).WithUserPermissions(discordgo.PermissionModerateMembers).RequiresDatabase()
}

// warnHandler handles the /mod warn command
func warnHandler(ctx *discord.CommandContext) error {
	go func() {
		// Middleware de recuperación por si algo crashea dentro de la goroutine
		defer errors.RecoverMiddleware()()

		// 1. Obtener argumentos
		targetUser := ctx.GetUserOption("usuario")
		reason := ctx.GetStringOption("razon")
		if reason == "" {
			reason = "Razón no proporcionada"
		}

		// 2. Validaciones básicas de usuario
		if targetUser == nil {
			ctx.ReplyEphemeral("❌ Debes especificar un usuario válido.")
			return
		}
		if targetUser.ID == ctx.User().ID {
			ctx.ReplyEphemeral("❌ No puedes advertirte a ti misma.")
			return
		}
		if targetUser.Bot {
			ctx.ReplyEphemeral("❌ No puedes advertir a un bot.")
			return
		}

		// Obtener Guild
		guild := ctx.Guild()
		if guild == nil {
			ctx.ReplyEphemeral("❌ Error obteniendo información del servidor.")
			return
		}

		if targetUser.ID == guild.OwnerID {
			ctx.ReplyEphemeral("❌ No puedes advertir al dueño del servidor.")
			return
		}

		// 3. Validación de Jerarquía de Roles
		targetMember, err := ctx.Session.GuildMember(ctx.Interaction.GuildID, targetUser.ID)
		if err != nil {
			ctx.ReplyEphemeral("❌ No se pudo obtener la información del miembro objetivo.")
			return
		}

		executorPosition := getHighestRolePosition(guild, ctx.Member())
		targetPosition := getHighestRolePosition(guild, targetMember)

		// Si el objetivo tiene un rol igual o superior y NO eres el dueño
		if targetPosition >= executorPosition && ctx.User().ID != guild.OwnerID {
			ctx.ReplyEphemeral("❌ No puedes advertir a un usuario con un rol mayor o igual al tuyo.")
			return
		}

		// 4. Feedback inicial (Público)
		embedProcess := &discordgo.MessageEmbed{
			Title:       "⚠️ Advertencia en proceso...",
			Description: fmt.Sprintf("Advirtiendo a **%s**...\n\nEspere un momento...", targetUser.String()),
			Color:       0xFFFF00, // Yellow
			Footer: &discordgo.MessageEmbedFooter{
				Text:    fmt.Sprintf("Solicitado por %s", ctx.User().String()),
				IconURL: ctx.User().AvatarURL(""),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		if err := ctx.ReplyEmbed(embedProcess); err != nil {
			logger.Error(fmt.Sprintf("Error enviando reply inicial: %v", err), "CMD-Warn")
			return
		}

		// 5. Preparar datos para DB
		warnID := uuid.New().String()
		shortID := strings.ReplaceAll(warnID, "-", "")[:6]

		newWarn := models.Warn{
			Reason:    reason,
			Moderator: ctx.User().ID,
			ID:        shortID,
			Timestamp: time.Now().Unix(),
		}

		// 6. Operación DB
		dm := database.NewDataManager[models.WarnsDocument]("Warns", database.Get())

		query := bson.M{"guildId": ctx.Interaction.GuildID, "userId": targetUser.ID}
		doc, err := dm.Get(query)

		if err != nil {
			logger.Error(fmt.Sprintf("Error DB Warn: %v", err), "CMD-Warn")
			ctx.EditReply("❌ Error al acceder a la base de datos.")
			return
		}

		if doc == nil {
			newDoc := models.WarnsDocument{
				GuildID: ctx.Interaction.GuildID,
				UserID:  targetUser.ID,
				Warns:   []models.Warn{newWarn},
			}
			_, err = dm.Set(query, newDoc)
		} else {
			doc.Warns = append(doc.Warns, newWarn)
			_, err = dm.Set(query, doc)
		}

		if err != nil {
			logger.Error(fmt.Sprintf("Error guardando Warn: %v", err), "CMD-Warn")
			embedError := &discordgo.MessageEmbed{
				Title:       "❌ Error al advertir",
				Description: fmt.Sprintf("No se pudo guardar la advertencia en la base de datos.\nError: `%v`", err),
				Color:       0xFF0000,
			}
			ctx.EditReplyEmbed(embedError)
			return
		}

		// 7. Embed de Éxito
		embedSuccess := &discordgo.MessageEmbed{
			Title:       "✅ Usuario advertido con éxito",
			Description: fmt.Sprintf("El usuario **%s** ha sido advertido correctamente.\n\n**Razón:** %s\n**ID de Advertencia:** `%s`", targetUser.String(), reason, shortID),
			Color:       0x00FF00, // Green
			Footer: &discordgo.MessageEmbedFooter{
				Text:    fmt.Sprintf("Solicitado por %s", ctx.User().String()),
				IconURL: ctx.User().AvatarURL(""),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		ctx.EditReplyEmbed(embedSuccess)

		// 8. Enviar MD al usuario
		embedDM := &discordgo.MessageEmbed{
			Title: "⚠️ - Has recibido una advertencia",
			Color: 0xFFFF00,
			Description: fmt.Sprintf(
				"⚒️ - **Servidor:** %s (%s)\n"+
					"🔨 - **Razón:** %s\n\n"+
					"🕒 - **Fecha:** <t:%d:F>",
				guild.Name, guild.ID, reason, time.Now().Unix(),
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
			msg, _ := ctx.Session.ChannelMessageSend(ctx.Interaction.ChannelID, fmt.Sprintf("⚠️ No se pudo enviar un mensaje directo a **%s**.", targetUser.String()))
			go func() {
				time.Sleep(5 * time.Second)
				err := ctx.Session.ChannelMessageDelete(ctx.Interaction.ChannelID, msg.ID)
				if err != nil {
					return
				}
			}()
		}
	}()

	// Retornamos nil inmediatamente para liberar al CommandHandler principal
	return nil
}

// getHighestRolePosition calcula la posición más alta de los roles de un miembro
func getHighestRolePosition(guild *discordgo.Guild, member *discordgo.Member) int {
	highest := 0
	// Crear un mapa rápido de ID de rol -> Posición para no iterar demasiado
	roleMap := make(map[string]int)
	for _, r := range guild.Roles {
		roleMap[r.ID] = r.Position
	}

	for _, roleID := range member.Roles {
		if pos, ok := roleMap[roleID]; ok {
			if pos > highest {
				highest = pos
			}
		}
	}
	return highest
}
