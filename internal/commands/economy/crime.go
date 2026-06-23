package economy

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createCrimeCommand() *discord.Command {
	return discord.NewCommand(
		"crime",
		"🔫 | Comete un crimen (cuidado con la policía)",
		"economy",
		crimeHandler,
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

func crimeHandler(ctx *discord.CommandContext) error {
	ecoType := ctx.GetStringOption("tipo")
	userID := ctx.Interaction.Member.User.ID
	guildID := ctx.Interaction.GuildID

	cooldownDuration := 10 * time.Minute

	var isReady bool
	var remaining time.Duration
	var err error

	if ecoType == "local" {
		isReady, remaining, err = database.CooldownLocal(guildID, userID, "crime", cooldownDuration)
	} else {
		isReady, remaining, err = database.CooldownStars(userID, "crime", cooldownDuration)
	}

	if err != nil {
		ctx.Reply("❌ Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		ctx.Reply(fmt.Sprintf("❌ La policía te está buscando. Escóndete por **%d minutos y %d segundos**.", int(remaining.Minutes()), int(remaining.Seconds())%60))
		return nil
	}

	// 40% chance of success
	success := rand.Float64() < 0.40

	if ecoType == "local" {
		_ = database.SetCooldownLocal(guildID, userID, "crime")
		
		if success {
			amount := int64(rand.Intn(400) + 200)
			database.AddLocalBalance(guildID, userID, amount, false)
			ctx.Reply(fmt.Sprintf("🔪 Robaste una tienda y escapaste con **💵 %d monedas**.", amount))
		} else {
			fine := int64(rand.Intn(200) + 100)
			database.AddLocalBalance(guildID, userID, -fine, false) // Subtract money
			ctx.Reply(fmt.Sprintf("🚔 Te atraparon intentando robar una ancianita. Pagaste una fianza de **💵 %d monedas**.", fine))
		}
	} else {
		_ = database.SetCooldownStars(userID, "crime")
		
		if success {
			amount := int64(rand.Intn(300) + 150)
			database.AddStars(userID, amount, false)
			ctx.Reply(fmt.Sprintf("🔪 Hackeaste el banco intergaláctico y obtuviste **🌟 %d estrellas**.", amount))
		} else {
			fine := int64(rand.Intn(150) + 50)
			database.AddStars(userID, -fine, false)
			ctx.Reply(fmt.Sprintf("🚔 La patrulla espacial te pilló contrabandeando. Pagaste una multa de **🌟 %d estrellas**.", fine))
		}
	}
	return nil
}
