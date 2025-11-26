package premium

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
)

// CreateRedeemCommand creates the /premium redeem command
func CreateRedeemCommand() *discord.Command {
	return discord.NewCommand(
		"redeem",
		"Canjea un código premium",
		"premium",
		redeemHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "codigo",
			Description: "El código premium a canjear",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tipo",
			Description: "Tipo de código (user/guild)",
			Required:    false,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Usuario",
					Value: "user",
				},
				{
					Name:  "Servidor",
					Value: "guild",
				},
			},
		},
	)
}

func redeemHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()

		code := ctx.GetStringOption("codigo")
		codeType := ctx.GetStringOption("tipo")

		if codeType == "" {
			premiumCode, err := database.GetPremiumCode(code)
			if err != nil {
				if err == database.ErrCodeNotFound {
					sendErrorEmbed(ctx, "Código inválido", "El código proporcionado no existe o es inválido.")
					return
				}
				logger.Error(fmt.Sprintf("Error obteniendo código: %v", err), "Premium")
				sendErrorEmbed(ctx, "Error", "Hubo un error al verificar el código.")
				return
			}
			codeType = string(premiumCode.Type)
		}

		if codeType == "user" {
			handleUserRedeem(ctx, code)
		} else if codeType == "guild" {
			handleGuildRedeem(ctx, code)
		} else {
			sendErrorEmbed(ctx, "Tipo inválido", "El tipo de código debe ser 'user' o 'guild'.")
		}
	}()
	return nil
}

func handleUserRedeem(ctx *discord.CommandContext, code string) {
	premiumCode, err := database.GetPremiumCode(code)
	if err != nil {
		if err == database.ErrCodeNotFound {
			sendErrorEmbed(ctx, "Código inválido", "El código proporcionado no existe o es inválido.")
			return
		}
		logger.Error(fmt.Sprintf("Error obteniendo código: %v", err), "Premium")
		sendErrorEmbed(ctx, "Error", "Hubo un error al verificar el código.")
		return
	}

	if premiumCode.Type != models.PremiumCodeTypeUser {
		sendErrorEmbed(ctx, "Tipo incorrecto", "Este código es para servidores, no para usuarios.")
		return
	}

	userID := getUserID(ctx)

	isPremium, existingPremium, err := database.IsUserPremium(userID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error verificando premium: %v", err), "Premium")
		sendErrorEmbed(ctx, "Error", "Hubo un error al verificar tu estado premium.")
		return
	}

	if isPremium {
		var expiresText string
		if existingPremium.Permanent {
			expiresText = "permanente"
		} else {
			expiresAt := time.UnixMilli(existingPremium.ExpiresAt)
			expiresText = fmt.Sprintf("<t:%d:R>", expiresAt.Unix())
		}
		sendErrorEmbed(ctx, "Ya eres premium", fmt.Sprintf("Ya tienes premium activo que expira %s.", expiresText))
		return
	}

	redeemedCode, err := database.RedeemPremiumCode(code, userID)
	if err != nil {
		if err == database.ErrCodeAlreadyClaimed {
			sendErrorEmbed(ctx, "Código ya usado", "Este código ya ha sido canjeado por otro usuario.")
			return
		}
		logger.Error(fmt.Sprintf("Error canjeando código: %v", err), "Premium")
		sendErrorEmbed(ctx, "Error", "Hubo un error al canjear el código.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "✨ Premium Activado",
		Description: "¡Has canjeado exitosamente tu código premium!",
		Color:       0xFFD700,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Tipo",
				Value:  "Premium de Usuario",
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "PancyBot Premium",
		},
	}

	if redeemedCode.Permanent {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Duración",
			Value:  "⭐ Permanente",
			Inline: true,
		})
	} else {
		duration := time.Duration(redeemedCode.DurationDays) * 24 * time.Hour
		expiresAt := time.Now().Add(duration)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Duración",
			Value:  fmt.Sprintf("%d días", redeemedCode.DurationDays),
			Inline: true,
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Expira",
			Value:  fmt.Sprintf("<t:%d:F>", expiresAt.Unix()),
			Inline: false,
		})
	}

	sendSuccessEmbed(ctx, embed)

	logger.Info(fmt.Sprintf("Usuario %s (%s) canjeó código premium: %s", getUserName(ctx), userID, code), "Premium")
}

