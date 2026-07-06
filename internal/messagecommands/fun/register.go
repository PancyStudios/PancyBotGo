package fun

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register fun text commands
func Register() {
	messagecommands.RegisterCommand("8ball", "Comando 8ball", "pan!8ball", "Fun", eightBallCommand)
	messagecommands.RegisterCommand("ppt", "Comando ppt", "pan!ppt", "Fun", pptCommand)
	messagecommands.RegisterCommand("ascii", "Comando ascii", "pan!ascii", "Fun", asciiCommand)
	messagecommands.RegisterCommand("dog", "Comando dog", "pan!dog", "Fun", dogCommand)
}
