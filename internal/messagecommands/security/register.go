package security

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("antibots", "Comando antibots", "pan!antibots", "General", antibotsCommand)
	messagecommands.RegisterCommand("antiraid", "Comando antiraid", "pan!antiraid", "General", antiraidCommand)
	messagecommands.RegisterCommand("verification", "Comando verification", "pan!verification", "General", verificationCommand)
}
