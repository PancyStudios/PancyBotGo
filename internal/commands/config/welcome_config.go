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

func welcomeSubcommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "welcome",
		Description: "Configura el sistema de bienvenidas",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "enable",
				Description: "Activar o desactivar las bienvenidas",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "Canal donde se enviarán las bienvenidas",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "Mensaje de bienvenida (usa {user} para mencionar al usuario)",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "is_dm",
				Description: "¿Enviar por mensaje directo en lugar del canal?",
				Required:    false,
			},
		},
	}
}

func handleWelcome(ctx *discord.CommandContext, options []*discordgo.ApplicationCommandInteractionDataOption) error {
	var enable bool
	var channelID string
	var message string
	var isDM bool

	for _, opt := range options {
		switch opt.Name {
		case "enable":
			enable = opt.BoolValue()
		case "channel":
			channelID = opt.ChannelValue(nil).ID
		case "message":
			message = opt.StringValue()
		case "is_dm":
			isDM = opt.BoolValue()
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
		// Document doesn't exist, will be upserted by Set
		guildDoc = &models.GuildDocument{ID: guildID}
	}

	guildDoc.Greetings.Welcome.Enable = enable

	if channelID != "" {
		guildDoc.Greetings.Welcome.Channel = channelID
	}
	if message != "" {
		guildDoc.Greetings.Welcome.Message = message
	}

	// If it was provided, update it (discordgo options don't have a WasProvided method easily accessible here,
	// but default false is fine if we check if it was actually in the options)
	isDMProvided := false
	for _, opt := range options {
		if opt.Name == "is_dm" {
			isDMProvided = true
			break
		}
	}
	if isDMProvided {
		guildDoc.Greetings.Welcome.IsDM = isDM
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
	sb.WriteString(fmt.Sprintf("✅ Bienvenidas **%s**.", status))
	if enable {
		if guildDoc.Greetings.Welcome.Channel != "" && !guildDoc.Greetings.Welcome.IsDM {
			sb.WriteString(fmt.Sprintf("\nCanal: <#%s>", guildDoc.Greetings.Welcome.Channel))
		}
		if guildDoc.Greetings.Welcome.IsDM {
			sb.WriteString("\nModo: **Mensaje Directo (DM)**")
		}
		if guildDoc.Greetings.Welcome.Message != "" {
			sb.WriteString(fmt.Sprintf("\nMensaje: `%s`", guildDoc.Greetings.Welcome.Message))
		}
	}

	return ctx.Reply(sb.String())
}
