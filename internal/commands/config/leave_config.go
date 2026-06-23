package config

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func farewellSubcommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "farewell",
		Description: "⚙️ | Configura el sistema de despedidas",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "enable",
				Description: "⚙️ | Activar o desactivar las despedidas",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "⚙️ | Canal donde se enviarán las despedidas",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "⚙️ | Mensaje de despedida (usa {user} para el nombre del usuario)",
				Required:    false,
			},
		},
	}
}

func handleFarewell(ctx *discord.CommandContext, options []*discordgo.ApplicationCommandInteractionDataOption) error {
	var enable bool
	var channelID string
	var message string

	for _, opt := range options {
		switch opt.Name {
		case "enable":
			enable = opt.BoolValue()
		case "channel":
			channelID = opt.ChannelValue(nil).ID
		case "message":
			message = opt.StringValue()
		}
	}

	guildID := ctx.Interaction.GuildID
	if guildID == "" {
		return ctx.ReplyEphemeral("❌ Este comando solo puede usarse en un servidor.")
	}

	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"_id": guildID})
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error obteniendo configuración: %v", err))
	}

	if guildDoc == nil {
		guildDoc = &models.GuildDocument{ID: guildID}
	}

	guildDoc.Greetings.Farewell.Enable = enable

	if channelID != "" {
		guildDoc.Greetings.Farewell.Channel = channelID
	}
	if message != "" {
		guildDoc.Greetings.Farewell.Message = message
	}

	_, err = database.GlobalGuildDM.Set(bson.M{"_id": guildID}, guildDoc)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error guardando configuración: %v", err))
	}

	status := "desactivadas"
	if enable {
		status = "activadas"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("✅ Despedidas **%s**.", status))
	if enable {
		if guildDoc.Greetings.Farewell.Channel != "" {
			sb.WriteString(fmt.Sprintf("\nCanal: <#%s>", guildDoc.Greetings.Farewell.Channel))
		}
		if guildDoc.Greetings.Farewell.Message != "" {
			sb.WriteString(fmt.Sprintf("\nMensaje: `%s`", guildDoc.Greetings.Farewell.Message))
		}
	}

	return ctx.Reply(sb.String())
}
