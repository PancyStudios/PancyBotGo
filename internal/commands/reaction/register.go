package reaction

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// RegisterReactionCommands registers all the reaction commands
func RegisterReactionCommands(client *discord.ExtendedClient) {
	hugCmd := createReactionCommand("hug", "Abraza a otro usuario", "abrazó", "", true)
	kissCmd := createReactionCommand("kiss", "Besa a otro usuario", "besó", "", true)
	patCmd := createReactionCommand("pat", "Acaricia a otro usuario", "acarició", "", true)
	slapCmd := createReactionCommand("slap", "Abofetea a otro usuario", "abofeteó", "", true)
	biteCmd := createReactionCommand("bite", "Muerde a otro usuario", "mordió", "", true)

	// Single target actions (no requires target, or optional)
	// Waifu.pics has "cry", "smile", "happy", "dance", "cringe"
	cryCmd := createReactionCommand("cry", "Ponte a llorar", "", "se puso a llorar 😢", false)
	danceCmd := createReactionCommand("dance", "Ponte a bailar", "", "se puso a bailar 💃", false)

	reactionGroup := client.CommandHandler.BuildCommandGroup(
		"reaccion",
		"Comandos de reacciones de anime",
		hugCmd,
		kissCmd,
		patCmd,
		slapCmd,
		biteCmd,
		cryCmd,
		danceCmd,
	)

	client.CommandHandler.AddGlobalCommand(reactionGroup)
}
