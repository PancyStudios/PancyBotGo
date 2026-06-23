package economy

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createSlutCommand() *discord.Command {
	return discord.NewCommand(
		"slut",
		"👠 | Trabaja en las calles (alto riesgo)",
		"economy",
		slutHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tipo",
			Description: "💰 | Elige economía Local o Global",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{Name: "Local (Servidor)", Value: "local"},
				{Name: "Global (Estrellas)", Value: "global"},
			},
		},
	)
}

func slutHandler(ctx *discord.CommandContext) error {
	ecoType := ctx.GetStringOption("tipo")
	userID := ctx.Interaction.Member.User.ID
	guildID := ctx.Interaction.GuildID

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
		ctx.Reply("❌ Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		ctx.Reply(fmt.Sprintf("❌ Aún te duelen las caderas. Descansa por **%d minutos y %d segundos**.", int(remaining.Minutes()), int(remaining.Seconds())%60))
		return nil
	}

	// 60% chance of success
	success := rand.Float64() < 0.60

	if ecoType == "local" {
		_ = database.SetCooldownLocal(guildID, userID, "slut")
		
		if success {
			amount := int64(rand.Intn(500) + 100)
			database.AddLocalBalance(guildID, userID, amount, false)
			ctx.Reply(fmt.Sprintf("💋 Te fue excelente en la esquina y te pagaron **💵 %d monedas**.", amount))
		} else {
			fine := int64(rand.Intn(100) + 50)
			database.AddLocalBalance(guildID, userID, -fine, false) // Subtract money
			ctx.Reply(fmt.Sprintf("🚔 Te asaltaron en el callejón. Perdiste **💵 %d monedas**.", fine))
		}
	} else {
		_ = database.SetCooldownStars(userID, "slut")
		
		if success {
			amount := int64(rand.Intn(400) + 100)
			database.AddStars(userID, amount, false)
			ctx.Reply(fmt.Sprintf("💋 Conseguiste un Sugar Alien que te donó **🌟 %d estrellas**.", amount))
		} else {
			fine := int64(rand.Intn(100) + 20)
			database.AddStars(userID, -fine, false)
			ctx.Reply(fmt.Sprintf("🚔 Te arrestó la patrulla del espacio por exhibicionismo. Pagaste **🌟 %d estrellas** de multa.", fine))
		}
	}
	return nil
}
