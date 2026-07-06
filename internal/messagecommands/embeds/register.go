package embeds

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("embedcreate", "Comando embedcreate", "pan!embedcreate", "Embeds", createEmbedCommand)
	messagecommands.RegisterCommand("embeddelete", "Comando embeddelete", "pan!embeddelete", "Embeds", deleteEmbedCommand)
	messagecommands.RegisterCommand("embededit", "Comando embededit", "pan!embededit", "Embeds", editEmbedCommand)
	messagecommands.RegisterCommand("embedsend", "Comando embedsend", "pan!embedsend", "Embeds", sendEmbedCommand)
}
