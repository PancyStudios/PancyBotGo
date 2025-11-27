package dev

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

// CreateBlacklistAddCommand creates the /dev blacklist add command
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
			Description: "Tipo de entrada a bloquear",
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
			Description: "ID del usuario o servidor a bloquear",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "razon",
			Description: "Razón del bloqueo",
			Required:    false,
		},
	)
}

// CreateBlacklistRemoveCommand creates the /dev blacklist remove command
func CreateBlacklistRemoveCommand() *discord.Command {
	return discord.NewCommand(
		"remove",
		"Elimina un usuario o servidor de la blacklist",
		"dev",
		blacklistRemoveHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tipo",
			Description: "Tipo de entrada a desbloquear",
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
			Description: "ID del usuario o servidor a desbloquear",
			Required:    true,
		},
	)
}

func blacklistAddHandler(ctx *discord.CommandContext) error {
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
		tipo := ctx.GetStringOption("tipo")
		id := ctx.GetStringOption("id")
		razon := ctx.GetStringOption("razon")

		if razon == "" {
			razon = "Sin razón especificada"
		}

		// Determinar el tipo de blacklist
		var blacklistType models.BlacklistType
		if tipo == "user" {
			blacklistType = models.BlacklistTypeUser
		} else {
			blacklistType = models.BlacklistTypeGuild
		}

		// Añadir a la blacklist
		entry, err := database.AddToBlacklist(id, blacklistType, razon, userID)
		if err != nil {
			if err == database.ErrBlacklistEntryExists {
				sendErrorEmbed(ctx, "Error", fmt.Sprintf("❌ El %s `%s` ya está en la blacklist.", getBlacklistTypeName(tipo), id))
				return
			}
			logger.Error(fmt.Sprintf("Error añadiendo a blacklist: %v", err), "DevBlacklist")
			sendErrorEmbed(ctx, "Error", "❌ Error al añadir a la blacklist.")
			return
		}

		// Crear embed de confirmación
		embed := &discordgo.MessageEmbed{
			Title:       "🚫 Añadido a la Blacklist",
			Description: fmt.Sprintf("El %s ha sido bloqueado correctamente.", getBlacklistTypeName(tipo)),
			Color:       0xFF0000, // Rojo
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Tipo",
					Value:  getBlacklistTypeEmoji(tipo) + " " + getBlacklistTypeName(tipo),
					Inline: true,
				},
				{
					Name:   "ID",
					Value:  fmt.Sprintf("`%s`", id),
					Inline: true,
				},
				{
					Name:   "Razón",
					Value:  entry.Reason,
					Inline: false,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Bloqueado por %s", getUserName(ctx)),
			},
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
			logger.Error(fmt.Sprintf("Error enviando respuesta: %v", err), "DevBlacklist")
			return
		}

		logger.Info(fmt.Sprintf("Usuario %s añadió %s %s a la blacklist", getUserName(ctx), tipo, id), "DevBlacklist")
	}()

	return nil
}

func blacklistRemoveHandler(ctx *discord.CommandContext) error {
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
		tipo := ctx.GetStringOption("tipo")
		id := ctx.GetStringOption("id")

		// Obtener información antes de eliminar
		entry, err := database.GetBlacklistEntry(id)
		if err != nil {
			if err == database.ErrBlacklistEntryNotFound {
				sendErrorEmbed(ctx, "Error", fmt.Sprintf("❌ El %s `%s` no está en la blacklist.", getBlacklistTypeName(tipo), id))
				return
			}
			logger.Error(fmt.Sprintf("Error obteniendo entrada de blacklist: %v", err), "DevBlacklist")
			sendErrorEmbed(ctx, "Error", "❌ Error al obtener la entrada de la blacklist.")
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
			Description: fmt.Sprintf("El %s ha sido desbloqueado correctamente.", getBlacklistTypeName(tipo)),
			Color:       0x00FF00, // Verde
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Tipo",
					Value:  getBlacklistTypeEmoji(string(entry.Type)) + " " + getBlacklistTypeName(string(entry.Type)),
					Inline: true,
				},
				{
					Name:   "ID",
					Value:  fmt.Sprintf("`%s`", id),
					Inline: true,
				},
				{
					Name:   "Razón Original",
					Value:  entry.Reason,
					Inline: false,
				},
				{
					Name:   "Bloqueado desde",
					Value:  fmt.Sprintf("<t:%d:R>", entry.CreatedAt.Unix()),
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Desbloqueado por %s", getUserName(ctx)),
			},
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
			logger.Error(fmt.Sprintf("Error enviando respuesta: %v", err), "DevBlacklist")
			return
		}

		logger.Info(fmt.Sprintf("Usuario %s eliminó %s %s de la blacklist", getUserName(ctx), tipo, id), "DevBlacklist")
	}()

	return nil
}

// getBlacklistTypeName devuelve el nombre legible del tipo
func getBlacklistTypeName(tipo string) string {
	if tipo == "user" {
		return "Usuario"
	}
	return "Servidor"
}

// getBlacklistTypeEmoji devuelve el emoji del tipo
func getBlacklistTypeEmoji(tipo string) string {
	if tipo == "user" {
		return "👤"
	}
	return "🏰"
}
