package dev

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("eval", "Comando eval", "pan!eval", "General", evalCommand)
	messagecommands.RegisterCommand("globalshop", "Comando globalshop", "pan!globalshop", "General", globalShopCommand)
	messagecommands.RegisterCommand("blacklist", "Comando blacklist", "pan!blacklist", "General", blacklistCommand)
	messagecommands.RegisterCommand("codegen", "Comando codegen", "pan!codegen", "General", codegenCommand)
	messagecommands.RegisterCommand("codelist", "Comando codelist", "pan!codelist", "General", codelistCommand)
	messagecommands.RegisterCommand("codedel", "Comando codedel", "pan!codedel", "General", codedelCommand)
}
