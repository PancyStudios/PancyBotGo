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

func autoroleSubcommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "autorole",
		Description: "Configura el sistema de auto-rol",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "enable",
				Description: "Activar o desactivar auto-rol",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "Rol a asignar (requerido si se activa)",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "delay",
				Description: "Retraso en ms antes de asignar (ej. 5000 para 5s)",
				Required:    false,
			},
		},
	}
}

func handleAutorole(ctx *discord.CommandContext, options []*discordgo.ApplicationCommandInteractionDataOption) error {
	var enable bool
	var roleID string
	var delay int

	for _, opt := range options {
		switch opt.Name {
		case "enable":
			enable = opt.BoolValue()
		case "role":
			if r := opt.RoleValue(nil, ""); r != nil {
				roleID = r.ID
			}
		case "delay":
			delay = int(opt.IntValue())
		}
	}

	if enable && roleID == "" {
		// Needs to read current role ID if not provided?
		// We'll require it to be passed if enabling fresh, but for simplicity let's just complain if not provided when enable=true
		return ctx.ReplyEphemeral("❌ Debes especificar un rol cuando activas el auto-rol.")
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

	guildDoc.Greetings.Autorole.Enable = enable

	if roleID != "" {
		guildDoc.Greetings.Autorole.Roles = []string{roleID} // Just setting one for now
	}

	// If delay was provided, set it. Otherwise leave it as is or default to 0
	delayProvided := false
	for _, opt := range options {
		if opt.Name == "delay" {
			delayProvided = true
			break
		}
	}
	if delayProvided {
		guildDoc.Greetings.Autorole.Delay = delay
	}

	_, err = database.GlobalGuildDM.Set(bson.M{"_id": guildID}, guildDoc)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error guardando configuración: %v", err))
	}

	status := "desactivado"
	if enable {
		status = "activado"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("✅ Auto-rol **%s**.", status))
	if enable {
		if len(guildDoc.Greetings.Autorole.Roles) > 0 {
			sb.WriteString(fmt.Sprintf("\nRol: <@&%s>", guildDoc.Greetings.Autorole.Roles[0]))
		}
		if guildDoc.Greetings.Autorole.Delay > 0 {
			sb.WriteString(fmt.Sprintf("\nRetraso: `%d ms`", guildDoc.Greetings.Autorole.Delay))
		}
	}

	return ctx.Reply(sb.String())
}
