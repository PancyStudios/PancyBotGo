package economy

import "github.com/PancyStudios/PancyBotGo/internal/messagecommands"

// Register economy text commands
func Register() {
	messagecommands.RegisterCommand("work", "Comando work", "pan!work", "General", workCommand)
	messagecommands.RegisterCommand("rob", "Comando rob", "pan!rob", "General", robCommand)
	messagecommands.RegisterCommand("slut", "Comando slut", "pan!slut", "General", slutCommand)
	messagecommands.RegisterCommand("crime", "Comando crime", "pan!crime", "General", crimeCommand)
	messagecommands.RegisterCommand("buy", "Comando buy", "pan!buy", "General", buyCommand)
	messagecommands.RegisterCommand("deposit", "Comando deposit", "pan!deposit", "General", depositCommand)
	messagecommands.RegisterCommand("withdraw", "Comando withdraw", "pan!withdraw", "General", withdrawCommand)
	messagecommands.RegisterCommand("pay", "Comando pay", "pan!pay", "General", payCommand)
	messagecommands.RegisterCommand("top", "Comando top", "pan!top", "General", topCommand)
	messagecommands.RegisterCommand("use", "Comando use", "pan!use", "General", useCommand)
	messagecommands.RegisterCommand("balance", "Comando balance", "pan!balance", "General", balanceCommand)
	messagecommands.RegisterCommand("daily", "Comando daily", "pan!daily", "General", dailyCommand)
	messagecommands.RegisterCommand("weekly", "Comando weekly", "pan!weekly", "General", weeklyCommand)
	messagecommands.RegisterCommand("inventory", "Comando inventory", "pan!inventory", "General", inventoryCommand)
	messagecommands.RegisterCommand("shop", "Comando shop", "pan!shop", "General", shopCommand)
	messagecommands.RegisterCommand("adminshop", "Comando adminshop", "pan!adminshop", "General", adminShopCommand)
}
