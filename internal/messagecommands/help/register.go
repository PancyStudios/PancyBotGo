package help

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("cmds", "Comando cmds", "pan!cmds", "General", cmdsCommand)
	messagecommands.RegisterCommand("help", "Comando help", "pan!help", "General", cmdsCommand)
}
