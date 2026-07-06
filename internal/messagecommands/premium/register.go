package premium

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("redeem", "Comando redeem", "pan!redeem", "Premium", redeemCommand)
}
