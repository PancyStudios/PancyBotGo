package ia

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("createimage", createImageCommand)
	messagecommands.RegisterCommand("getimage", getImageCommand)
}
