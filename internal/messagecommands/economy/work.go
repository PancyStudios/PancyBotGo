package economy

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func workCommand(ctx *messagecommands.MessageContext, isGlobal bool) error {

	userID := ctx.Message.Author.ID
	guildID := ctx.Message.GuildID

	cooldownDuration := 5 * time.Minute

	var isReady bool
	var remaining time.Duration
	var err error

	if !isGlobal {
		isReady, remaining, err = database.CooldownLocal(guildID, userID, "work", cooldownDuration)
	} else {
		isReady, remaining, err = database.CooldownStars(userID, "work", cooldownDuration)
	}

	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		_, err = ctx.ReplyError("Cooldown", fmt.Sprintf("❌ Estás cansado. Vuelve a trabajar en **%d minutos y %d segundos**.", int(remaining.Minutes()), int(remaining.Seconds())%60))
		return err
	}

	amount := int64(rand.Intn(200) + 50)

	if !isGlobal {
		_, err = database.AddLocalBalance(guildID, userID, amount, false)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ Error al procesar el pago local.")
			return err
		}
		_ = database.SetCooldownLocal(guildID, userID, "work")

		_, err = ctx.ReplySuccess("Trabajo Terminado", fmt.Sprintf("Has trabajado duro y ganaste **💵 %d monedas locales**.", amount))
		return err
	} else {
		_, err = database.AddStars(userID, amount, false)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ Error al procesar el pago global.")
			return err
		}
		_ = database.SetCooldownStars(userID, "work")

		_, err = ctx.ReplySuccess("Trabajo Terminado", fmt.Sprintf("Hiciste un viaje espacial y minaste **🌟 %d estrellas**.", amount))
		return err
	}
}
