package fun

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register fun text commands
func Register() {
	messagecommands.RegisterCommand("8ball", eightBallCommand)
	messagecommands.RegisterCommand("ppt", pptCommand)
	messagecommands.RegisterCommand("ascii", asciiCommand)
	messagecommands.RegisterCommand("dog", dogCommand)
}
