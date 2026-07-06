package fun

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register fun text commands
func Register() {
	messagecommands.RegisterCommand("8ball", "Comando 8ball", "pan!8ball", "General", eightBallCommand)
	messagecommands.RegisterCommand("ppt", "Comando ppt", "pan!ppt", "General", pptCommand)
	messagecommands.RegisterCommand("ascii", "Comando ascii", "pan!ascii", "General", asciiCommand)
	messagecommands.RegisterCommand("dog", "Comando dog", "pan!dog", "General", dogCommand)
}
