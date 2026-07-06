package economy

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func slutCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el tipo de economía.\nUso: `pan!slut <local|global>`")
		return err
	}

	ecoType := strings.ToLower(ctx.Args[0])
	if ecoType != "local" && ecoType != "global" {
		_, err := ctx.ReplyError("Uso Incorrecto", "El tipo de economía debe ser `local` o `global`.")
		return err
	}

	userID := ctx.Message.Author.ID
	guildID := ctx.Message.GuildID

	cooldownDuration := 15 * time.Minute

	var isReady bool
	var remaining time.Duration
	var err error

	if ecoType == "local" {
		isReady, remaining, err = database.CooldownLocal(guildID, userID, "slut", cooldownDuration)
	} else {
		isReady, remaining, err = database.CooldownStars(userID, "slut", cooldownDuration)
	}

	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		_, err = ctx.ReplyError("Cooldown", fmt.Sprintf("❌ Aún te duelen las caderas. Descansa por **%d minutos y %d segundos**.", int(remaining.Minutes()), int(remaining.Seconds())%60))
		return err
	}

	success := rand.Float64() < 0.60

	if ecoType == "local" {
		_ = database.SetCooldownLocal(guildID, userID, "slut")

		if success {
			amount := int64(rand.Intn(500) + 100)
			database.AddLocalBalance(guildID, userID, amount, false)
			_, err = ctx.ReplySuccess("Trabajo Terminado", fmt.Sprintf("💋 Te fue excelente en la esquina y te pagaron **💵 %d monedas**.", amount))
			return err
		} else {
			fine := int64(rand.Intn(100) + 50)
			database.AddLocalBalance(guildID, userID, -fine, false)
			_, err = ctx.ReplyError("Atrapado", fmt.Sprintf("🚔 Te asaltaron en el callejón. Perdiste **💵 %d monedas**.", fine))
			return err
		}
	} else {
		_ = database.SetCooldownStars(userID, "slut")

		if success {
			amount := int64(rand.Intn(400) + 100)
			database.AddStars(userID, amount, false)
			_, err = ctx.ReplySuccess("Trabajo Terminado", fmt.Sprintf("💋 Conseguiste un Sugar Alien que te donó **🌟 %d estrellas**.", amount))
			return err
		} else {
			fine := int64(rand.Intn(100) + 20)
			database.AddStars(userID, -fine, false)
			_, err = ctx.ReplyError("Atrapado", fmt.Sprintf("🚔 Te arrestó la patrulla del espacio por exhibicionismo. Pagaste **🌟 %d estrellas** de multa.", fine))
			return err
		}
	}
}
