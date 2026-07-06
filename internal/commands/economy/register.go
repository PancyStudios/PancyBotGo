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

	// Build the /ecog command group
	ecogGroup := client.CommandHandler.BuildCommandGroup(
		"ecog",
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

	// Create individual shop subcommands (Global)
	shopViewGlobal := createShopCommand(true)
	shopBuyGlobal := createBuyCommand(true)
	shopUseGlobal := createUseCommand(true)
	shopInvGlobal := createInventoryCommand(true)
	shopAdminGlobal := createAdminShopCommand(true)

	// Create individual shop subcommands (Local)
	shopViewLocal := createShopCommand(false)
	shopBuyLocal := createBuyCommand(false)
	shopUseLocal := createUseCommand(false)
	shopInvLocal := createInventoryCommand(false)
	shopAdminLocal := createAdminShopCommand(false)

	// Build the /shop command group
	shopGroup := client.CommandHandler.BuildCommandGroup(
		"shop",
		"🌟 Tienda global de objetos",
		shopViewGlobal,
		shopBuyGlobal,
		shopUseGlobal,
		shopInvGlobal,
		shopAdminGlobal,
	)

	// Build the /shopg command group
	shopgGroup := client.CommandHandler.BuildCommandGroup(
		"shopg",
		"🛒 Tienda local del servidor",
		shopViewLocal,
		shopBuyLocal,
		shopUseLocal,
		shopInvLocal,
		shopAdminLocal,
	)

	// Register global groups
	client.CommandHandler.AddGlobalCommand(ecoGroup)
	client.CommandHandler.AddGlobalCommand(ecogGroup)
	client.CommandHandler.AddGlobalCommand(shopGroup)
	client.CommandHandler.AddGlobalCommand(shopgGroup)
}
