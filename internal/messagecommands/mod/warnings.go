package mod

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func warningsCommand(ctx *messagecommands.MessageContext) error {
	targetUser := ctx.Message.Author
	isSelf := true

	if len(ctx.Args) > 0 {
		parsedUserID := ctx.ParseUser(0)
		if parsedUserID != "" {
			member, err := ctx.Session.GuildMember(ctx.Message.GuildID, parsedUserID)
			if err == nil {
				targetUser = member.User
				isSelf = (targetUser.ID == ctx.Message.Author.ID)
			}
		}
	}

	if !isSelf && !ctx.HasPermission(discordgo.PermissionManageMessages) {
		_, err := ctx.ReplyError("Acceso Denegado", "❌ No tienes permisos para ver la lista de advertencias de otro usuario.")
		return err
	}

	dm := database.GlobalWarnDM
	query := bson.M{"guildId": ctx.Message.GuildID, "userId": targetUser.ID}

	doc, err := dm.Get(query)
	if err != nil {
		logger.Error(fmt.Sprintf("Error DB Warnings: %v", err), "CMD-Warnings")
		_, err = ctx.ReplyError("Error", "❌ Error al consultar la base de datos.")
		return err
	}

	embedClear := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🔖 - Lista de advertencias de %s", targetUser.Username),
		Color:       0x00FF00,
		Description: fmt.Sprintf("No se han encontrado advertencias del usuario en este servidor\n\n> 💫 - **Cantidad de advertencias:** 0\n> 🕒 - **Fecha de consulta:** <t:%d>", time.Now().Unix()),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "💫 - Developed by PancyStudios",
		},
	}

	if doc == nil || len(doc.Warns) == 0 {
		_, err = ctx.ReplyEmbed(embedClear)
		return err
	}

	description := fmt.Sprintf("Se han encontrado advertencias del usuario en este servidor:\n\n> 💫 - **Cantidad de advertencias:** %d\n> 🕒 - **Fecha de consulta:** <t:%d>\n\n", len(doc.Warns), time.Now().Unix())

	for _, warn := range doc.Warns {
		description += fmt.Sprintf("> 📝 - **ID de la advertencia:** `%s`\n> 🔨 - **Moderador a cargo:** <@%s>\n> 📅 - **Fecha:** <t:%d>\n> 📄 - **Razón:** %s\n\n", warn.ID, warn.Moderator, warn.Timestamp, warn.Reason)
	}

	if len(description) > 4000 {
		description = description[:4000] + "...\n\n**La lista ha sido truncada por ser demasiado larga.**"
	}

	embedSuccess := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🔖 - Lista de advertencias de %s", targetUser.Username),
		Color:       0xFF0000,
		Description: description,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "💫 - Developed by PancyStudios",
		},
	}

	_, err = ctx.ReplyEmbed(embedSuccess)
	return err
}
