package levels

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

var rankCommand = &discord.Command{
	Name:        "rank",
	Description: "🏅 | Muestra tu nivel y experiencia actual",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "🌟 | Usuario del que quieres ver el rango (opcional)",
			Required:    false,
		},
	},
	Run: func(ctx *discord.CommandContext) error {
		guildID := ctx.Interaction.GuildID
		if guildID == "" {
			return ctx.ReplyEphemeral("❌ Este comando solo se puede usar en un servidor.")
		}

		// Verificar si el sistema está activado
		guildData, err := database.GlobalGuildDM.Get(bson.M{"id": guildID})
		if err != nil || guildData == nil || !guildData.Levels.Enable {
			return ctx.ReplyEphemeral("❌ El sistema de niveles está desactivado en este servidor.")
		}

		targetUser := ctx.User()
		userOpt := ctx.GetOption("usuario")
		if userOpt != nil && userOpt.UserValue(ctx.Session) != nil {
			targetUser = userOpt.UserValue(ctx.Session)
		}

		if targetUser.Bot {
			return ctx.ReplyEphemeral("🤖 Los bots no tienen niveles.")
		}

		profile, err := database.GetLocalLevelProfile(guildID, targetUser.ID)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al obtener perfil de nivel: %v", err), "RankCommand")
			return ctx.ReplyEphemeral("❌ Ocurrió un error al obtener el rango.")
		}

		// Calcular XP necesaria para el próximo nivel
		nextLevel := profile.Level + 1
		requiredXP := nextLevel * nextLevel * 100
		
		// Calcular XP base del nivel actual
		currentLevelXP := int64(0)
		if profile.Level > 0 {
			currentLevelXP = profile.Level * profile.Level * 100
		}

		// Progreso dentro del nivel actual
		xpInCurrentLevel := profile.XP - currentLevelXP
		xpNeededForNext := requiredXP - currentLevelXP

		progressPercent := float64(xpInCurrentLevel) / float64(xpNeededForNext) * 100
		
		progressBar := createProgressBar(progressPercent, 10)

		embed := discord.NewEmbed().
			SetTitle(fmt.Sprintf("🌟 Rango de %s", targetUser.Username)).
			SetDescription(fmt.Sprintf("¡Sigue chateando para subir de nivel!\n\n**Nivel Actual:** %d\n**Experiencia:** %d / %d XP\n\n%s (%.1f%%)", profile.Level, profile.XP, requiredXP, progressBar, progressPercent)).
			SetColor(0x00FFFF). // Cyan
			SetThumbnail(targetUser.AvatarURL("128")).
			SetFooter(fmt.Sprintf("Mensajes totales: %d", profile.TotalMessages), "").
			Build()

		return ctx.ReplyEmbed(embed)
	},
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
