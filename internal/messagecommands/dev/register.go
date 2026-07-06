package dev

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("eval", "Comando eval", "pan!eval", "Dev", evalCommand)
	messagecommands.RegisterCommand("globalshop", "Comando globalshop", "pan!globalshop", "Dev", globalShopCommand)
	messagecommands.RegisterCommand("blacklist", "Comando blacklist", "pan!blacklist", "Dev", blacklistCommand)
	messagecommands.RegisterCommand("codegen", "Comando codegen", "pan!codegen", "Dev", codegenCommand)
	messagecommands.RegisterCommand("codelist", "Comando codelist", "pan!codelist", "Dev", codelistCommand)
	messagecommands.RegisterCommand("codedel", "Comando codedel", "pan!codedel", "Dev", codedelCommand)
}
