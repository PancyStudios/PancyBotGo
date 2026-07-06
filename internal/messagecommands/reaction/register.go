package reaction

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("hug", "Comando hug", "pan!hug", "Reaction", createReactionCommand("hug", "abrazó", "está dando abrazos", true))
	messagecommands.RegisterCommand("kiss", "Comando kiss", "pan!kiss", "Reaction", createReactionCommand("kiss", "besó", "está dando besos", true))
	messagecommands.RegisterCommand("pat", "Comando pat", "pan!pat", "Reaction", createReactionCommand("pat", "acarició", "está dando caricias", true))
	messagecommands.RegisterCommand("slap", "Comando slap", "pan!slap", "Reaction", createReactionCommand("slap", "abofeteó", "está dando bofetadas", true))
	messagecommands.RegisterCommand("bite", "Comando bite", "pan!bite", "Reaction", createReactionCommand("bite", "mordió", "está mordiendo", true))
	messagecommands.RegisterCommand("cuddle", "Comando cuddle", "pan!cuddle", "Reaction", createReactionCommand("cuddle", "se acurrucó con", "se está acurrucando", true))
	messagecommands.RegisterCommand("cry", "Comando cry", "pan!cry", "Reaction", createReactionCommand("cry", "lloró con", "está llorando", false))
	messagecommands.RegisterCommand("dance", "Comando dance", "pan!dance", "Reaction", createReactionCommand("dance", "bailó con", "está bailando", false))
	messagecommands.RegisterCommand("happy", "Comando happy", "pan!happy", "Reaction", createReactionCommand("happy", "sonrió con", "está muy feliz", false))
	messagecommands.RegisterCommand("smile", "Comando smile", "pan!smile", "Reaction", createReactionCommand("smile", "le sonrió a", "está sonriendo", false))
	messagecommands.RegisterCommand("smug", "Comando smug", "pan!smug", "Reaction", createReactionCommand("smug", "miró con presunción a", "tiene una mirada presumida", false))
	messagecommands.RegisterCommand("blush", "Comando blush", "pan!blush", "Reaction", createReactionCommand("blush", "se sonrojó por", "está sonrojado/a", false))
	messagecommands.RegisterCommand("wink", "Comando wink", "pan!wink", "Reaction", createReactionCommand("wink", "le guiñó el ojo a", "guiñó un ojo", false))
}
