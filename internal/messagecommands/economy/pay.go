package economy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func payCommand(ctx *messagecommands.MessageContext, isGlobal bool) error {
	if len(ctx.Args) < 2 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el tipo de economía, el usuario y la cantidad.\nUso: `pan!pay <local|global> @usuario <cantidad>`")
		return err
	}

	ecoType := strings.ToLower(ctx.Args[0])
	if ecoType != "local" && ecoType != "global" {
		_, err := ctx.ReplyError("Uso Incorrecto", "El tipo de economía debe ser `local` o `global`.")
		return err
	}

	targetUserID := ctx.ParseUser(0)
	if targetUserID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un usuario válido.")
		return err
	}

	amount, err := strconv.ParseInt(ctx.Args[1], 10, 64)
	if err != nil || amount <= 0 {
		_, err := ctx.ReplyError("Error", "❌ La cantidad debe ser un número mayor a 0.")
		return err
	}

	userID := ctx.Message.Author.ID
	guildID := ctx.Message.GuildID

	if targetUserID == userID {
		_, err := ctx.ReplyError("Error", "❌ No puedes transferirte dinero a ti mismo.")
		return err
	}

	targetMember, err := ctx.Session.GuildMember(guildID, targetUserID)
	if err == nil && targetMember.User.Bot {
		_, err := ctx.ReplyError("Error", "❌ Los bots no tienen economía.")
		return err
	}

	if !isGlobal {
		err = database.TransferLocalBalance(guildID, userID, targetUserID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				_, err = ctx.ReplyError("Error", "❌ No tienes suficientes monedas locales en tu cartera.")
			} else {
				_, err = ctx.ReplyError("Error", "❌ Error al procesar la transferencia local.")
			}
			return err
		}
		_, err = ctx.ReplySuccess("Transferencia Exitosa", fmt.Sprintf("Has transferido **💵 %d** monedas locales a <@%s>.", amount, targetUserID))
		return err
	} else {
		err = database.TransferStars(userID, targetUserID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				_, err = ctx.ReplyError("Error", "❌ No tienes suficientes estrellas en tu cartera.")
			} else {
				_, err = ctx.ReplyError("Error", "❌ Error al procesar la transferencia estelar.")
			}
			return err
		}
		_, err = ctx.ReplySuccess("Transferencia Exitosa", fmt.Sprintf("Has transferido **🌟 %d** estrellas a <@%s>.", amount, targetUserID))
		return err
	}
}
