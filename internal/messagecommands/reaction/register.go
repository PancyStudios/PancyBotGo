package reaction

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func RegisterAll() {
	messagecommands.RegisterCommand("hug", createReactionCommand("hug", "abrazó", "está dando abrazos", true))
	messagecommands.RegisterCommand("kiss", createReactionCommand("kiss", "besó", "está dando besos", true))
	messagecommands.RegisterCommand("pat", createReactionCommand("pat", "acarició", "está dando caricias", true))
	messagecommands.RegisterCommand("slap", createReactionCommand("slap", "abofeteó", "está dando bofetadas", true))
	messagecommands.RegisterCommand("bite", createReactionCommand("bite", "mordió", "está mordiendo", true))
	messagecommands.RegisterCommand("cuddle", createReactionCommand("cuddle", "se acurrucó con", "se está acurrucando", true))
	messagecommands.RegisterCommand("cry", createReactionCommand("cry", "lloró con", "está llorando", false))
	messagecommands.RegisterCommand("dance", createReactionCommand("dance", "bailó con", "está bailando", false))
	messagecommands.RegisterCommand("happy", createReactionCommand("happy", "sonrió con", "está muy feliz", false))
	messagecommands.RegisterCommand("smile", createReactionCommand("smile", "le sonrió a", "está sonriendo", false))
	messagecommands.RegisterCommand("smug", createReactionCommand("smug", "miró con presunción a", "tiene una mirada presumida", false))
	messagecommands.RegisterCommand("blush", createReactionCommand("blush", "se sonrojó por", "está sonrojado/a", false))
	messagecommands.RegisterCommand("wink", createReactionCommand("wink", "le guiñó el ojo a", "guiñó un ojo", false))
}
