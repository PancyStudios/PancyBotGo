package config

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func logsSubcommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "logs",
		Description: "Establece el canal de Logs del servidor",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "Canal donde se enviarán los logs",
				Required:    true,
			},
		},
	}
}

func handleLogs(ctx *discord.CommandContext, options []*discordgo.ApplicationCommandInteractionDataOption) error {
	var channelID string

	for _, opt := range options {
		if opt.Name == "channel" {
			channelID = opt.ChannelValue(nil).ID
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

	guildDoc.Configuration.LogsChannel = channelID

	_, err = database.GlobalGuildDM.Set(bson.M{"_id": guildID}, guildDoc)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error guardando configuración: %v", err))
	}

	return ctx.Reply(fmt.Sprintf("✅ Canal de logs establecido a <#%s>.", channelID))
}
