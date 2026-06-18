package economy

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// Register registers all economy commands
func Register(client *discord.ExtendedClient) {
	client.CommandHandler.RegisterCommand(createBalanceCommand())
	client.CommandHandler.RegisterCommand(createWorkCommand())
	client.CommandHandler.RegisterCommand(createDepositCommand())
	client.CommandHandler.RegisterCommand(createWithdrawCommand())
	client.CommandHandler.RegisterCommand(createPayCommand())
	client.CommandHandler.RegisterCommand(createDailyCommand())
	client.CommandHandler.RegisterCommand(createWeeklyCommand())
	client.CommandHandler.RegisterCommand(createShopCommand())
	client.CommandHandler.RegisterCommand(createAdminShopCommand())
	client.CommandHandler.RegisterCommand(createBuyCommand())
	client.CommandHandler.RegisterCommand(createUseCommand())
	client.CommandHandler.RegisterCommand(createInventoryCommand())
	client.CommandHandler.RegisterCommand(createCrimeCommand())
	client.CommandHandler.RegisterCommand(createSlutCommand())
	client.CommandHandler.RegisterCommand(createRobCommand())
	client.CommandHandler.RegisterCommand(createTopCommand())
}
