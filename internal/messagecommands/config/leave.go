package config

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func leaveCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionAdministrator) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso de Administrador para configurar las despedidas.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!leave <enable/disable> [#canal] [mensaje...]`")
		return err
	}

	action := strings.ToLower(ctx.Args[0])
	enable := false
	if action == "enable" || action == "on" || action == "true" {
		enable = true
	} else if action == "disable" || action == "off" || action == "false" {
		enable = false
	} else {
		_, err := ctx.ReplyError("Uso Incorrecto", "El primer argumento debe ser `enable` o `disable`.")
		return err
	}

	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Message.GuildID})
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error obteniendo configuración: %v", err))
		return err
	}

	if guildDoc == nil {
		guildDoc = &models.GuildDocument{ID: ctx.Message.GuildID}
	}

	guildDoc.Greetings.Farewell.Enable = enable

	if len(ctx.Args) > 1 {
		channelID := ctx.ParseChannel(1)
		if channelID != "" {
			guildDoc.Greetings.Farewell.Channel = channelID
			if len(ctx.Args) > 2 {
				guildDoc.Greetings.Farewell.Message = strings.Join(ctx.Args[2:], " ")
			}
		} else {
			guildDoc.Greetings.Farewell.Message = strings.Join(ctx.Args[1:], " ")
		}
	}

	_, err = database.GlobalGuildDM.Set(bson.M{"id": ctx.Message.GuildID}, guildDoc)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error guardando configuración: %v", err))
		return err
	}

	status := "desactivadas"
	if enable {
		status = "activadas"
	}

	msg := fmt.Sprintf("✅ Despedidas **%s**.", status)
	if enable {
		if guildDoc.Greetings.Farewell.Channel != "" {
			msg += fmt.Sprintf("\nCanal: <#%s>", guildDoc.Greetings.Farewell.Channel)
		}
		if guildDoc.Greetings.Farewell.Message != "" {
			msg += fmt.Sprintf("\nMensaje: `%s`", guildDoc.Greetings.Farewell.Message)
		}
	}

	_, err = ctx.ReplySuccess("Despedidas Configuradas", msg)
	return err
}
