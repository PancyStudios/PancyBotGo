package levels

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("leaderboard", "Comando leaderboard", "pan!leaderboard", "General", leaderboardCommand)
	messagecommands.RegisterCommand("rank", "Comando rank", "pan!rank", "General", rankCommand)
	messagecommands.RegisterCommand("rewards", "Comando rewards", "pan!rewards", "General", rewardsCommand)
	messagecommands.RegisterCommand("togglelevels", "Comando togglelevels", "pan!togglelevels", "General", toggleCommand)
}
