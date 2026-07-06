package embeds

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("embedcreate", createEmbedCommand)
	messagecommands.RegisterCommand("embeddelete", deleteEmbedCommand)
	messagecommands.RegisterCommand("embededit", editEmbedCommand)
	messagecommands.RegisterCommand("embedsend", sendEmbedCommand)
}
