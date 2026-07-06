package dev

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func codedelCommand(ctx *messagecommands.MessageContext) error {
	if !isDev(ctx.Message.Author.ID) {
		_, err := ctx.ReplyError("Acceso Denegado", "Este comando es solo para la desarrolladora.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el código a eliminar.\nUso: `pan!codedel <codigo>`")
		return err
	}

	code := ctx.Args[0]

	_, err := database.GetPremiumCode(code)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("El código `%s` no existe.", code))
		return err
	}

	err = database.DeletePremiumCode(code)
	if err != nil {
		_, err = ctx.ReplyError("Error", "Error al eliminar el código.")
		return err
	}

	_, err = ctx.ReplySuccess("Código Eliminado", fmt.Sprintf("✅ El código `%s` ha sido eliminado.", code))
	return err
}
