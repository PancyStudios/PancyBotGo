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
					Description: "Canal donde se hará el ping",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "rol",
					Description: "Rol que será mencionado",
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
		var channelID, roleID string
		for _, opt := range subcommand.Options {
			if opt.Name == "canal" {
				channelID = opt.ChannelValue(nil).ID
			}
			if opt.Name == "rol" {
				roleID = opt.RoleValue(nil, "").ID
			}
		}

		err := addPojConfig(ctx.Interaction.GuildID, channelID, roleID)
		if err != nil {
			return ctx.ReplyEphemeral("❌ Hubo un error al guardar.")
		}
		return ctx.Reply(fmt.Sprintf("✅ PoJ añadido: el rol <@&%s> será mencionado en <#%s>.", roleID, channelID))

	case "remove":
		var channelID string
		for _, opt := range subcommand.Options {
			if opt.Name == "canal" {
				channelID = opt.ChannelValue(nil).ID
			}
		}
		
		err := removePojConfig(ctx.Interaction.GuildID, channelID)
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
			list += fmt.Sprintf("• Canal: <#%s> | Rol: <@&%s>\n", poj.ChannelID, poj.RoleID)
		}

		return ctx.Reply(list)
	}

	return nil
}

func addPojConfig(guildID, channelID, roleID string) error {
	doc, err := database.GlobalGuildDM.Get(bson.M{"id": guildID})
	if err != nil {
		doc = models.NewGuildDocument(guildID)
	}

	found := false
	for i, poj := range doc.PingOnJoin {
		if poj.ChannelID == channelID {
			doc.PingOnJoin[i].RoleID = roleID
			found = true
			break
		}
	}
	if !found {
		doc.PingOnJoin = append(doc.PingOnJoin, models.PingOnJoinConfig{
			ChannelID: channelID,
			RoleID:    roleID,
		})
	}

	return database.GlobalGuildDM.Update(bson.M{"id": guildID}, doc)
}

func removePojConfig(guildID, channelID string) error {
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
	return database.GlobalGuildDM.Update(bson.M{"id": guildID}, doc)
}

// HandleMessagePojCommand handles the pan!poj text command
func HandleMessagePojCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Require ManageServer permission
	perms, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil || perms&discordgo.PermissionManageGuild != discordgo.PermissionManageGuild {
		s.ChannelMessageSend(m.ChannelID, "❌ No tienes permisos para usar este comando. Necesitas `Gestionar Servidor`.")
		return
	}

	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "ℹ️ Uso: `pan!poj <add|remove|list>`")
		return
	}

	subcommand := args[0]
	switch subcommand {
	case "add":
		if len(args) < 3 {
			s.ChannelMessageSend(m.ChannelID, "❌ Uso correcto: `pan!poj add <#canal> <@&rol>`")
			return
		}
		
		// Clean up mentions: <#1234> -> 1234, <@&1234> -> 1234
		channelID := args[1]
		if len(channelID) > 4 && channelID[:2] == "<#" {
			channelID = channelID[2 : len(channelID)-1]
		}
		roleID := args[2]
		if len(roleID) > 5 && roleID[:3] == "<@&" {
			roleID = roleID[3 : len(roleID)-1]
		}

		err := addPojConfig(m.GuildID, channelID, roleID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "❌ Hubo un error al guardar.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ PoJ añadido: el rol <@&%s> será mencionado en <#%s>.", roleID, channelID))

	case "remove":
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "❌ Uso correcto: `pan!poj remove <#canal>`")
			return
		}
		
		channelID := args[1]
		if len(channelID) > 4 && channelID[:2] == "<#" {
			channelID = channelID[2 : len(channelID)-1]
		}

		err := removePojConfig(m.GuildID, channelID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "❌ Hubo un error al remover.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Configuración de PoJ removida del canal <#%s>.", channelID))

	case "list":
		doc, err := database.GlobalGuildDM.Get(bson.M{"id": m.GuildID})
		if err != nil || len(doc.PingOnJoin) == 0 {
			s.ChannelMessageSend(m.ChannelID, "ℹ️ No hay configuraciones de Ping On Join activas.")
			return
		}

		list := "🔔 **Lista de Ping On Join (PoJ)**\n"
		for _, poj := range doc.PingOnJoin {
			list += fmt.Sprintf("• Canal: <#%s> | Rol: <@&%s>\n", poj.ChannelID, poj.RoleID)
		}

		s.ChannelMessageSend(m.ChannelID, list)
	default:
		s.ChannelMessageSend(m.ChannelID, "❌ Subcomando desconocido. Usa `add`, `remove` o `list`.")
	}
}
