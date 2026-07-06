package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func autoroleCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionAdministrator) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso de Administrador para configurar el auto-rol.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!autorole <enable/disable> [@rol] [delay_ms]`\nEjemplo: `pan!autorole enable @UsuarioNuevo 5000`")
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

	var roleID string
	if len(ctx.Args) > 1 {
		roleID = ctx.ParseRole(1)
	}

	if enable && roleID == "" {
		_, err := ctx.ReplyError("Error", "❌ Debes especificar un rol cuando activas el auto-rol.")
		return err
	}

	delay := 0
	if len(ctx.Args) > 2 {
		d, err := strconv.Atoi(ctx.Args[2])
		if err == nil {
			delay = d
		}
	}

	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Message.GuildID})
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error obteniendo configuración: %v", err))
		return err
	}

	if guildDoc == nil {
		guildDoc = &models.GuildDocument{ID: ctx.Message.GuildID}
	}

	guildDoc.Greetings.Autorole.Enable = enable

	if roleID != "" {
		guildDoc.Greetings.Autorole.Roles = []string{roleID}
	}

	if len(ctx.Args) > 2 {
		guildDoc.Greetings.Autorole.Delay = delay
	}

	_, err = database.GlobalGuildDM.Set(bson.M{"id": ctx.Message.GuildID}, guildDoc)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error guardando configuración: %v", err))
		return err
	}

	status := "desactivado"
	if enable {
		status = "activado"
	}

	msg := fmt.Sprintf("✅ Auto-rol **%s**.", status)
	if enable {
		msg += fmt.Sprintf("\n**Rol:** <@&%s>", roleID)
		if delay > 0 {
			msg += fmt.Sprintf("\n**Retraso:** %d ms", delay)
		}
	}

	_, err = ctx.ReplySuccess("Auto-rol Configurado", msg)
	return err
}
