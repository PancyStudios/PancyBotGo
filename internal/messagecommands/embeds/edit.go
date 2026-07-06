package embeds

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func editEmbedCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el ID del mensaje.\nUso: `pan!embededit <mensaje_id>`")
		return err
	}
	messageID := ctx.Args[0]

	msg, err := ctx.Session.ChannelMessage(ctx.Message.ChannelID, messageID)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ No pude encontrar el mensaje. Asegúrate de estar en el canal donde se envió y de que el ID sea correcto.")
		return err
	}

	if msg.Author.ID != ctx.Session.State.User.ID {
		_, err = ctx.ReplyError("Error", "❌ Solo puedo editar mis propios mensajes.")
		return err
	}

	if len(msg.Embeds) == 0 {
		_, err = ctx.ReplyError("Error", "❌ Este mensaje no contiene ningún embed.")
		return err
	}

	user := ctx.Message.Author
	saveBuilderState(user.ID, msg.Embeds[0])

	_, err = ctx.Reply("✅ Embed cargado con éxito. Ahora ejecuta `pan!embedcreate` para modificar sus propiedades y usar `pan!embedsend`.")
	return err
}
