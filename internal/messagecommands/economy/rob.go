package economy

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func robCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) < 2 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el tipo de economía y el usuario.\nUso: `pan!rob <local|global> @usuario`")
		return err
	}

	ecoType := strings.ToLower(ctx.Args[0])
	if ecoType != "local" && ecoType != "global" {
		_, err := ctx.ReplyError("Uso Incorrecto", "El tipo de economía debe ser `local` o `global`.")
		return err
	}

	targetUserID := ctx.ParseUser(1)
	if targetUserID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un usuario válido para robar.")
		return err
	}

	userID := ctx.Message.Author.ID
	guildID := ctx.Message.GuildID

	if targetUserID == userID {
		_, err := ctx.ReplyError("Error", "❌ No puedes robarte a ti mismo.")
		return err
	}

	targetMember, err := ctx.Session.GuildMember(guildID, targetUserID)
	if err == nil && targetMember.User.Bot {
		_, err := ctx.ReplyError("Error", "❌ Los bots no tienen dinero, ni bolsillos.")
		return err
	}

	cooldownDuration := 30 * time.Minute

	var isReady bool
	var remaining time.Duration

	if ecoType == "local" {
		isReady, remaining, err = database.CooldownLocal(guildID, userID, "rob", cooldownDuration)
	} else {
		isReady, remaining, err = database.CooldownStars(userID, "rob", cooldownDuration)
	}

	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error al comprobar el cooldown.")
		return err
	}

	if !isReady {
		_, err = ctx.ReplyError("Cooldown", fmt.Sprintf("❌ Tienes que esperar para planear tu próximo golpe. Vuelve en **%d minutos**.", int(remaining.Minutes())))
		return err
	}

	success := rand.Float64() < 0.40 // 40% base success rate

	if ecoType == "local" {
		targetProfile, err := database.GetLocalProfile(guildID, targetUserID)
		if err != nil || targetProfile.Wallet < 100 {
			_, err = ctx.ReplyError("Error", "❌ Ese usuario no tiene dinero que valga la pena robar (Mínimo 100).")
			return err
		}
		myProfile, _ := database.GetLocalProfile(guildID, userID)
		if myProfile.Wallet < 100 {
			_, err = ctx.ReplyError("Error", "❌ Necesitas al menos 100 monedas locales en tu cartera para cubrir posibles fianzas.")
			return err
		}

		_ = database.SetCooldownLocal(guildID, userID, "rob")

		if success {
			percent := (rand.Float64() * 0.3) + 0.1
			stolen := int64(float64(targetProfile.Wallet) * percent)

			database.AddLocalBalance(guildID, targetUserID, -stolen, false)
			database.AddLocalBalance(guildID, userID, stolen, false)
			_, err = ctx.ReplySuccess("¡Robo Exitoso!", fmt.Sprintf("🦹 ¡Éxito! Le robaste **💵 %d monedas** a <@%s>.", stolen, targetUserID))
			return err
		} else {
			fine := int64(float64(myProfile.Wallet) * 0.25)
			if fine < 10 {
				fine = 10
			}
			database.AddLocalBalance(guildID, userID, -fine, false)
			database.AddLocalBalance(guildID, targetUserID, fine, false)
			_, err = ctx.ReplyError("¡Atrapado!", fmt.Sprintf("🚔 ¡Te atraparon intentando robarle a <@%s>! Tuviste que pagarle **💵 %d monedas** como multa.", targetUserID, fine))
			return err
		}
	} else {
		targetProfile, err := database.GetGlobalProfile(targetUserID)
		if err != nil || targetProfile.StarsWallet < 100 {
			_, err = ctx.ReplyError("Error", "❌ Ese usuario no tiene estrellas suficientes en la cartera (Mínimo 100).")
			return err
		}
		myProfile, _ := database.GetGlobalProfile(userID)
		if myProfile.StarsWallet < 100 {
			_, err = ctx.ReplyError("Error", "❌ Necesitas al menos 100 estrellas en tu cartera para cubrir posibles fianzas.")
			return err
		}

		_ = database.SetCooldownStars(userID, "rob")

		if success {
			percent := (rand.Float64() * 0.3) + 0.1
			stolen := int64(float64(targetProfile.StarsWallet) * percent)

			database.AddStars(targetUserID, -stolen, false)
			database.AddStars(userID, stolen, false)
			_, err = ctx.ReplySuccess("¡Robo Exitoso!", fmt.Sprintf("🦹 ¡Éxito! Le robaste **🌟 %d estrellas** a <@%s>.", stolen, targetUserID))
			return err
		} else {
			fine := int64(float64(myProfile.StarsWallet) * 0.25)
			if fine < 10 {
				fine = 10
			}
			database.AddStars(userID, -fine, false)
			database.AddStars(targetUserID, fine, false)
			_, err = ctx.ReplyError("¡Atrapado!", fmt.Sprintf("🚔 ¡Te atraparon robando estrellas de <@%s>! Fuiste multado por **🌟 %d estrellas**, las cuales se le entregaron a tu víctima.", targetUserID, fine))
			return err
		}
	}
}