func handleGuildRedeem(ctx *discord.CommandContext, code string) {
	premiumCode, err := database.GetPremiumCode(code)
	if err != nil {
		if err == database.ErrCodeNotFound {
			sendErrorEmbed(ctx, "Código inválido", "El código proporcionado no existe o es inválido.")
			return
		}
		logger.Error(fmt.Sprintf("Error obteniendo código: %v", err), "Premium")
		sendErrorEmbed(ctx, "Error", "Hubo un error al verificar el código.")
		return
	}

	if premiumCode.Type != models.PremiumCodeTypeGuild {
		sendErrorEmbed(ctx, "Tipo incorrecto", "Este código es para usuarios, no para servidores.")
		return
	}

	guildID := ctx.Interaction.GuildID
	if guildID == "" {
		sendErrorEmbed(ctx, "Solo en servidores", "Los códigos de servidor solo pueden canjearse dentro de un servidor.")
		return
	}

	userID := getUserID(ctx)

	member, err := ctx.Session.GuildMember(guildID, userID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error obteniendo miembro: %v", err), "Premium")
		sendErrorEmbed(ctx, "Error", "Hubo un error al verificar tus permisos.")
		return
	}

	guild, err := ctx.Session.Guild(guildID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error obteniendo guild: %v", err), "Premium")
		sendErrorEmbed(ctx, "Error", "Hubo un error al verificar el servidor.")
		return
	}

	hasPermission := false
	for _, roleID := range member.Roles {
		role, err := ctx.Session.State.Role(guildID, roleID)
		if err == nil && (role.Permissions&discordgo.PermissionAdministrator != 0) {
			hasPermission = true
			break
		}
	}

	if !hasPermission && member.User.ID != guild.OwnerID {
		sendErrorEmbed(ctx, "Sin permisos", "Solo los administradores pueden canjear códigos premium para el servidor.")
		return
	}

	isPremium, existingPremium, err := database.IsGuildPremium(guildID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error verificando premium del servidor: %v", err), "Premium")
		sendErrorEmbed(ctx, "Error", "Hubo un error al verificar el estado premium del servidor.")
		return
	}

	if isPremium {
		var expiresText string
		if existingPremium.Permanent {
			expiresText = "permanente"
		} else {
			expiresAt := time.UnixMilli(existingPremium.ExpiresAt)
			expiresText = fmt.Sprintf("<t:%d:R>", expiresAt.Unix())
		}
		sendErrorEmbed(ctx, "Servidor ya premium", fmt.Sprintf("Este servidor ya tiene premium activo que expira %s.", expiresText))
		return
	}

	redeemedCode, err := database.RedeemPremiumCodeForGuild(code, guildID, userID)
	if err != nil {
		if err == database.ErrCodeAlreadyClaimed {
			sendErrorEmbed(ctx, "Código ya usado", "Este código ya ha sido canjeado.")
			return
		}
		logger.Error(fmt.Sprintf("Error canjeando código: %v", err), "Premium")
		sendErrorEmbed(ctx, "Error", "Hubo un error al canjear el código.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "✨ Premium Activado",
		Description: fmt.Sprintf("¡El servidor **%s** ahora tiene premium activo!", guild.Name),
		Color:       0xFFD700,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Tipo",
				Value:  "Premium de Servidor",
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "PancyBot Premium",
		},
	}

	if redeemedCode.Permanent {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Duración",
			Value:  "⭐ Permanente",
			Inline: true,
		})
	} else {
		duration := time.Duration(redeemedCode.DurationDays) * 24 * time.Hour
		expiresAt := time.Now().Add(duration)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Duración",
			Value:  fmt.Sprintf("%d días", redeemedCode.DurationDays),
			Inline: true,
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Expira",
			Value:  fmt.Sprintf("<t:%d:F>", expiresAt.Unix()),
			Inline: false,
		})
	}

	sendSuccessEmbed(ctx, embed)

	logger.Info(fmt.Sprintf("Usuario %s (%s) canjeó código premium de servidor para guild %s (%s): %s",
		getUserName(ctx), userID, guild.Name, guildID, code), "Premium")
}

func sendErrorEmbed(ctx *discord.CommandContext, title, description string) {
	embed := &discordgo.MessageEmbed{
		Title:       "❌ " + title,
		Description: description,
		Color:       0xFF0000,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Error enviando embed de error: %v", err), "Premium")
	}
}

func sendSuccessEmbed(ctx *discord.CommandContext, embed *discordgo.MessageEmbed) {
	err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Error enviando embed de éxito: %v", err), "Premium")
	}
}

func getUserID(ctx *discord.CommandContext) string {
	if ctx.Interaction.Member != nil && ctx.Interaction.Member.User != nil {
		return ctx.Interaction.Member.User.ID
	}
	if ctx.Interaction.User != nil {
		return ctx.Interaction.User.ID
	}
	return ""
}

func getUserName(ctx *discord.CommandContext) string {
	if ctx.Interaction.Member != nil && ctx.Interaction.Member.User != nil {
		return ctx.Interaction.Member.User.Username
	}
	if ctx.Interaction.User != nil {
		return ctx.Interaction.User.Username
	}
	return "Unknown"
}
