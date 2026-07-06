package security

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("antibots", "Comando antibots", "pan!antibots", "Security", antibotsCommand)
	messagecommands.RegisterCommand("antiraid", "Comando antiraid", "pan!antiraid", "Security", antiraidCommand)
	messagecommands.RegisterCommand("verification", "Comando verification", "pan!verification", "Security", verificationCommand)
}
