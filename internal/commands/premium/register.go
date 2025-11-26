package premium

import "github.com/PancyStudios/PancyBotGo/pkg/discord"

// Register registers all premium commands as /premium subcommands
func Register(client *discord.ExtendedClient) {
	// Create individual subcommands
	redeemCmd := CreateRedeemCommand()

	// Build the /premium command group with all subcommands
	premiumGroup := client.CommandHandler.BuildCommandGroup(
		"premium",
		"Comandos premium",
		redeemCmd,
	)

	// Register the command group
	client.CommandHandler.AddGlobalCommand(premiumGroup)
}
