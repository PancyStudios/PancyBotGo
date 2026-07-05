package levels

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

var leaderboardCommand = &discord.Command{
	Name:        "leaderboard",
	Description: "🏆 | Muestra los usuarios con más nivel en el servidor",
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

		topProfiles, err := database.GetTopLevels(guildID, 10)
		if err != nil {
			logger.Error(fmt.Sprintf("Error obteniendo leaderboard para %s: %v", guildID, err), "Leaderboard")
			return ctx.ReplyEphemeral("❌ Ocurrió un error al obtener la clasificación.")
		}

		if len(topProfiles) == 0 {
			return ctx.ReplyEphemeral("📉 Aún no hay usuarios con experiencia en este servidor.")
		}

		description := "¡Estos son los usuarios más activos del servidor!\n\n"

		for i, profile := range topProfiles {
			medal := "🏅"
			switch i {
			case 0:
				medal = "🥇"
			case 1:
				medal = "🥈"
			case 2:
				medal = "🥉"
			}
			
			description += fmt.Sprintf("%s **#%d** <@%s> - Nivel %d (%d XP)\n", medal, i+1, profile.UserID, profile.Level, profile.XP)
		}

		embed := &discordgo.MessageEmbed{
			Title:       "🏆 Tabla de Clasificación de Niveles",
			Description: description,
			Color:       0xFFD700, // Gold
		}

		return ctx.ReplyEmbed(embed)
	},
}
