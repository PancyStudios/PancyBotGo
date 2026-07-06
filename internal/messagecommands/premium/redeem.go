package premium

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
)

func redeemCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!redeem <codigo>`")
		return err
	}
	code := ctx.Args[0]

	premiumCode, err := database.GetPremiumCode(code)
	if err != nil {
		if err == database.ErrCodeNotFound {
			_, err = ctx.ReplyError("Código inválido", "El código proporcionado no existe o es inválido.")
			return err
		}
		_, err = ctx.ReplyError("Error", "Hubo un error al verificar el código.")
		return err
	}

	if premiumCode.Type == models.PremiumCodeTypeUser {
		return handleUserRedeem(ctx, premiumCode)
	} else if premiumCode.Type == models.PremiumCodeTypeGuild {
		return handleGuildRedeem(ctx, premiumCode)
	}

	_, err = ctx.ReplyError("Error", "Tipo de código desconocido.")
	return err
}

func handleUserRedeem(ctx *messagecommands.MessageContext, premiumCode *models.PremiumCode) error {
	userID := ctx.Message.Author.ID

	isPremium, existingPremium, err := database.IsUserPremium(userID)
	if err != nil {
		_, err = ctx.ReplyError("Error", "Error al verificar tu estado premium.")
		return err
	}

	if isPremium && existingPremium != nil {
		if existingPremium.Permanent {
			_, err = ctx.ReplyError("Ya eres Premium", "Ya tienes una suscripción Premium permanente. ¡Disfruta de tus beneficios sin preocuparte!")
			return err
		}
	}

	_, err = database.RedeemPremiumCode(premiumCode.Code, userID)
	if err != nil {
		if err == database.ErrCodeAlreadyClaimed {
			_, err = ctx.ReplyError("Código usado", "Este código ya ha sido canjeado por alguien más.")
			return err
		}
		_, err = ctx.ReplyError("Error", "Hubo un problema al canjear el código.")
		return err
	}

	duracionMsg := fmt.Sprintf("por %d días", premiumCode.DurationDays)
	if premiumCode.Permanent {
		duracionMsg = "permanentemente"
	}

	_, err = ctx.ReplySuccess("¡Gracias por apoyar el bot!", fmt.Sprintf("✅ **Has canjeado exitosamente un código Premium.**\n\nDisfruta de tus beneficios %s. Para ver tus beneficios, puedes usar el comando `/premium info` o ver en la página web.", duracionMsg))
	return err
}

func handleGuildRedeem(ctx *messagecommands.MessageContext, premiumCode *models.PremiumCode) error {
	guildID := ctx.Message.GuildID
	if guildID == "" {
		_, err := ctx.ReplyError("Error", "Los códigos de servidor solo se pueden canjear dentro de un servidor.")
		return err
	}

	isPremium, existingPremium, err := database.IsGuildPremium(guildID)
	if err != nil {
		_, err = ctx.ReplyError("Error", "Error al verificar el estado premium del servidor.")
		return err
	}

	if isPremium && existingPremium != nil {
		if existingPremium.Permanent {
			_, err = ctx.ReplyError("Servidor ya Premium", "Este servidor ya tiene una suscripción Premium permanente.")
			return err
		}
	}

	_, err = database.RedeemPremiumCodeForGuild(premiumCode.Code, guildID, ctx.Message.Author.ID)
	if err != nil {
		if err == database.ErrCodeAlreadyClaimed {
			_, err = ctx.ReplyError("Código usado", "Este código ya ha sido canjeado por alguien más.")
			return err
		}
		_, err = ctx.ReplyError("Error", "Hubo un problema al canjear el código para este servidor.")
		return err
	}

	duracionMsg := fmt.Sprintf("por %d días", premiumCode.DurationDays)
	if premiumCode.Permanent {
		duracionMsg = "permanentemente"
	}

	_, err = ctx.ReplySuccess("¡Gracias por apoyar el bot!", fmt.Sprintf("✅ **Has canjeado exitosamente un código Premium para este servidor.**\n\nEl servidor disfrutará de sus beneficios %s.", duracionMsg))
	return err
}
