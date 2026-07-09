package reaction

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	// Requires target
	messagecommands.RegisterCommand("hug", "Comando hug", "pan!hug", "Reaction", createReactionCommand("hug", "abrazó a", "está dando abrazos", true))
	messagecommands.RegisterCommand("kiss", "Comando kiss", "pan!kiss", "Reaction", createReactionCommand("kiss", "besó a", "está dando besos", true))
	messagecommands.RegisterCommand("pat", "Comando pat", "pan!pat", "Reaction", createReactionCommand("pat", "acarició a", "está dando caricias", true))
	messagecommands.RegisterCommand("slap", "Comando slap", "pan!slap", "Reaction", createReactionCommand("slap", "abofeteó a", "está dando bofetadas", true))
	messagecommands.RegisterCommand("bite", "Comando bite", "pan!bite", "Reaction", createReactionCommand("bite", "mordió a", "está mordiendo", true))
	messagecommands.RegisterCommand("cuddle", "Comando cuddle", "pan!cuddle", "Reaction", createReactionCommand("cuddle", "se acurrucó con", "se está acurrucando", true))
	messagecommands.RegisterCommand("punch", "Comando punch", "pan!punch", "Reaction", createReactionCommand("punch", "golpeó a", "está dando golpes", true))
	messagecommands.RegisterCommand("poke", "Comando poke", "pan!poke", "Reaction", createReactionCommand("poke", "le dio un toquecito a", "está llamando la atención", true))
	messagecommands.RegisterCommand("lick", "Comando lick", "pan!lick", "Reaction", createReactionCommand("lick", "lamió a", "está lamiendo", true))
	messagecommands.RegisterCommand("brofist", "Comando brofist", "pan!brofist", "Reaction", createReactionCommand("brofist", "chocó los puños con", "está chocando puños", true))
	messagecommands.RegisterCommand("tickle", "Comando tickle", "pan!tickle", "Reaction", createReactionCommand("tickle", "le hizo cosquillas a", "está haciendo cosquillas", true))
	messagecommands.RegisterCommand("nom", "Comando nom", "pan!nom", "Reaction", createReactionCommand("nom", "mordisqueó a", "está mordisqueando", true))

	// No target
	messagecommands.RegisterCommand("cry", "Comando cry", "pan!cry", "Reaction", createReactionCommand("cry", "lloró con", "se puso a llorar 😢", false))
	messagecommands.RegisterCommand("dance", "Comando dance", "pan!dance", "Reaction", createReactionCommand("dance", "bailó con", "se puso a bailar 💃", false))
	messagecommands.RegisterCommand("happy", "Comando happy", "pan!happy", "Reaction", createReactionCommand("happy", "sonrió con", "está muy feliz ✨", false))
	messagecommands.RegisterCommand("smile", "Comando smile", "pan!smile", "Reaction", createReactionCommand("smile", "le sonrió a", "está sonriendo 🙂", false))
	messagecommands.RegisterCommand("smug", "Comando smug", "pan!smug", "Reaction", createReactionCommand("smug", "miró con presunción a", "tiene una mirada presumida 😏", false))
	messagecommands.RegisterCommand("blush", "Comando blush", "pan!blush", "Reaction", createReactionCommand("blush", "se sonrojó por", "se ha sonrojado 😳", false))
	messagecommands.RegisterCommand("wink", "Comando wink", "pan!wink", "Reaction", createReactionCommand("wink", "le guiñó el ojo a", "guiñó un ojo 😉", false))
	messagecommands.RegisterCommand("laugh", "Comando laugh", "pan!laugh", "Reaction", createReactionCommand("laugh", "se rio de", "se está riendo a carcajadas 😂", false))
	messagecommands.RegisterCommand("sigh", "Comando sigh", "pan!sigh", "Reaction", createReactionCommand("sigh", "suspiró por", "soltó un suspiro 😮‍💨", false))
	messagecommands.RegisterCommand("sleep", "Comando sleep", "pan!sleep", "Reaction", createReactionCommand("sleep", "se durmió con", "se quedó dormido 😴", false))
	messagecommands.RegisterCommand("pout", "Comando pout", "pan!pout", "Reaction", createReactionCommand("pout", "le hizo un puchero a", "hizo un puchero 😠", false))
	messagecommands.RegisterCommand("shrug", "Comando shrug", "pan!shrug", "Reaction", createReactionCommand("shrug", "se encogió de hombros ante", "se encogió de hombros 🤷", false))
	messagecommands.RegisterCommand("confused", "Comando confused", "pan!confused", "Reaction", createReactionCommand("confused", "se confundió con", "está muy confundido 😵‍💫", false))
}
