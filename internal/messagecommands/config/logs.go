package config

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func logsCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionAdministrator) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso de Administrador para configurar los logs.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un canal.\nUso: `pan!logs <#canal>`")
		return err
	}

	channelID := ctx.ParseChannel(0)
	if channelID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un canal válido.")
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

	guildDoc.Configuration.LogsChannel = channelID

	_, err = database.GlobalGuildDM.Set(bson.M{"id": ctx.Message.GuildID}, guildDoc)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error guardando configuración: %v", err))
		return err
	}

	_, err = ctx.ReplySuccess("Logs Configurados", fmt.Sprintf("✅ Canal de logs establecido a <#%s>.", channelID))
	return err
}
