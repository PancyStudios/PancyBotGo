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

func welcomeCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionAdministrator) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso de Administrador para configurar las bienvenidas.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!welcome <enable/disable> [#canal] [is_dm: true/false] [mensaje...]`")
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

	guildDoc.Greetings.Welcome.Enable = enable

	if len(ctx.Args) > 1 {
		channelID := ctx.ParseChannel(1)
		if channelID != "" {
			guildDoc.Greetings.Welcome.Channel = channelID
		}
	}

	if len(ctx.Args) > 2 {
		dmArg := strings.ToLower(ctx.Args[2])
		if dmArg == "true" || dmArg == "false" {
			guildDoc.Greetings.Welcome.IsDM = (dmArg == "true")

			if len(ctx.Args) > 3 {
				guildDoc.Greetings.Welcome.Message = strings.Join(ctx.Args[3:], " ")
			}
		} else {
			guildDoc.Greetings.Welcome.Message = strings.Join(ctx.Args[2:], " ")
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

	msg := fmt.Sprintf("✅ Bienvenidas **%s**.", status)
	if enable {
		if guildDoc.Greetings.Welcome.Channel != "" {
			msg += fmt.Sprintf("\n**Canal:** <#%s>", guildDoc.Greetings.Welcome.Channel)
		}
		msg += fmt.Sprintf("\n**Por DM:** %v", guildDoc.Greetings.Welcome.IsDM)
		if guildDoc.Greetings.Welcome.Message != "" {
			msg += fmt.Sprintf("\n**Mensaje:** %s", guildDoc.Greetings.Welcome.Message)
		}
	}

	_, err = ctx.ReplySuccess("Bienvenidas Configuradas", msg)
	return err
}
