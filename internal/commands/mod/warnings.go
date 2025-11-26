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

// createWarningsCommand creates the /mod warns subcommand
func createWarningsCommand() *discord.Command {
	return discord.NewCommand(
		"warns", // Nombre del comando igual que en TS
		"Lista de advertencias de un usuario",
		"mod",
		warningsHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "[STAFF] Usuario a buscar (opcional)",
			Required:    false,
		},
	).RequiresDatabase()
}

func warningsHandler(ctx *discord.CommandContext) error {
	// Goroutine para no bloquear el hilo principal
	go func() {
		defer errors.RecoverMiddleware()()

		// 1. Determinar objetivo y permisos
		targetUser := ctx.GetUserOption("usuario")
		isSelf := false

		// Verificar permisos del ejecutor (ManageMessages es el equivalente común a Moderador)
		perms, err := ctx.Session.UserChannelPermissions(ctx.User().ID, ctx.Interaction.ChannelID)
		if err != nil {
			perms = 0
		}
		isModerator := (perms & discordgo.PermissionManageMessages) != 0

		if targetUser == nil {
			targetUser = ctx.User()
			isSelf = true
		}

		// Si intenta ver advertencias de otro y no es moderador
		if !isSelf && !isModerator {
			ctx.ReplyEphemeral("❌ No tienes permisos para ver la lista de advertencias de otro usuario.")
			return
		}

		// 2. Feedback inicial (Efímero como en TS)
		embedLoading := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("🔖 - Lista de advertencias de %s", targetUser.Username),
			Description: "Espere un momento mientras obtenemos las advertencias...\n\n> 💫 - **Cantidad de advertencias:** Desconocido\n> 🕒 - **Fecha de consulta:** Cargando...",
			Color:       0x3498db, // Blue
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "💫 - Developed by PancyStudios",
				IconURL: ctx.Guild().IconURL(""),
			},
		}

		if err := ReplyEphemeralEmbed(ctx, embedLoading); err != nil {
			logger.Error(fmt.Sprintf("Error enviando reply inicial warnings: %v", err), "CMD-Warnings")
			return
		}

		// 3. Consulta DB
		dm := database.NewDataManager[models.WarnsDocument]("Warns", database.Get())
		query := bson.M{"guildId": ctx.Interaction.GuildID, "userId": targetUser.ID}

		doc, err := dm.Get(query)

		// Embed base para "Sin advertencias"
		embedClear := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("🔖 - Lista de advertencias de %s", targetUser.Username),
			Color:       0x00FF00, // Green
			Description: fmt.Sprintf("No se han encontrado advertencias del usuario en este servidor\n\n> 💫 - **Cantidad de advertencias:** 0\n> 🕒 - **Fecha de consulta:** <t:%d>", time.Now().Unix()),
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "💫 - Developed by PancyStudios",
				IconURL: ctx.Guild().IconURL(""),
			},
		}

		if err != nil {
			logger.Error(fmt.Sprintf("Error DB Warnings: %v", err), "CMD-Warnings")
			ctx.EditReply("❌ Error al consultar la base de datos.")
			return
		}

		if doc == nil || len(doc.Warns) == 0 {
			ctx.EditReplyEmbed(embedClear)
			return
		}

		// 4. Construir lista de advertencias
		embedList := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("🔖 - Lista de advertencias de %s (%s)", targetUser.Username, targetUser.ID),
			Color: 0xFFA500, // Orange
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "💫 - Developed by PancyStudios",
				IconURL: ctx.Guild().IconURL(""),
			},
		}

		var description string

		// Iterar warnings
		for _, warn := range doc.Warns {
			modName := "Desconocido"

			// Lógica de visualización del moderador (TS logic)
			if isModerator {
				// Intentar obtener nombre del mod
				modUser, err := ctx.Session.User(warn.Moderator)
				if err == nil {
					modName = fmt.Sprintf("%s#%s", modUser.Username, modUser.Discriminator)
				} else {
					modName = warn.Moderator // Fallback al ID si no se encuentra
				}
			} else {
				modName = "Oculto"
			}

			description += fmt.Sprintf("> **Advertencia:** %s \n> **Moderador:** %s \n> **ID:** %s \n\n", warn.Reason, modName, warn.ID)
		}

		description += fmt.Sprintf("> 💫 - **Cantidad de advertencias:** %d \n> 🕒 - **Fecha de consulta:** <t:%d>", len(doc.Warns), time.Now().Unix())

		embedList.Description = description

		// 5. Enviar respuesta final
		ctx.EditReplyEmbed(embedList)
	}()

	return nil
}

func ReplyEphemeralEmbed(ctx *discord.CommandContext, embed *discordgo.MessageEmbed) error {
	return ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}
