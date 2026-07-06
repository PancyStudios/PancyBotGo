package economy

import "github.com/PancyStudios/PancyBotGo/internal/messagecommands"

// Register economy text commands
func Register() {
	messagecommands.RegisterCommand("work", workCommand)
	messagecommands.RegisterCommand("rob", robCommand)
	messagecommands.RegisterCommand("slut", slutCommand)
	messagecommands.RegisterCommand("crime", crimeCommand)
	messagecommands.RegisterCommand("buy", buyCommand)
	messagecommands.RegisterCommand("deposit", depositCommand)
	messagecommands.RegisterCommand("withdraw", withdrawCommand)
	messagecommands.RegisterCommand("pay", payCommand)
	messagecommands.RegisterCommand("top", topCommand)
	messagecommands.RegisterCommand("use", useCommand)
	messagecommands.RegisterCommand("balance", balanceCommand)
	messagecommands.RegisterCommand("daily", dailyCommand)
	messagecommands.RegisterCommand("weekly", weeklyCommand)
	messagecommands.RegisterCommand("inventory", inventoryCommand)
	messagecommands.RegisterCommand("shop", shopCommand)
	messagecommands.RegisterCommand("adminshop", adminShopCommand)
}
