package economy

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createRobCommand(isGlobal bool) *discord.Command {
	return discord.NewCommand(
		"rob",
		"🥷 | Intenta robarle monedas a otro usuario",
		"economy",
		func(ctx *discord.CommandContext) error {
			return robHandler(ctx, isGlobal)
		},
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
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "victima",
			Description: "💰 | El usuario al que quieres robar",
			Required:    true,
		},
	)
}

func robHandler(ctx *discord.CommandContext, isGlobal bool) error {
	

	targetUser := ctx.GetUserOption("victima")
	userID := ctx.Interaction.Member.User.ID
	guildID := ctx.Interaction.GuildID

	if targetUser == nil || targetUser.ID == userID {
		ctx.Reply("❌ No puedes robarte a ti mismo.")
		return nil
	}
	if targetUser.Bot {
		ctx.Reply("❌ Los bots no tienen dinero, ni bolsillos.")
		return nil
	}

	cooldownDuration := 30 * time.Minute

	var isReady bool
	var remaining time.Duration
	var err error

	if !isGlobal {
		isReady, remaining, err = database.CooldownLocal(guildID, userID, "rob", cooldownDuration)
	} else {
		isReady, remaining, err = database.CooldownStars(userID, "rob", cooldownDuration)
	}

	if err != nil {
		ctx.Reply("❌ Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		ctx.Reply(fmt.Sprintf("❌ Tienes que esperar para planear tu próximo golpe. Vuelve en **%d minutos**.", int(remaining.Minutes())))
		return nil
	}

	// Calculate robbery logic
	success := rand.Float64() < 0.40 // 40% base success rate

	if !isGlobal {
		targetProfile, err := database.GetLocalProfile(guildID, targetUser.ID)
		if err != nil || targetProfile.Wallet < 100 {
			ctx.Reply("❌ Ese usuario no tiene dinero que valga la pena robar (Mínimo 100).")
			return nil
		}
		myProfile, _ := database.GetLocalProfile(guildID, userID)
		if myProfile.Wallet < 100 {
			ctx.Reply("❌ Necesitas al menos 100 monedas locales en tu cartera para cubrir posibles fianzas.")
			return nil
		}

		_ = database.SetCooldownLocal(guildID, userID, "rob")

		if success {
			// Steal between 10% to 40% of their wallet
			percent := (rand.Float64() * 0.3) + 0.1
			stolen := int64(float64(targetProfile.Wallet) * percent)

			database.AddLocalBalance(guildID, targetUser.ID, -stolen, false)
			database.AddLocalBalance(guildID, userID, stolen, false)
			ctx.Reply(fmt.Sprintf("🦹 ¡Éxito! Le robaste **💵 %d monedas** a %s.", stolen, targetUser.Mention()))
		} else {
			fine := int64(float64(myProfile.Wallet) * 0.25) // Pay 25% of your wallet
			if fine < 10 {
				fine = 10
			}
			database.AddLocalBalance(guildID, userID, -fine, false)
			database.AddLocalBalance(guildID, targetUser.ID, fine, false) // Give it to the victim as compensation
			ctx.Reply(fmt.Sprintf("🚔 ¡Te atraparon intentando robarle a %s! Tuviste que pagarle **💵 %d monedas** como multa.", targetUser.Mention(), fine))
		}
	} else {
		targetProfile, err := database.GetGlobalProfile(targetUser.ID)
		if err != nil || targetProfile.StarsWallet < 100 {
			ctx.Reply("❌ Ese usuario no tiene estrellas suficientes en la cartera (Mínimo 100).")
			return nil
		}
		myProfile, _ := database.GetGlobalProfile(userID)
		if myProfile.StarsWallet < 100 {
			ctx.Reply("❌ Necesitas al menos 100 estrellas en tu cartera para cubrir posibles fianzas.")
			return nil
		}

		_ = database.SetCooldownStars(userID, "rob")

		if success {
			percent := (rand.Float64() * 0.3) + 0.1
			stolen := int64(float64(targetProfile.StarsWallet) * percent)

			database.AddStars(targetUser.ID, -stolen, false)
			database.AddStars(userID, stolen, false)
			ctx.Reply(fmt.Sprintf("🦹 ¡Éxito! Le robaste **🌟 %d estrellas** a %s.", stolen, targetUser.Mention()))
		} else {
			fine := int64(float64(myProfile.StarsWallet) * 0.25)
			if fine < 10 {
				fine = 10
			}
			database.AddStars(userID, -fine, false)
			database.AddStars(targetUser.ID, fine, false)
			ctx.Reply(fmt.Sprintf("🚔 ¡Te atraparon robando estrellas de %s! Fuiste multado por **🌟 %d estrellas**, las cuales se le entregaron a tu víctima.", targetUser.Mention(), fine))
		}
	}

	return nil
}
