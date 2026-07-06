package levels

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("leaderboard", "Comando leaderboard", "pan!leaderboard", "Levels", leaderboardCommand)
	messagecommands.RegisterCommand("rank", "Comando rank", "pan!rank", "Levels", rankCommand)
	messagecommands.RegisterCommand("rewards", "Comando rewards", "pan!rewards", "Levels", rewardsCommand)
	messagecommands.RegisterCommand("togglelevels", "Comando togglelevels", "pan!togglelevels", "Levels", toggleCommand)
}
