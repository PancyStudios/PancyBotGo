package levels

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("leaderboard", leaderboardCommand)
	messagecommands.RegisterCommand("rank", rankCommand)
	messagecommands.RegisterCommand("rewards", rewardsCommand)
	messagecommands.RegisterCommand("togglelevels", toggleCommand)
}
