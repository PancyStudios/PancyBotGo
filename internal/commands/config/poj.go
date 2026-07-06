package config

import (
	"fmt"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

// createPojCommand creates the /poj command
func createPojCommand() *discord.Command {
	cmd := discord.NewCommand(
		"poj",
		"🔔 | Configura el sistema Ping On Join",
		"poj",
		handlePojCommand,
	).WithUserPermissions(discordgo.PermissionManageGuild)

	cmd.Options = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "Añade un Ping On Join para un canal",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "canal",
					Description: "Canal donde se hará el Ping",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "remove",
			Description: "Elimina la configuración de PoJ de un canal",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "canal",
					Description: "Canal a remover",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "Muestra las configuraciones activas de PoJ",
		},
	}

	return cmd
}

func handlePojCommand(ctx *discord.CommandContext) error {
	options := ctx.Interaction.ApplicationCommandData().Options
	if len(options) == 0 {
		return ctx.ReplyEphemeral("❌ Comando inválido.")
	}

	subcommand := options[0]
	switch subcommand.Name {
	case "add":
		channelID := ctx.GetChannelOption("canal").ID

		err := AddPojConfig(ctx.Interaction.GuildID, channelID)
		if err != nil {
			return ctx.Reply("❌ Hubo un error al guardar la configuración de PoJ.")
		}
		return ctx.Reply(fmt.Sprintf("✅ PoJ añadido: el usuario nuevo será mencionado en <#%s>.", channelID))

	case "remove":
		var channelID string
		for _, opt := range subcommand.Options {
			if opt.Name == "canal" {
				channelID = opt.ChannelValue(nil).ID
			}
		}

		err := RemovePojConfig(ctx.Interaction.GuildID, channelID)
		if err != nil {
			return ctx.ReplyEphemeral("❌ Hubo un error al remover.")
		}
		return ctx.Reply(fmt.Sprintf("✅ Configuración de PoJ removida del canal <#%s>.", channelID))

	case "list":
		doc, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Interaction.GuildID})
		if err != nil || len(doc.PingOnJoin) == 0 {
			return ctx.Reply("ℹ️ No hay configuraciones de Ping On Join activas.")
		}

		list := "🔔 **Lista de Ping On Join (PoJ)**\n"
		for _, poj := range doc.PingOnJoin {
			list += fmt.Sprintf("• Canal: <#%s>\n", poj.ChannelID)
		}

		return ctx.Reply(list)
	}

	return nil
}

// AddPojConfig exports the Poj logic
func AddPojConfig(guildID, channelID string) error {
	doc, err := database.GlobalGuildDM.Get(bson.M{"id": guildID})
	if err != nil {
		doc = &models.GuildDocument{ID: guildID}
	}

	found := false
	for _, poj := range doc.PingOnJoin {
		if poj.ChannelID == channelID {
			found = true
			break
		}
	}
	if !found {
		doc.PingOnJoin = append(doc.PingOnJoin, models.PingOnJoinConfig{
			ChannelID: channelID,
		})
	}

	_, err = database.GlobalGuildDM.Set(bson.M{"id": guildID}, doc)
	return err
}

// RemovePojConfig exports the Poj logic
func RemovePojConfig(guildID, channelID string) error {
	doc, err := database.GlobalGuildDM.Get(bson.M{"id": guildID})
	if err != nil {
		return nil
	}

	var newPoj []models.PingOnJoinConfig
	for _, poj := range doc.PingOnJoin {
		if poj.ChannelID != channelID {
			newPoj = append(newPoj, poj)
		}
	}
	_, err = database.GlobalGuildDM.Set(bson.M{"id": guildID}, doc)
	return err
}


