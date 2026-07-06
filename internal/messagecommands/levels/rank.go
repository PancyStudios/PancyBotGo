package levels

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

func rankCommand(ctx *messagecommands.MessageContext) error {
	guildID := ctx.Message.GuildID

	guildData, err := database.GlobalGuildDM.Get(bson.M{"id": guildID})
	if err != nil || guildData == nil || !guildData.Levels.Enable {
		_, err = ctx.ReplyError("Error", "❌ El sistema de niveles está desactivado en este servidor.")
		return err
	}

	targetUser := ctx.Message.Author
	if len(ctx.Args) > 0 {
		parsedUserID := ctx.ParseUser(0)
		if parsedUserID != "" {
			member, err := ctx.Session.GuildMember(guildID, parsedUserID)
			if err == nil {
				targetUser = member.User
			}
		}
	}

	if targetUser.Bot {
		_, err = ctx.ReplyError("Error", "🤖 Los bots no tienen niveles.")
		return err
	}

	profile, err := database.GetLocalLevelProfile(guildID, targetUser.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error al obtener perfil de nivel: %v", err), "RankCommand")
		_, err = ctx.ReplyError("Error", "❌ Ocurrió un error al obtener el rango.")
		return err
	}

	nextLevel := profile.Level + 1
	requiredXP := nextLevel * nextLevel * 100

	currentLevelXP := int64(0)
	if profile.Level > 0 {
		currentLevelXP = profile.Level * profile.Level * 100
	}

	xpInCurrentLevel := profile.XP - currentLevelXP
	xpNeededForNext := requiredXP - currentLevelXP

	progressPercent := float64(xpInCurrentLevel) / float64(xpNeededForNext) * 100

	progressBar := createProgressBar(progressPercent, 10)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🌟 Rango de %s", targetUser.Username),
		Description: fmt.Sprintf("¡Sigue chateando para subir de nivel!\n\n**Nivel Actual:** %d\n**Experiencia:** %d / %d XP\n\n%s (%.1f%%)", profile.Level, profile.XP, requiredXP, progressBar, progressPercent),
		Color:       0x00FFFF,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: targetUser.AvatarURL("128"),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Mensajes totales: %d", profile.TotalMessages),
		},
	}

	_, err = ctx.ReplyEmbed(embed)
	return err
}

func createProgressBar(percent float64, length int) string {
	filledBlocks := int((percent / 100.0) * float64(length))
	if filledBlocks > length {
		filledBlocks = length
	}
	if filledBlocks < 0 {
		filledBlocks = 0
	}

	emptyBlocks := length - filledBlocks

	bar := ""
	for i := 0; i < filledBlocks; i++ {
		bar += "🟦"
	}
	for i := 0; i < emptyBlocks; i++ {
		bar += "⬛"
	}
	return bar
}
