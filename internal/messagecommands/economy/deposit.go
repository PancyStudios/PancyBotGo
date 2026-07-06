package economy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func depositCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) < 2 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el tipo de economía y la cantidad.\nUso: `pan!deposit <local|global> <cantidad>`")
		return err
	}

	ecoType := strings.ToLower(ctx.Args[0])
	if ecoType != "local" && ecoType != "global" {
		_, err := ctx.ReplyError("Uso Incorrecto", "El tipo de economía debe ser `local` o `global`.")
		return err
	}

	amount, err := strconv.ParseInt(ctx.Args[1], 10, 64)
	if err != nil || amount <= 0 {
		_, err := ctx.ReplyError("Error", "❌ La cantidad debe ser un número mayor a 0.")
		return err
	}

	userID := ctx.Message.Author.ID
	guildID := ctx.Message.GuildID

	if ecoType == "local" {
		err = database.DepositLocal(guildID, userID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				_, err = ctx.ReplyError("Error", "❌ No tienes suficientes monedas locales en tu cartera.")
			} else if err == database.ErrBankFull {
				_, err = ctx.ReplyError("Error", "❌ El banco local no tiene suficiente capacidad para ese depósito.")
			} else {
				_, err = ctx.ReplyError("Error", "❌ Error al depositar.")
			}
			return err
		}
		_, err = ctx.ReplySuccess("Depósito Exitoso", fmt.Sprintf("Has depositado **💵 %d** a tu banco local.", amount))
		return err
	} else {
		err = database.DepositStars(userID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				_, err = ctx.ReplyError("Error", "❌ No tienes suficientes estrellas en tu cartera.")
			} else if err == database.ErrBankFull {
				_, err = ctx.ReplyError("Error", "❌ Tu banco estelar está al límite de su capacidad.")
			} else {
				_, err = ctx.ReplyError("Error", "❌ Error al depositar.")
			}
			return err
		}
		_, err = ctx.ReplySuccess("Depósito Exitoso", fmt.Sprintf("Has depositado **🌟 %d** a tu banco estelar.", amount))
		return err
	}
}
