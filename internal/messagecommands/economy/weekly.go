package economy

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func weeklyCommand(ctx *messagecommands.MessageContext) error {
	userID := ctx.Message.Author.ID

	cooldownDuration := 7 * 24 * time.Hour

	isReady, remaining, err := database.CooldownStars(userID, "weekly", cooldownDuration)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		_, err = ctx.ReplyError("Cooldown", fmt.Sprintf("❌ Ya reclamaste tu recompensa semanal. Vuelve en **%d días y %d horas**.", int(remaining.Hours()/24), int(remaining.Hours())%24))
		return err
	}

	amount := int64(10000)

	_, err = database.AddStars(userID, amount, false)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error al procesar la recompensa.")
		return err
	}

	_ = database.SetCooldownStars(userID, "weekly")

	_, err = ctx.ReplySuccess("Recompensa Semanal", fmt.Sprintf("¡Increíble! Has reclamado tu jugosa recompensa semanal de **🌟 %d estrellas**.", amount))
	return err
}
