package economy

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// Register registers all economy commands into groups
func Register(client *discord.ExtendedClient) {
	// Create individual eco subcommands (Global)
	balanceGlobal := createBalanceCommand(true)
	workGlobal := createWorkCommand(true)
	depositGlobal := createDepositCommand(true)
	withdrawGlobal := createWithdrawCommand(true)
	payGlobal := createPayCommand(true)
	dailyGlobal := createDailyCommand(true)
	weeklyGlobal := createWeeklyCommand(true)
	crimeGlobal := createCrimeCommand(true)
	slutGlobal := createSlutCommand(true)
	robGlobal := createRobCommand(true)
	topGlobal := createTopCommand(true)

	// Create individual eco subcommands (Local)
	balanceLocal := createBalanceCommand(false)
	workLocal := createWorkCommand(false)
	depositLocal := createDepositCommand(false)
	withdrawLocal := createWithdrawCommand(false)
	payLocal := createPayCommand(false)
	dailyLocal := createDailyCommand(false)
	weeklyLocal := createWeeklyCommand(false)
	crimeLocal := createCrimeCommand(false)
	slutLocal := createSlutCommand(false)
	robLocal := createRobCommand(false)
	topLocal := createTopCommand(false)

	// Build the /eco command group
	ecoGroup := client.CommandHandler.BuildCommandGroup(
		"eco",
		"🌟 Sistema de economía global (Estrellas)",
		balanceGlobal,
		workGlobal,
		depositGlobal,
		withdrawGlobal,
		payGlobal,
		dailyGlobal,
		weeklyGlobal,
		crimeGlobal,
		slutGlobal,
		robGlobal,
		topGlobal,
	)

	// Build the /ecol command group
	ecolGroup := client.CommandHandler.BuildCommandGroup(
		"ecol",
		"💵 Sistema de economía local (Servidor)",
		balanceLocal,
		workLocal,
		depositLocal,
		withdrawLocal,
		payLocal,
		dailyLocal,
		weeklyLocal,
		crimeLocal,
		slutLocal,
		robLocal,
		topLocal,
	)

	// Create unified shop subcommands
	shopView := createShopCommand()
	shopBuy := createBuyCommand()
	shopUse := createUseCommand()
	shopInv := createInventoryCommand()
	shopAdmin := createAdminShopCommand()

	// Build the /shop command group
	shopGroup := client.CommandHandler.BuildCommandGroup(
		"shop",
		"🛒 Tienda de objetos",
		shopView,
		shopBuy,
		shopUse,
		shopInv,
		shopAdmin,
	)

	// Register global groups
	client.CommandHandler.AddGlobalCommand(ecoGroup)
	client.CommandHandler.AddGlobalCommand(ecolGroup)
	client.CommandHandler.AddGlobalCommand(shopGroup)
}
