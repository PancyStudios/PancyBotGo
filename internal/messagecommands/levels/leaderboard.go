package levels

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

func leaderboardCommand(ctx *messagecommands.MessageContext) error {
	guildID := ctx.Message.GuildID

	guildData, err := database.GlobalGuildDM.Get(bson.M{"id": guildID})
	if err != nil || guildData == nil || !guildData.Levels.Enable {
		_, err = ctx.ReplyError("Error", "❌ El sistema de niveles está desactivado en este servidor.")
		return err
	}

	topProfiles, err := database.GetTopLevels(guildID, 10)
	if err != nil {
		logger.Error(fmt.Sprintf("Error obteniendo leaderboard para %s: %v", guildID, err), "Leaderboard")
		_, err = ctx.ReplyError("Error", "❌ Ocurrió un error al obtener la clasificación.")
		return err
	}

	if len(topProfiles) == 0 {
		_, err = ctx.ReplySuccess("Clasificación", "📉 Aún no hay usuarios con experiencia en este servidor.")
		return err
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
		Color:       0xFFD700,
	}

	_, err = ctx.ReplyEmbed(embed)
	return err
}
