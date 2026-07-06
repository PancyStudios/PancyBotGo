package economy

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

func createDailyCommand(isGlobal bool) *discord.Command {
	return discord.NewCommand(
		"daily",
		"📅 | Reclama tu recompensa diaria",
		"economy",
		func(ctx *discord.CommandContext) error {
			return dailyHandler(ctx, isGlobal)
		},
	)
}

func dailyHandler(ctx *discord.CommandContext, isGlobal bool) error {
	userID := ctx.Interaction.Member.User.ID

	cooldownDuration := 24 * time.Hour

	isReady, remaining, err := database.CooldownStars(userID, "daily", cooldownDuration)
	if err != nil {
		ctx.Reply("❌ " + "Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		ctx.Reply("❌ " + fmt.Sprintf("Ya reclamaste tu recompensa diaria. Vuelve en **%d horas y %d minutos**.", int(remaining.Hours()), int(remaining.Minutes())%60))
		return nil
	}

	amount := int64(1000)

	_, err = database.AddStars(userID, amount, false)
	if err != nil {
		ctx.Reply("❌ " + "Error al procesar la recompensa.")
		return err
	}

	_ = database.SetCooldownStars(userID, "daily")

	ctx.Reply(fmt.Sprintf("¡Felicidades! Has reclamado tu recompensa diaria de **🌟 %d estrellas**.", amount))
	return nil
}
