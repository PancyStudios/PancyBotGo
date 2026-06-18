package fun

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// Register registers all fun commands as /fun subcommands
func Register(client *discord.ExtendedClient) {
	eightBallCmd := create8BallCommand()
	pptCmd := createPPTCommand()
	asciiCmd := createAsciiCommand()
	dogCmd := createDogCommand()

	funGroup := client.CommandHandler.BuildCommandGroup(
		"fun",
		"Comandos de Diversión y Minijuegos",
		eightBallCmd,
		pptCmd,
		asciiCmd,
		dogCmd,
	)

	client.CommandHandler.AddGlobalCommand(funGroup)
}
