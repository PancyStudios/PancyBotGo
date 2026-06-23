package levels

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// RegisterCommands registers the levels commands
func RegisterCommands(client *discord.ExtendedClient) {
	levelsGroup := client.CommandHandler.BuildCommandGroup(
		"levels",
		"🌟 | Sistema de niveles y experiencia",
		rankCommand,
		leaderboardCommand,
	)

	client.CommandHandler.AddGlobalCommand(levelsGroup)
}
