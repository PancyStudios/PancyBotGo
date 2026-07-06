package ia

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("createimage", "Comando createimage", "pan!createimage", "Ia", createImageCommand)
	messagecommands.RegisterCommand("getimage", "Comando getimage", "pan!getimage", "Ia", getImageCommand)
}
