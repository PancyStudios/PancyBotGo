package security

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("antibots", antibotsCommand)
	messagecommands.RegisterCommand("antiraid", antiraidCommand)
	messagecommands.RegisterCommand("verification", verificationCommand)
}
