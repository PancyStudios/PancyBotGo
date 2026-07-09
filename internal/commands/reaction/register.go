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
	cuddleCmd := createReactionCommand("cuddle", "Acurrúcate con otro usuario", "se acurrucó con", "", true)
	punchCmd := createReactionCommand("punch", "Golpea a otro usuario", "golpeó a", "", true)
	pokeCmd := createReactionCommand("poke", "Llama la atención de otro usuario", "le dio un toquecito a", "", true)
	lickCmd := createReactionCommand("lick", "Lame a otro usuario", "lamió a", "", true)
	brofistCmd := createReactionCommand("brofist", "Choca los puños con otro usuario", "chocó los puños con", "", true)
	tickleCmd := createReactionCommand("tickle", "Hazle cosquillas a otro usuario", "le hizo cosquillas a", "", true)
	nomCmd := createReactionCommand("nom", "Muerde suavemente a otro usuario", "mordisqueó a", "", true)

	cryCmd := createReactionCommand("cry", "Ponte a llorar", "", "se puso a llorar 😢", false)
	danceCmd := createReactionCommand("dance", "Ponte a bailar", "", "se puso a bailar 💃", false)
	happyCmd := createReactionCommand("happy", "Demuestra tu felicidad", "", "está muy feliz ✨", false)
	smileCmd := createReactionCommand("smile", "Sonríe alegremente", "", "está sonriendo 🙂", false)
	smugCmd := createReactionCommand("smug", "Pon una mirada presumida", "", "tiene una mirada presumida 😏", false)
	blushCmd := createReactionCommand("blush", "Sonrójate", "", "se ha sonrojado 😳", false)
	winkCmd := createReactionCommand("wink", "Guiña un ojo", "", "guiñó un ojo 😉", false)
	laughCmd := createReactionCommand("laugh", "Ríete a carcajadas", "", "se está riendo a carcajadas 😂", false)
	sighCmd := createReactionCommand("sigh", "Suelta un suspiro", "", "soltó un suspiro 😮‍💨", false)
	sleepCmd := createReactionCommand("sleep", "Vete a dormir", "", "se quedó dormido 😴", false)
	poutCmd := createReactionCommand("pout", "Haz un puchero", "", "hizo un puchero 😠", false)
	shrugCmd := createReactionCommand("shrug", "Encógete de hombros", "", "se encogió de hombros 🤷", false)
	confusedCmd := createReactionCommand("confused", "Muestra confusión", "", "está muy confundido 😵‍💫", false)

	reactionGroup := client.CommandHandler.BuildCommandGroup(
		"reaccion",
		"Comandos de reacciones de anime",
		hugCmd,
		kissCmd,
		patCmd,
		slapCmd,
		biteCmd,
		cuddleCmd,
		punchCmd,
		pokeCmd,
		lickCmd,
		brofistCmd,
		tickleCmd,
		nomCmd,
		cryCmd,
		danceCmd,
		happyCmd,
		smileCmd,
		smugCmd,
		blushCmd,
		winkCmd,
		laughCmd,
		sighCmd,
		sleepCmd,
		poutCmd,
		shrugCmd,
		confusedCmd,
	)

	client.CommandHandler.AddGlobalCommand(reactionGroup)
}
