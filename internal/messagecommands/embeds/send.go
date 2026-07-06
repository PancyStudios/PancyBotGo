package embeds

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
)

func sendEmbedCommand(ctx *messagecommands.MessageContext) error {
	user := ctx.Message.Author

	builderMutex.RLock()
	embedState, exists := builderStateMap[user.ID]
	builderMutex.RUnlock()

	if !exists {
		_, err := ctx.ReplyError("Error", "❌ No tienes ningún embed en construcción. Usa `pan!embedcreate` primero.")
		return err
	}

	targetChannelID := ctx.Message.ChannelID
	parsedChannel := ctx.ParseChannel(0)
	if parsedChannel != "" {
		targetChannelID = parsedChannel
	}

	_, err := ctx.Session.ChannelMessageSendEmbed(targetChannelID, embedState)
	if err != nil {
		logger.Error(fmt.Sprintf("Error sending embed: %v", err), "Embeds")
		_, err = ctx.ReplyError("Error", "❌ No pude enviar el embed. Verifica que tengo permisos de escribir y enviar embeds en ese canal.")
		return err
	}

	_, err = ctx.ReplySuccess("Embed Enviado", fmt.Sprintf("✅ Embed enviado exitosamente a <#%s>", targetChannelID))
	return err
}
