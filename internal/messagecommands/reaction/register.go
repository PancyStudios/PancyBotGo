package reaction

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("hug", "Comando hug", "pan!hug", "General", createReactionCommand("hug", "abrazó", "está dando abrazos", true))
	messagecommands.RegisterCommand("kiss", "Comando kiss", "pan!kiss", "General", createReactionCommand("kiss", "besó", "está dando besos", true))
	messagecommands.RegisterCommand("pat", "Comando pat", "pan!pat", "General", createReactionCommand("pat", "acarició", "está dando caricias", true))
	messagecommands.RegisterCommand("slap", "Comando slap", "pan!slap", "General", createReactionCommand("slap", "abofeteó", "está dando bofetadas", true))
	messagecommands.RegisterCommand("bite", "Comando bite", "pan!bite", "General", createReactionCommand("bite", "mordió", "está mordiendo", true))
	messagecommands.RegisterCommand("cuddle", "Comando cuddle", "pan!cuddle", "General", createReactionCommand("cuddle", "se acurrucó con", "se está acurrucando", true))
	messagecommands.RegisterCommand("cry", "Comando cry", "pan!cry", "General", createReactionCommand("cry", "lloró con", "está llorando", false))
	messagecommands.RegisterCommand("dance", "Comando dance", "pan!dance", "General", createReactionCommand("dance", "bailó con", "está bailando", false))
	messagecommands.RegisterCommand("happy", "Comando happy", "pan!happy", "General", createReactionCommand("happy", "sonrió con", "está muy feliz", false))
	messagecommands.RegisterCommand("smile", "Comando smile", "pan!smile", "General", createReactionCommand("smile", "le sonrió a", "está sonriendo", false))
	messagecommands.RegisterCommand("smug", "Comando smug", "pan!smug", "General", createReactionCommand("smug", "miró con presunción a", "tiene una mirada presumida", false))
	messagecommands.RegisterCommand("blush", "Comando blush", "pan!blush", "General", createReactionCommand("blush", "se sonrojó por", "está sonrojado/a", false))
	messagecommands.RegisterCommand("wink", "Comando wink", "pan!wink", "General", createReactionCommand("wink", "le guiñó el ojo a", "guiñó un ojo", false))
}
