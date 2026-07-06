package help

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("cmds", cmdsCommand)
	messagecommands.RegisterCommand("help", cmdsCommand)
}
