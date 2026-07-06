package help

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("cmds", "Comando cmds", "pan!cmds", "Help", cmdsCommand)
	messagecommands.RegisterCommand("help", "Comando help", "pan!help", "Help", cmdsCommand)
}
