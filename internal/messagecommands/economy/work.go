package economy

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func workCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el tipo de economía.\nUso: `pan!work <local|global>`")
		return err
	}

	ecoType := strings.ToLower(ctx.Args[0])
	if ecoType != "local" && ecoType != "global" {
		_, err := ctx.ReplyError("Uso Incorrecto", "El tipo de economía debe ser `local` o `global`.")
		return err
	}

	userID := ctx.Message.Author.ID
	guildID := ctx.Message.GuildID

	cooldownDuration := 5 * time.Minute

	var isReady bool
	var remaining time.Duration
	var err error

	if ecoType == "local" {
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

	if ecoType == "local" {
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
