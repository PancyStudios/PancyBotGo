package economy

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createWorkCommand(isGlobal bool) *discord.Command {
	return discord.NewCommand(
		"work",
		"💼 | Trabaja honradamente para ganar monedas",
		"economy",
		func(ctx *discord.CommandContext) error {
			return workHandler(ctx, isGlobal)
		},
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tipo",
			Description: "💰 | Elige si quieres ganar economía Local o Global",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{Name: "Local (Servidor)", Value: "local"},
				{Name: "Global (Estrellas)", Value: "global"},
			},
		},
	)
}

func workHandler(ctx *discord.CommandContext, isGlobal bool) error {
	
	userID := ctx.Interaction.Member.User.ID
	guildID := ctx.Interaction.GuildID

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
		ctx.Reply("❌ " + "Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		ctx.Reply("❌ " + fmt.Sprintf("Estás cansado. Vuelve a trabajar en **%d minutos y %d segundos**.", int(remaining.Minutes()), int(remaining.Seconds())%60))
		return nil
	}

	amount := int64(rand.Intn(200) + 50)

	if !isGlobal {
		_, err = database.AddLocalBalance(guildID, userID, amount, false)
		if err != nil {
			ctx.Reply("❌ " + "Error al procesar el pago local.")
			return err
		}
		_ = database.SetCooldownLocal(guildID, userID, "work")

		ctx.Reply(fmt.Sprintf("Has trabajado duro y ganaste **💵 %d monedas locales**.", amount))
	} else {
		_, err = database.AddStars(userID, amount, false)
		if err != nil {
			ctx.Reply("❌ " + "Error al procesar el pago global.")
			return err
		}
		_ = database.SetCooldownStars(userID, "work")

		ctx.Reply(fmt.Sprintf("Hiciste un viaje espacial y minaste **🌟 %d estrellas**.", amount))
	}
	return nil
}
