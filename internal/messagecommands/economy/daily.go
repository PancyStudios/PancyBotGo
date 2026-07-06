package economy

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func dailyCommand(ctx *messagecommands.MessageContext) error {
	userID := ctx.Message.Author.ID

	cooldownDuration := 24 * time.Hour

	isReady, remaining, err := database.CooldownStars(userID, "daily", cooldownDuration)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		_, err = ctx.ReplyError("Cooldown", fmt.Sprintf("❌ Ya reclamaste tu recompensa diaria. Vuelve en **%d horas y %d minutos**.", int(remaining.Hours()), int(remaining.Minutes())%60))
		return err
	}

	amount := int64(1000)

	_, err = database.AddStars(userID, amount, false)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error al procesar la recompensa.")
		return err
	}

	_ = database.SetCooldownStars(userID, "daily")

	_, err = ctx.ReplySuccess("Recompensa Diaria", fmt.Sprintf("¡Felicidades! Has reclamado tu recompensa diaria de **🌟 %d estrellas**.", amount))
	return err
}
