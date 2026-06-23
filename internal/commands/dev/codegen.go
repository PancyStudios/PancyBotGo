package dev

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
)

// CreateCodeGenCommand creates the /dev codegen command
func CreateCodeGenCommand() *discord.Command {
	return discord.NewCommand(
		"codegen",
		"✨ | Genera códigos premium (Solo desarrolladores)",
		"dev",
		codegenHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tipo",
			Description: "💻 | Tipo de código premium",
			Required:    true,
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
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "duracion",
			Description: "💻 | Duración en días (0 para permanente)",
			Required:    false,
			MinValue:    float64Ptr(0),
			MaxValue:    3650, // ~10 años
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "cantidad",
			Description: "💻 | Cantidad de códigos a generar (1-10)",
			Required:    false,
			MinValue:    float64Ptr(1),
			MaxValue:    10,
		},
	)
}

func codegenHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()

		// Obtener el ID del usuario de manera segura
		userID := ""
		if ctx.Interaction.Member != nil && ctx.Interaction.Member.User != nil {
			userID = ctx.Interaction.Member.User.ID
		} else if ctx.Interaction.User != nil {
			userID = ctx.Interaction.User.ID
		}
		if userID != "852683369899622430" {
			sendErrorEmbed(ctx, "Acceso Denegado", "❌ Este comando solo está disponible para desarrolladores.")
			return
		}

		// Obtener opciones
		codeType := ctx.GetStringOption("tipo")
		duration := ctx.GetIntOption("duracion")
		cantidad := ctx.GetIntOption("cantidad")

		permanent := false
		if duration == 0 {
			permanent = true
		}

		// Generar códigos
		var generatedCodes []string
		var failedCodes []string

		for i := 0; i < int(cantidad); i++ {
			code := generateRandomCode()

			var premiumType models.PremiumCodeType
			if codeType == "user" {
				premiumType = models.PremiumCodeTypeUser
			} else {
				premiumType = models.PremiumCodeTypeGuild
			}

			_, err := database.CreatePremiumCode(
				code,
				premiumType,
				int(duration),
				permanent,
				userID,
			)

			if err != nil {
				logger.Error(fmt.Sprintf("Error generando código: %v", err), "DevCodeGen")
				failedCodes = append(failedCodes, code)
			} else {
				generatedCodes = append(generatedCodes, code)
			}
		}

		// Crear embed de respuesta
		embed := &discordgo.MessageEmbed{
			Title: "🎫 Códigos Premium Generados",
			Color: 0x00FF00, // Verde
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Tipo",
					Value:  getTipeName(codeType),
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Generado por %s", getUserName(ctx)),
			},
		}

		if permanent {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Duración",
				Value:  "⭐ Permanente",
				Inline: true,
			})
		} else {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Duración",
				Value:  fmt.Sprintf("%d días", duration),
				Inline: true,
			})
		}

		if len(generatedCodes) > 0 {
			codesText := ""
			for _, code := range generatedCodes {
				codesText += fmt.Sprintf("`%s`\n", code)
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("✅ Códigos Generados (%d)", len(generatedCodes)),
				Value:  codesText,
				Inline: false,
			})
		}

		if len(failedCodes) > 0 {
			embed.Color = 0xFFA500 // Naranja
			failedText := ""
			for _, code := range failedCodes {
				failedText += fmt.Sprintf("`%s`\n", code)
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("❌ Códigos Fallidos (%d)", len(failedCodes)),
				Value:  failedText,
				Inline: false,
			})
		}

		// Enviar respuesta efímera (solo visible para el usuario que ejecutó el comando)
		err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Flags:  discordgo.MessageFlagsEphemeral, // Mensaje efímero
			},
		})

		if err != nil {
			logger.Error(fmt.Sprintf("Error enviando respuesta: %v", err), "DevCodeGen")
		}

		logger.Info(fmt.Sprintf("Usuario %s generó %d códigos premium de tipo %s",
			getUserName(ctx), len(generatedCodes), codeType), "DevCodeGen")
	}()

	return nil
}

// generateRandomCode genera un código aleatorio de 16 caracteres
func generateRandomCode() string {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback a un código basado en tiempo si falla
		return fmt.Sprintf("PANC-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("PANC-%s", strings.ToUpper(hex.EncodeToString(bytes)[:12]))
}

// getTipeName devuelve el nombre legible del tipo
func getTipeName(tipo string) string {
	if tipo == "user" {
		return "👤 Usuario"
	}
	return "🏰 Servidor"
}

// float64Ptr helper para convertir float64 a puntero
func float64Ptr(f float64) *float64 {
	return &f
}

// getUserName obtiene el nombre del usuario de manera segura
func getUserName(ctx *discord.CommandContext) string {
	if ctx.Interaction.Member != nil && ctx.Interaction.Member.User != nil {
		return ctx.Interaction.Member.User.Username
	}
	if ctx.Interaction.User != nil {
		return ctx.Interaction.User.Username
	}
	return "Unknown"
}

// sendErrorEmbed envía un embed de error
func sendErrorEmbed(ctx *discord.CommandContext, title, description string) {
	embed := &discordgo.MessageEmbed{
		Title:       "❌ " + title,
		Description: description,
		Color:       0xFF0000, // Rojo
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Error enviando embed de error: %v", err), "DevCodeGen")
	}
}
