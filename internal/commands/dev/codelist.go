package dev

import (
	"fmt"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// CreateCodeListCommand creates the /dev codelist command
func CreateCodeListCommand() *discord.Command {
	return discord.NewCommand(
		"codelist",
		"Lista todos los códigos premium generados (Solo desarrolladores)",
		"dev",
		codelistHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "filtro",
			Description: "Filtrar códigos por estado",
			Required:    false,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Todos",
					Value: "all",
				},
				{
					Name:  "Disponibles",
					Value: "available",
				},
				{
					Name:  "Canjeados",
					Value: "claimed",
				},
			},
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tipo",
			Description: "Filtrar códigos por tipo",
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

func codelistHandler(ctx *discord.CommandContext) error {
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

		// Obtener opciones
		filtro := "all"
		if ctx.HasOption("filtro") {
			filtro = ctx.GetStringOption("filtro")
		}

		tipoFiltro := ""
		if ctx.HasOption("tipo") {
			tipoFiltro = ctx.GetStringOption("tipo")
		}

		// Obtener códigos de la base de datos
		codes, err := database.GetAllPremiumCodes()
		if err != nil {
			logger.Error(fmt.Sprintf("Error obteniendo códigos: %v", err), "DevCodeList")
			sendErrorEmbed(ctx, "Error", "❌ Error al obtener los códigos premium.")
			return
		}

		// Filtrar códigos
		var filteredCodes []*database.PremiumCode
		for _, code := range codes {
			// Filtro por estado
			if filtro == "available" && code.IsClaimed {
				continue
			}
			if filtro == "claimed" && !code.IsClaimed {
				continue
			}

			// Filtro por tipo
			if tipoFiltro != "" && string(code.Type) != tipoFiltro {
				continue
			}

			filteredCodes = append(filteredCodes, code)
		}

		// Si no hay códigos
		if len(filteredCodes) == 0 {
			embed := &discordgo.MessageEmbed{
				Title:       "📋 Lista de Códigos Premium",
				Description: "No se encontraron códigos con los filtros especificados.",
				Color:       0xFFFF00, // Amarillo
				Timestamp:   time.Now().Format(time.RFC3339),
			}

			ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
					Flags:  discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		// Crear embeds (máximo 10 códigos por embed debido al límite de Discord)
		var embeds []*discordgo.MessageEmbed
		const codesPerEmbed = 10

		for i := 0; i < len(filteredCodes); i += codesPerEmbed {
			end := i + codesPerEmbed
			if end > len(filteredCodes) {
				end = len(filteredCodes)
			}

			embed := &discordgo.MessageEmbed{
				Title: fmt.Sprintf("📋 Lista de Códigos Premium (%d-%d de %d)", i+1, end, len(filteredCodes)),
				Color: 0x00BFFF, // Azul claro
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Filtros Aplicados",
						Value:  getFilterDescription(filtro, tipoFiltro),
						Inline: false,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			for _, code := range filteredCodes[i:end] {
				fieldValue := formatCodeInfo(code)
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   fmt.Sprintf("`%s`", code.Code),
					Value:  fieldValue,
					Inline: false,
				})
			}

			embeds = append(embeds, embed)
		}

		// Enviar respuesta (solo el primer embed por ahora, Discord tiene límites)
		err = ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embeds[0]},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})

		if err != nil {
			logger.Error(fmt.Sprintf("Error enviando respuesta: %v", err), "DevCodeList")
			return
		}

		// Si hay más embeds, enviarlos como follow-up
		for i := 1; i < len(embeds); i++ {
			_, err = ctx.Session.FollowupMessageCreate(ctx.Interaction.Interaction, true, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{embeds[i]},
				Flags:  discordgo.MessageFlagsEphemeral,
			})
			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando follow-up: %v", err), "DevCodeList")
				break
			}
			time.Sleep(100 * time.Millisecond) // Pequeño delay para evitar rate limits
		}

		logger.Info(fmt.Sprintf("Usuario %s listó %d códigos premium", getUserName(ctx), len(filteredCodes)), "DevCodeList")
	}()

	return nil
}

// formatCodeInfo formatea la información de un código
func formatCodeInfo(code *database.PremiumCode) string {
	var parts []string

	// Tipo
	if code.Type == "user" {
		parts = append(parts, "**Tipo:** 👤 Usuario")
	} else {
		parts = append(parts, "**Tipo:** 🏰 Servidor")
	}

	// Estado
	if code.IsClaimed {
		parts = append(parts, "**Estado:** ✅ Canjeado")
		if code.ClaimedBy != "" {
			parts = append(parts, fmt.Sprintf("**Canjeado por:** <@%s>", code.ClaimedBy))
		}
		if !code.ClaimedAt.IsZero() {
			parts = append(parts, fmt.Sprintf("**Canjeado el:** <t:%d:R>", code.ClaimedAt.Unix()))
		}
	} else {
		parts = append(parts, "**Estado:** 🎫 Disponible")
	}

	// Duración
	if code.Permanent {
		parts = append(parts, "**Duración:** ⭐ Permanente")
	} else {
		parts = append(parts, fmt.Sprintf("**Duración:** %d días", code.DurationDays))
	}

	// Fecha de creación
	if !code.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("**Creado:** <t:%d:R>", code.CreatedAt.Unix()))
	}

	// Creado por
	if code.CreatedBy != "" {
		parts = append(parts, fmt.Sprintf("**Creado por:** <@%s>", code.CreatedBy))
	}

	return strings.Join(parts, "\n")
}

// getFilterDescription devuelve una descripción de los filtros aplicados
func getFilterDescription(filtro, tipo string) string {
	var parts []string

	switch filtro {
	case "available":
		parts = append(parts, "🎫 Solo disponibles")
	case "claimed":
		parts = append(parts, "✅ Solo canjeados")
	default:
		parts = append(parts, "📋 Todos los códigos")
	}

	if tipo != "" {
		if tipo == "user" {
			parts = append(parts, "👤 Usuario")
		} else {
			parts = append(parts, "🏰 Servidor")
		}
	}

	if len(parts) == 0 {
		return "Sin filtros"
	}

	return strings.Join(parts, " | ")
}
