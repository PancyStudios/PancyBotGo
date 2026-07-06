package config

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func createSuggestConfigCommand() *discord.Command {
	return &discord.Command{
		Name:        "config-suggest",
		Description: "⚙️ | Configura el canal de sugerencias",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionChannel,
				Name:         "canal",
				Description:  "⚙️ | El canal donde se enviarán las sugerencias",
				Required:     true,
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
			},
		},
		Run: func(ctx *discord.CommandContext) error {
			channelOpt := ctx.GetChannelOption("canal")
			if channelOpt == nil {
				return ctx.ReplyEphemeral("❌ Canal inválido.")
			}
			channelID := channelOpt.ID

			guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Interaction.GuildID})
			if err != nil {
				return ctx.ReplyEphemeral("❌ Ocurrió un error al cargar la configuración.")
			}
			if guildDoc == nil {
				guildDoc = &models.GuildDocument{ID: ctx.Interaction.GuildID}
			}

			guildDoc.Configuration.SubData.SuggestChannel = channelID
			_, err = database.GlobalGuildDM.Set(bson.M{"id": ctx.Interaction.GuildID}, guildDoc)

			if err != nil {
				logger.Error(fmt.Sprintf("Error actualizando base de datos: %v", err), "Config")
				return ctx.ReplyEphemeral("❌ Ocurrió un error al guardar la configuración.")
			}

			return ctx.ReplyEphemeral(fmt.Sprintf("✅ Canal de sugerencias configurado en <#%s>", channelID))
		},
	}
}

func createConfessConfigCommand() *discord.Command {
	return &discord.Command{
		Name:        "config-confess",
		Description: "⚙️ | Configura el canal de confesiones",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionChannel,
				Name:         "canal",
				Description:  "⚙️ | El canal donde se enviarán las confesiones",
				Required:     true,
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
			},
		},
		Run: func(ctx *discord.CommandContext) error {
			channelOpt := ctx.GetChannelOption("canal")
			if channelOpt == nil {
				return ctx.ReplyEphemeral("❌ Canal inválido.")
			}
			channelID := channelOpt.ID

			guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Interaction.GuildID})
			if err != nil {
				return ctx.ReplyEphemeral("❌ Ocurrió un error al cargar la configuración.")
			}
			if guildDoc == nil {
				guildDoc = &models.GuildDocument{ID: ctx.Interaction.GuildID}
			}

			guildDoc.Configuration.SubData.ConfessionChannel = channelID
			_, err = database.GlobalGuildDM.Set(bson.M{"id": ctx.Interaction.GuildID}, guildDoc)

			if err != nil {
				logger.Error(fmt.Sprintf("Error actualizando base de datos: %v", err), "Config")
				return ctx.ReplyEphemeral("❌ Ocurrió un error al guardar la configuración.")
			}

			return ctx.ReplyEphemeral(fmt.Sprintf("✅ Canal de confesiones configurado en <#%s>", channelID))
		},
	}
}

func createVerifyChannelCommand() *discord.Command {
	return &discord.Command{
		Name:        "config-verifychannel",
		Description: "⚙️ | Configura el canal de verificación",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionChannel,
				Name:         "canal",
				Description:  "⚙️ | El canal donde se enviará el panel de verificación",
				Required:     true,
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
			},
		},
		Run: func(ctx *discord.CommandContext) error {
			channelOpt := ctx.GetChannelOption("canal")
			if channelOpt == nil {
				return ctx.ReplyEphemeral("❌ Canal inválido.")
			}
			channelID := channelOpt.ID

			guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Interaction.GuildID})
			if err != nil {
				return ctx.ReplyEphemeral("❌ Ocurrió un error al cargar la configuración.")
			}
			if guildDoc == nil {
				guildDoc = &models.GuildDocument{ID: ctx.Interaction.GuildID}
			}

			guildDoc.Configuration.SubData.VerifyChannel = channelID
			_, err = database.GlobalGuildDM.Set(bson.M{"id": ctx.Interaction.GuildID}, guildDoc)

			if err != nil {
				logger.Error(fmt.Sprintf("Error actualizando base de datos: %v", err), "Config")
				return ctx.ReplyEphemeral("❌ Ocurrió un error al guardar la configuración.")
			}

			return ctx.ReplyEphemeral(fmt.Sprintf("✅ Canal de verificación configurado en <#%s>", channelID))
		},
	}
}

func createVerifyRoleCommand() *discord.Command {
	return &discord.Command{
		Name:        "config-verifyrole",
		Description: "⚙️ | Configura el rol que se dará al verificarse",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "rol",
				Description: "⚙️ | El rol de verificado",
				Required:    true,
			},
		},
		Run: func(ctx *discord.CommandContext) error {
			roleOpt := ctx.GetRoleOption("rol")
			if roleOpt == nil {
				return ctx.ReplyEphemeral("❌ Rol inválido.")
			}
			roleID := roleOpt.ID

			guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Interaction.GuildID})
			if err != nil {
				return ctx.ReplyEphemeral("❌ Ocurrió un error al cargar la configuración.")
			}
			if guildDoc == nil {
				guildDoc = &models.GuildDocument{ID: ctx.Interaction.GuildID}
			}

			guildDoc.Configuration.SubData.VerifyRole = roleID
			_, err = database.GlobalGuildDM.Set(bson.M{"id": ctx.Interaction.GuildID}, guildDoc)

			if err != nil {
				logger.Error(fmt.Sprintf("Error actualizando base de datos: %v", err), "Config")
				return ctx.ReplyEphemeral("❌ Ocurrió un error al guardar la configuración.")
			}

			return ctx.ReplyEphemeral(fmt.Sprintf("✅ Rol de verificación configurado como <@&%s>", roleID))
		},
	}
}

func createSendVerifyCommand() *discord.Command {
	return &discord.Command{
		Name:        "config-sendverify",
		Description: "⚙️ | Envía el panel de verificación al canal configurado",
		Run: func(ctx *discord.CommandContext) error {
			guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Interaction.GuildID})
			if err != nil || guildDoc == nil || guildDoc.Configuration.SubData.VerifyChannel == "" || guildDoc.Configuration.SubData.VerifyRole == "" {
				return ctx.ReplyEphemeral("❌ Debes configurar primero el canal (`/config verifychannel`) y el rol (`/config verifyrole`).")
			}

			channelID := guildDoc.Configuration.SubData.VerifyChannel

			embed := &discordgo.MessageEmbed{
				Title:       "🔒 Sistema de Verificación",
				Description: "⚙️ | Haz clic en el botón de abajo para verificarte y acceder al servidor.",
				Color:       0x2ECC71, // Green
			}

			btn := discordgo.Button{
				Label:    "Verificarse",
				Style:    discordgo.SuccessButton,
				CustomID: "btn_verify_user",
				Emoji:    &discordgo.ComponentEmoji{Name: "✅"},
			}

			row := discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{btn},
			}

			_, err = ctx.Session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Embeds:     []*discordgo.MessageEmbed{embed},
				Components: []discordgo.MessageComponent{row},
			})

			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando panel de verificación: %v", err), "Config")
				return ctx.ReplyEphemeral("❌ No pude enviar el panel de verificación al canal. Revisa mis permisos.")
			}

			return ctx.ReplyEphemeral(fmt.Sprintf("✅ Panel enviado correctamente a <#%s>", channelID))
		},
	}
}
