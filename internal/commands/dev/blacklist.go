package dev

import (
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

// CreateBlacklistAddCommand crea el comando /dev blacklist add
func CreateBlacklistAddCommand() *discord.Command {
	return discord.NewCommand(
		"add",
		"Añade un usuario o servidor a la blacklist",
		"dev",
		blacklistAddHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tipo",
			Description: "Tipo de blacklist",
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
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "id",
			Description: "ID del usuario o servidor",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "razon",
			Description: "Razón del blacklist",
			Required:    false,
		},
	)
}

// CreateBlacklistRemoveCommand crea el comando /dev blacklist remove
func CreateBlacklistRemoveCommand() *discord.Command {
	return discord.NewCommand(
		"remove",
		"Elimina un usuario o servidor de la blacklist",
		"dev",
		blacklistRemoveHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "id",
			Description:  "ID del usuario o servidor",
			Required:     true,
			Autocomplete: true,
		},
	).WithAutoComplete(blacklistRemoveAutoComplete)
}

// CreateBlacklistListCommand crea el comando /dev blacklist list
func CreateBlacklistListCommand() *discord.Command {
	return discord.NewCommand(
		"list",
		"Lista todas las entradas de la blacklist",
		"dev",
		blacklistListHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tipo",
			Description: "Filtrar por tipo",
			Required:    false,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Todos",
					Value: "all",
				},
				{
					Name:  "Usuarios",
					Value: "user",
				},
				{
					Name:  "Servidores",
					Value: "guild",
				},
			},
		},
	)
}

func blacklistAddHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()

		// Verificar permisos de desarrollador
		if !isDev(ctx.User().ID) {
			sendErrorEmbed(ctx, "Acceso Denegado", "❌ Este comando solo está disponible para desarrolladores.")
			return
		}

		// Obtener opciones
		tipo := ctx.GetStringOption("tipo")
		id := ctx.GetStringOption("id")
		razon := ctx.GetStringOption("razon")

		if razon == "" {
			razon = "No especificada"
		}

		// Validar tipo
		var blacklistType models.BlacklistType
		if tipo == "user" {
			blacklistType = models.BlacklistTypeUser
		} else if tipo == "guild" {
			blacklistType = models.BlacklistTypeGuild
		} else {
			sendErrorEmbed(ctx, "Error", "❌ Tipo inválido. Usa 'user' o 'guild'.")
			return
		}

		// Añadir a la blacklist
		_, err := database.AddToBlacklist(id, blacklistType, razon, ctx.User().ID)
		if err != nil {
			if err == database.ErrAlreadyBlacklisted {
				sendErrorEmbed(ctx, "Error", "❌ Este ID ya está en la blacklist.")
				return
			}
			logger.Error(fmt.Sprintf("Error añadiendo a blacklist: %v", err), "DevBlacklist")
			sendErrorEmbed(ctx, "Error", "❌ Error al añadir a la blacklist.")
			return
		}

		// Si es un guild y el bot está en él, salir
		if blacklistType == models.BlacklistTypeGuild {
			guild, err := ctx.Session.Guild(id)
			if err == nil && guild != nil {
				// Intentar enviar mensaje al owner
				owner, err := ctx.Session.User(guild.OwnerID)
				if err == nil && owner != nil {
					channel, err := ctx.Session.UserChannelCreate(owner.ID)
					if err == nil {
						embed := &discordgo.MessageEmbed{
							Title:       "🚫 Servidor en Blacklist",
							Description: fmt.Sprintf("El servidor **%s** ha sido añadido a la blacklist y el bot se retirará automáticamente.", guild.Name),
							Color:       0xFF0000,
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Razón",
									Value: razon,
								},
							},
							Timestamp: time.Now().Format(time.RFC3339),
						}
						ctx.Session.ChannelMessageSendEmbed(channel.ID, embed)
					}
				}

				// Salir del servidor
				ctx.Session.GuildLeave(id)
				logger.Info(fmt.Sprintf("Bot salió del servidor blacklisted: %s (%s)", guild.Name, id), "DevBlacklist")
			}
		}

		// Crear embed de confirmación
		embed := &discordgo.MessageEmbed{
			Title:       "✅ Añadido a la Blacklist",
			Description: fmt.Sprintf("El %s con ID `%s` ha sido añadido a la blacklist.", getTypeName(blacklistType), id),
			Color:       0x00FF00,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Tipo",
					Value:  getTypeName(blacklistType),
					Inline: true,
				},
				{
					Name:   "ID",
					Value:  fmt.Sprintf("`%s`", id),
					Inline: true,
				},
				{
					Name:   "Razón",
					Value:  razon,
					Inline: false,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Añadido por %s", getUserName(ctx)),
			},
		}

		ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})

		logger.Info(fmt.Sprintf("Usuario %s añadió a blacklist: %s (%s)", getUserName(ctx), id, tipo), "DevBlacklist")
	}()

	return nil
}

func blacklistRemoveHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()

		// Verificar permisos de desarrollador
		if !isDev(ctx.User().ID) {
			sendErrorEmbed(ctx, "Acceso Denegado", "❌ Este comando solo está disponible para desarrolladores.")
			return
		}

		// Obtener ID
		id := ctx.GetStringOption("id")
		if id == "" {
			sendErrorEmbed(ctx, "Error", "❌ Debes especificar un ID válido.")
			return
		}

		// Obtener entrada antes de eliminar
		entry, err := database.GetBlacklistEntry(id)
		if err != nil {
			sendErrorEmbed(ctx, "Error", "❌ Este ID no está en la blacklist.")
			return
		}

		// Eliminar de la blacklist
		err = database.RemoveFromBlacklist(id)
		if err != nil {
			logger.Error(fmt.Sprintf("Error eliminando de blacklist: %v", err), "DevBlacklist")
			sendErrorEmbed(ctx, "Error", "❌ Error al eliminar de la blacklist.")
			return
		}

		// Crear embed de confirmación
		embed := &discordgo.MessageEmbed{
			Title:       "✅ Eliminado de la Blacklist",
			Description: fmt.Sprintf("El %s con ID `%s` ha sido eliminado de la blacklist.", getTypeName(entry.Type), id),
			Color:       0x00FF00,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Tipo",
					Value:  getTypeName(entry.Type),
					Inline: true,
				},
				{
					Name:   "ID",
					Value:  fmt.Sprintf("`%s`", id),
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Eliminado por %s", getUserName(ctx)),
			},
		}

		if entry.Reason != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Razón Original",
				Value:  entry.Reason,
				Inline: false,
			})
		}

		ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})

		logger.Info(fmt.Sprintf("Usuario %s eliminó de blacklist: %s", getUserName(ctx), id), "DevBlacklist")
	}()

	return nil
}

func blacklistListHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()

		// Verificar permisos de desarrollador
		if !isDev(ctx.User().ID) {
			sendErrorEmbed(ctx, "Acceso Denegado", "❌ Este comando solo está disponible para desarrolladores.")
			return
		}

		// Obtener filtro
		filtro := "all"
		if ctx.HasOption("tipo") {
			filtro = ctx.GetStringOption("tipo")
		}

		// Obtener entradas
		var entries []*models.Blacklist
		var err error

		if filtro == "all" {
			entries, err = database.GetAllBlacklist()
		} else {
			var blacklistType models.BlacklistType
			if filtro == "user" {
				blacklistType = models.BlacklistTypeUser
			} else {
				blacklistType = models.BlacklistTypeGuild
			}
			entries, err = database.GetBlacklistByType(blacklistType)
		}

		if err != nil {
			logger.Error(fmt.Sprintf("Error obteniendo blacklist: %v", err), "DevBlacklist")
			sendErrorEmbed(ctx, "Error", "❌ Error al obtener la blacklist.")
			return
		}

		// Si no hay entradas
		if len(entries) == 0 {
			embed := &discordgo.MessageEmbed{
				Title:       "📋 Blacklist",
				Description: "No hay entradas en la blacklist con los filtros especificados.",
				Color:       0xFFFF00,
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

		// Crear embeds (máximo 10 entradas por embed)
		var embeds []*discordgo.MessageEmbed
		const entriesPerEmbed = 10

		for i := 0; i < len(entries); i += entriesPerEmbed {
			end := i + entriesPerEmbed
			if end > len(entries) {
				end = len(entries)
			}

			embed := &discordgo.MessageEmbed{
				Title:     fmt.Sprintf("📋 Blacklist (%d-%d de %d)", i+1, end, len(entries)),
				Color:     0xFF0000,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			if filtro != "all" {
				embed.Description = fmt.Sprintf("Filtrando por: %s", getFilterName(filtro))
			}

			for _, entry := range entries[i:end] {
				fieldValue := formatBlacklistEntry(entry)
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   fmt.Sprintf("%s `%s`", getTypeEmoji(entry.Type), entry.ID),
					Value:  fieldValue,
					Inline: false,
				})
			}

			embeds = append(embeds, embed)
		}

		// Enviar primera respuesta
		ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embeds[0]},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})

		// Enviar follow-ups si hay más embeds
		for i := 1; i < len(embeds); i++ {
			ctx.Session.FollowupMessageCreate(ctx.Interaction.Interaction, true, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{embeds[i]},
				Flags:  discordgo.MessageFlagsEphemeral,
			})
			time.Sleep(100 * time.Millisecond)
		}

		logger.Info(fmt.Sprintf("Usuario %s listó la blacklist (%d entradas)", getUserName(ctx), len(entries)), "DevBlacklist")
	}()

	return nil
}

// blacklistRemoveAutoComplete maneja el autocompletado para blacklist remove
func blacklistRemoveAutoComplete(ctx *discord.CommandContext) {
	data := ctx.Interaction.ApplicationCommandData()

	// Obtener el valor actual
	var focusedValue string
	for _, opt := range data.Options {
		if opt.Focused {
			focusedValue = opt.StringValue()
			break
		}
	}

	// Obtener todas las entradas
	entries, err := database.GetAllBlacklist()
	if err != nil {
		logger.Error(fmt.Sprintf("Error obteniendo blacklist para autocompletado: %v", err), "DevBlacklist")
		ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
		return
	}

	// Filtrar y crear opciones
	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, entry := range entries {
		if len(choices) >= 25 {
			break
		}

		// Filtrar por búsqueda
		if focusedValue == "" || strings.Contains(strings.ToLower(entry.ID), strings.ToLower(focusedValue)) {
			name := fmt.Sprintf("%s %s", getTypeEmoji(entry.Type), entry.ID)
			if len(name) > 100 {
				name = name[:97] + "..."
			}

			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  name,
				Value: entry.ID,
			})
		}
	}

	// Si no hay resultados
	if len(choices) == 0 {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  "No se encontraron entradas",
			Value: "none",
		})
	}

	ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

// Helper functions

func getTypeName(blacklistType models.BlacklistType) string {
	if blacklistType == models.BlacklistTypeUser {
		return "usuario"
	}
	return "servidor"
}

func getTypeEmoji(blacklistType models.BlacklistType) string {
	if blacklistType == models.BlacklistTypeUser {
		return "👤"
	}
	return "🏰"
}

func getFilterName(filtro string) string {
	switch filtro {
	case "user":
		return "Usuarios"
	case "guild":
		return "Servidores"
	default:
		return "Todos"
	}
}

func formatBlacklistEntry(entry *models.Blacklist) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("**Tipo:** %s", getTypeName(entry.Type)))

	if entry.Reason != "" {
		parts = append(parts, fmt.Sprintf("**Razón:** %s", entry.Reason))
	}

	if entry.AddedBy != "" {
		parts = append(parts, fmt.Sprintf("**Añadido por:** <@%s>", entry.AddedBy))
	}

	if !entry.AddedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("**Fecha:** <t:%d:R>", entry.AddedAt.Unix()))
	}

	return strings.Join(parts, "\n")
}
