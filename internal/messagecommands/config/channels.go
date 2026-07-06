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

func channelsCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionAdministrator) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso de Administrador para configurar canales.")
		return err
	}

	if len(ctx.Args) < 2 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!channels <suggest|confess|verify> <#canal>`")
		return err
	}

	configType := strings.ToLower(ctx.Args[0])
	channelID := ctx.ParseChannel(1)

	if channelID == "" {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un canal válido.")
		return err
	}

	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Message.GuildID})
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Ocurrió un error al cargar la configuración: %v", err))
		return err
	}
	if guildDoc == nil {
		guildDoc = &models.GuildDocument{ID: ctx.Message.GuildID}
	}

	if configType == "suggest" {
		guildDoc.Configuration.SubData.SuggestChannel = channelID
	} else if configType == "confess" {
		guildDoc.Configuration.SubData.ConfessionChannel = channelID
	} else if configType == "verify" {
		guildDoc.Configuration.SubData.VerifyChannel = channelID
	} else {
		_, err = ctx.ReplyError("Uso Incorrecto", "El tipo de configuración debe ser `suggest`, `confess` o `verify`.")
		return err
	}

	_, err = database.GlobalGuildDM.Set(bson.M{"id": ctx.Message.GuildID}, guildDoc)

	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Ocurrió un error al guardar la configuración.")
		return err
	}

	_, err = ctx.ReplySuccess("Canal Configurado", fmt.Sprintf("✅ Canal de **%s** configurado en <#%s>", configType, channelID))
	return err
}
