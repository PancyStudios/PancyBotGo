package embeds

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("embedcreate", "Comando embedcreate", "pan!embedcreate", "General", createEmbedCommand)
	messagecommands.RegisterCommand("embeddelete", "Comando embeddelete", "pan!embeddelete", "General", deleteEmbedCommand)
	messagecommands.RegisterCommand("embededit", "Comando embededit", "pan!embededit", "General", editEmbedCommand)
	messagecommands.RegisterCommand("embedsend", "Comando embedsend", "pan!embedsend", "General", sendEmbedCommand)
}
