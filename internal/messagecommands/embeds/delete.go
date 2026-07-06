package embeds

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func deleteEmbedCommand(ctx *messagecommands.MessageContext) error {
	user := ctx.Message.Author
	clearBuilderState(user.ID)
	_, err := ctx.Reply("🗑️ El embed en el que estabas trabajando ha sido descartado.")
	return err
}
