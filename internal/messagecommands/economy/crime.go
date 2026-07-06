package economy

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func crimeCommand(ctx *messagecommands.MessageContext, isGlobal bool) error {

	userID := ctx.Message.Author.ID
	guildID := ctx.Message.GuildID

	cooldownDuration := 10 * time.Minute

	var isReady bool
	var remaining time.Duration
	var err error

	if !isGlobal {
		isReady, remaining, err = database.CooldownLocal(guildID, userID, "crime", cooldownDuration)
	} else {
		isReady, remaining, err = database.CooldownStars(userID, "crime", cooldownDuration)
	}

	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		_, err = ctx.ReplyError("Cooldown", fmt.Sprintf("❌ La policía te está buscando. Escóndete por **%d minutos y %d segundos**.", int(remaining.Minutes()), int(remaining.Seconds())%60))
		return err
	}

	success := rand.Float64() < 0.40

	if !isGlobal {
		_ = database.SetCooldownLocal(guildID, userID, "crime")

		if success {
			amount := int64(rand.Intn(400) + 200)
			database.AddLocalBalance(guildID, userID, amount, false)
			_, err = ctx.ReplySuccess("Crimen Exitoso", fmt.Sprintf("🔪 Robaste una tienda y escapaste con **💵 %d monedas**.", amount))
			return err
		} else {
			fine := int64(rand.Intn(200) + 100)
			database.AddLocalBalance(guildID, userID, -fine, false)
			_, err = ctx.ReplyError("Atrapado", fmt.Sprintf("🚔 Te atraparon intentando robar una ancianita. Pagaste una fianza de **💵 %d monedas**.", fine))
			return err
		}
	} else {
		_ = database.SetCooldownStars(userID, "crime")

		if success {
			amount := int64(rand.Intn(300) + 150)
			database.AddStars(userID, amount, false)
			_, err = ctx.ReplySuccess("Crimen Exitoso", fmt.Sprintf("🔪 Hackeaste el banco intergaláctico y obtuviste **🌟 %d estrellas**.", amount))
			return err
		} else {
			fine := int64(rand.Intn(150) + 50)
			database.AddStars(userID, -fine, false)
			_, err = ctx.ReplyError("Atrapado", fmt.Sprintf("🚔 La patrulla espacial te pilló contrabandeando. Pagaste una multa de **🌟 %d estrellas**.", fine))
			return err
		}
	}
}
