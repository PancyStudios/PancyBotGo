package dev

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// CreateCodeDelCommand creates the /dev codedel command
func CreateCodeDelCommand() *discord.Command {
	return discord.NewCommand(
		"codedel",
		"✨ | Elimina un código premium (Solo desarrolladores)",
		"dev",
		codedelHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "codigo",
			Description: "💻 | Código premium a eliminar",
			Required:     true,
			Autocomplete: true,
		},
	).WithAutoComplete(codedelAutoComplete)
}

func codedelHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()

		// Verificar permisos de desarrollador
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

		// Obtener el código a eliminar
		code := ctx.GetStringOption("codigo")
		if code == "" {
			sendErrorEmbed(ctx, "Error", "❌ Debes especificar un código válido.")
			return
		}

		// Verificar que el código existe antes de eliminarlo
		codeData, err := database.GetPremiumCode(code)
		if err != nil {
			logger.Error(fmt.Sprintf("Error obteniendo código %s: %v", code, err), "DevCodeDel")
			sendErrorEmbed(ctx, "Error", fmt.Sprintf("❌ El código `%s` no existe.", code))
			return
		}

		// Eliminar el código
		err = database.DeletePremiumCode(code)
		if err != nil {
			logger.Error(fmt.Sprintf("Error eliminando código %s: %v", code, err), "DevCodeDel")
			sendErrorEmbed(ctx, "Error", fmt.Sprintf("❌ Error al eliminar el código `%s`.", code))
			return
		}

		// Crear embed de confirmación
		embed := &discordgo.MessageEmbed{
			Title:       "🗑️ Código Premium Eliminado",
			Description: fmt.Sprintf("El código `%s` ha sido eliminado correctamente.", code),
			Color:       0xFF0000, // Rojo
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Tipo",
					Value:  getCodeTypeName(string(codeData.Type)),
					Inline: true,
				},
				{
					Name:   "Estado",
					Value:  getCodeStatus(codeData.IsClaimed),
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Eliminado por %s", getUserName(ctx)),
			},
		}

		// Agregar información de duración
		if codeData.Permanent {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Duración",
				Value:  "⭐ Permanente",
				Inline: true,
			})
		} else {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Duración",
				Value:  fmt.Sprintf("%d días", codeData.DurationDays),
				Inline: true,
			})
		}

		// Agregar información de canje si fue canjeado
		if codeData.IsClaimed && codeData.ClaimedBy != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Canjeado por",
				Value:  fmt.Sprintf("<@%s>", codeData.ClaimedBy),
				Inline: true,
			})

			if !codeData.ClaimedAt.IsZero() {
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   "Canjeado el",
					Value:  fmt.Sprintf("<t:%d:R>", codeData.ClaimedAt.Unix()),
					Inline: true,
				})
			}
		}

		// Enviar respuesta
		err = ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})

		if err != nil {
			logger.Error(fmt.Sprintf("Error enviando respuesta: %v", err), "DevCodeDel")
			return
		}

		logger.Info(fmt.Sprintf("Usuario %s eliminó el código premium %s", getUserName(ctx), code), "DevCodeDel")
	}()

	return nil
}

// codedelAutoComplete maneja el autocompletado para el comando codedel
func codedelAutoComplete(ctx *discord.CommandContext) {
	data := ctx.Interaction.ApplicationCommandData()

	// Obtener el valor actual que el usuario está escribiendo
	var focusedValue string
	for _, opt := range data.Options {
		if opt.Focused {
			focusedValue = opt.StringValue()
			break
		}
	}

	// Obtener todos los códigos disponibles
	codes, err := database.GetAllPremiumCodes()
	if err != nil {
		logger.Error(fmt.Sprintf("Error obteniendo códigos para autocompletado: %v", err), "DevCodeDel")
		ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
		return
	}

	// Filtrar códigos que coincidan con lo que el usuario está escribiendo
	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, code := range codes {
		// Limitar a 25 resultados (límite de Discord)
		if len(choices) >= 25 {
			break
		}

		// Filtrar por el valor que el usuario está escribiendo
		if focusedValue == "" || containsIgnoreCase(code.Code, focusedValue) {
			status := "🎫"
			if code.IsClaimed {
				status = "✅"
			}

			typeIcon := "👤"
			if code.Type == "guild" {
				typeIcon = "🏰"
			}

			name := fmt.Sprintf("%s %s %s", status, typeIcon, code.Code)

			// Truncar si es muy largo
			if len(name) > 100 {
				name = name[:97] + "..."
			}

			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  name,
				Value: code.Code,
			})
		}
	}

	// Si no hay resultados, mostrar un mensaje
	if len(choices) == 0 {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  "No se encontraron códigos",
			Value: "none",
		})
	}

	// Enviar las opciones de autocompletado
	err = ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Error enviando autocompletado: %v", err), "DevCodeDel")
	}
}

// getCodeTypeName devuelve el nombre del tipo de código con emoji
func getCodeTypeName(codeType string) string {
	if codeType == "user" {
		return "👤 Usuario"
	}
	return "🏰 Servidor"
}

// getCodeStatus devuelve el estado del código con emoji
func getCodeStatus(isClaimed bool) string {
	if isClaimed {
		return "✅ Canjeado"
	}
	return "🎫 Disponible"
}

// containsIgnoreCase verifica si str contiene substr ignorando mayúsculas/minúsculas
func containsIgnoreCase(str, substr string) bool {
	if substr == "" {
		return true
	}

	strLower := ""
	substrLower := ""

	for _, r := range str {
		if r >= 'A' && r <= 'Z' {
			strLower += string(r + 32)
		} else {
			strLower += string(r)
		}
	}

	for _, r := range substr {
		if r >= 'A' && r <= 'Z' {
			substrLower += string(r + 32)
		} else {
			substrLower += string(r)
		}
	}

	// Implementación simple de contains
	if len(substrLower) > len(strLower) {
		return false
	}

	for i := 0; i <= len(strLower)-len(substrLower); i++ {
		if strLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}

	return false
}
