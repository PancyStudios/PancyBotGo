package economy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func withdrawCommand(ctx *messagecommands.MessageContext, isGlobal bool) error {
	if len(ctx.Args) < 1 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el tipo de economía y la cantidad.\nUso: `pan!withdraw <local|global> <cantidad>`")
		return err
	}

	ecoType := strings.ToLower(ctx.Args[0])
	if ecoType != "local" && ecoType != "global" {
		_, err := ctx.ReplyError("Uso Incorrecto", "El tipo de economía debe ser `local` o `global`.")
		return err
	}

	amount, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil || amount <= 0 {
		_, err := ctx.ReplyError("Error", "❌ La cantidad debe ser un número mayor a 0.")
		return err
	}

	userID := ctx.Message.Author.ID
	guildID := ctx.Message.GuildID

	if !isGlobal {
		err = database.WithdrawLocal(guildID, userID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				_, err = ctx.ReplyError("Error", "❌ No tienes suficientes monedas en el banco local.")
			} else {
				_, err = ctx.ReplyError("Error", "❌ Error al retirar.")
			}
			return err
		}
		_, err = ctx.ReplySuccess("Retiro Exitoso", fmt.Sprintf("Has retirado **💵 %d** de tu banco local.", amount))
		return err
	} else {
		err = database.WithdrawStars(userID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				_, err = ctx.ReplyError("Error", "❌ No tienes suficientes estrellas en el banco estelar.")
			} else {
				_, err = ctx.ReplyError("Error", "❌ Error al retirar.")
			}
			return err
		}
		_, err = ctx.ReplySuccess("Retiro Exitoso", fmt.Sprintf("Has retirado **🌟 %d** de tu banco estelar.", amount))
		return err
	}
}
