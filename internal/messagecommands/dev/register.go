package dev

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("eval", evalCommand)
	messagecommands.RegisterCommand("globalshop", globalShopCommand)
	messagecommands.RegisterCommand("blacklist", blacklistCommand)
	messagecommands.RegisterCommand("codegen", codegenCommand)
	messagecommands.RegisterCommand("codelist", codelistCommand)
	messagecommands.RegisterCommand("codedel", codedelCommand)
}
