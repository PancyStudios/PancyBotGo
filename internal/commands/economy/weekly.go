package economy

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

func createWeeklyCommand(isGlobal bool) *discord.Command {
	return discord.NewCommand(
		"weekly",
		"📆 | Reclama tu recompensa semanal",
		"economy",
		func(ctx *discord.CommandContext) error {
			return weeklyHandler(ctx, isGlobal)
		},
	)
}

func weeklyHandler(ctx *discord.CommandContext, isGlobal bool) error {
	userID := ctx.Interaction.Member.User.ID

	cooldownDuration := 7 * 24 * time.Hour

	isReady, remaining, err := database.CooldownStars(userID, "weekly", cooldownDuration)
	if err != nil {
		ctx.Reply("❌ Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		ctx.Reply(fmt.Sprintf("❌ Ya reclamaste tu recompensa semanal. Vuelve en **%d días y %d horas**.", int(remaining.Hours()/24), int(remaining.Hours())%24))
		return nil
	}

	amount := int64(10000) // Weekly gives 10k stars

	_, err = database.AddStars(userID, amount, false)
	if err != nil {
		ctx.Reply("❌ Error al procesar la recompensa.")
		return err
	}

	_ = database.SetCooldownStars(userID, "weekly")

	ctx.Reply(fmt.Sprintf("¡Increíble! Has reclamado tu jugosa recompensa semanal de **🌟 %d estrellas**.", amount))
	return nil
}
