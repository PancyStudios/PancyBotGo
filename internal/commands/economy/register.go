package economy

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// Register registers all economy commands into groups
func Register(client *discord.ExtendedClient) {
	// Create individual eco subcommands
	balanceCmd := createBalanceCommand()
	workCmd := createWorkCommand()
	depositCmd := createDepositCommand()
	withdrawCmd := createWithdrawCommand()
	payCmd := createPayCommand()
	dailyCmd := createDailyCommand()
	weeklyCmd := createWeeklyCommand()
	crimeCmd := createCrimeCommand()
	slutCmd := createSlutCommand()
	robCmd := createRobCommand()
	topCmd := createTopCommand()

	// Build the /eco command group
	ecoGroup := client.CommandHandler.BuildCommandGroup(
		"eco",
		"💰 Sistema de economía global y local",
		balanceCmd,
		workCmd,
		depositCmd,
		withdrawCmd,
		payCmd,
		dailyCmd,
		weeklyCmd,
		crimeCmd,
		slutCmd,
		robCmd,
		topCmd,
	)

	// Create individual shop subcommands
	shopViewCmd := createShopCommand()
	shopBuyCmd := createBuyCommand()
	shopUseCmd := createUseCommand()
	shopInvCmd := createInventoryCommand()
	shopAdminCmd := createAdminShopCommand()

	// Build the /shop command group
	shopGroup := client.CommandHandler.BuildCommandGroup(
		"shop",
		"🛒 Mercado de objetos y utilidades",
		shopViewCmd,
		shopBuyCmd,
		shopUseCmd,
		shopInvCmd,
		shopAdminCmd, // shopAdminCmd options are SubCommands, wait, a command can't have SubCommand inside SubCommandGroup in the root of a Slash Command if we build it as a top-level group.
		// Wait, BuildCommandGroup takes *discord.Command and turns them into SubCommands.
		// If admin_shop is already built with SubCommands (add/delete), it becomes a SubCommandGroup!
		// Discord allows: Root -> SubCommandGroup -> SubCommand.
	)

	// Register global groups
	client.CommandHandler.AddGlobalCommand(ecoGroup)
	client.CommandHandler.AddGlobalCommand(shopGroup)
}
