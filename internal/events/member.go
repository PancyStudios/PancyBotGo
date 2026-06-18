package events

import (
	"fmt"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

// RegisterMemberEvents registers all member-related event handlers
func RegisterMemberEvents(client *discord.ExtendedClient) {
	client.Session.AddHandler(onGuildMemberAdd)
	client.Session.AddHandler(onGuildMemberRemove)
	client.Session.AddHandler(onGuildMemberUpdate)
}

// onGuildMemberAdd is called when a new member joins the server
func onGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	logger.Info(fmt.Sprintf("👋 Nuevo miembro: %s#%s en servidor %s",
		m.User.Username, m.User.Discriminator, m.GuildID), "Member")

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error obteniendo servidor: %v", err), "Member")
		return
	}

	// Fetch guild settings from DB
	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"_id": m.GuildID})
	if err == nil && guildDoc != nil {

		// Antibots logic
		if m.User.Bot {
			if guildDoc.Protection.Antibots == "all" {
				s.GuildMemberDeleteWithReason(m.GuildID, m.User.ID, "Anti-Bots: Todos los bots están bloqueados")
				logger.Info(fmt.Sprintf("🤖 Bot %s expulsado por Anti-Bots (all)", m.User.Username), "Member")
				return
			} else if guildDoc.Protection.Antibots == "only_nv" {
				if m.User.PublicFlags&discordgo.UserFlagVerifiedBot == 0 {
					s.GuildMemberDeleteWithReason(m.GuildID, m.User.ID, "Anti-Bots: Bot no verificado")
					logger.Info(fmt.Sprintf("🤖 Bot %s expulsado por Anti-Bots (only_nv)", m.User.Username), "Member")
					return
				}
			} else if guildDoc.Protection.Antibots == "only_v" {
				if m.User.PublicFlags&discordgo.UserFlagVerifiedBot != 0 {
					s.GuildMemberDeleteWithReason(m.GuildID, m.User.ID, "Anti-Bots: Bot verificado no permitido")
					logger.Info(fmt.Sprintf("🤖 Bot %s expulsado por Anti-Bots (only_v)", m.User.Username), "Member")
					return
				}
			}
		}

		// Welcome logic
		if guildDoc.Greetings.Welcome.Enable {
			message := guildDoc.Greetings.Welcome.Message
			if message == "" {
				message = fmt.Sprintf("¡Bienvenido/a {user} a **%s**!", guild.Name)
			}

			// Replace {user} with mention
			message = strings.ReplaceAll(message, "{user}", fmt.Sprintf("<@%s>", m.User.ID))

			channelID := guildDoc.Greetings.Welcome.Channel
			if channelID == "" {
				channelID = guild.SystemChannelID
			}

			welcomeEmbed := &discordgo.MessageEmbed{
				Title:       "¡Bienvenido/a! 🎉",
				Description: message,
				Color:       0x00ff00,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: m.User.AvatarURL("128"),
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    fmt.Sprintf("Ahora somos %d miembros", guild.MemberCount),
					IconURL: guild.IconURL("64"),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			if guildDoc.Greetings.Welcome.IsDM {
				channel, err := s.UserChannelCreate(m.User.ID)
				if err == nil {
					s.ChannelMessageSendEmbed(channel.ID, welcomeEmbed)
				}
			} else if channelID != "" {
				s.ChannelMessageSendEmbed(channelID, welcomeEmbed)
			}
		}

		// Autorole logic
		if guildDoc.Greetings.Autorole.Enable && len(guildDoc.Greetings.Autorole.Roles) > 0 {
			applyRole := func() {
				for _, roleID := range guildDoc.Greetings.Autorole.Roles {
					err := s.GuildMemberRoleAdd(m.GuildID, m.User.ID, roleID)
					if err != nil {
						logger.Error(fmt.Sprintf("Error asignando autorol %s a %s: %v", roleID, m.User.ID, err), "Member")
					} else {
						logger.Debug(fmt.Sprintf("✅ Autorol %s asignado a %s", roleID, m.User.ID), "Member")
					}
				}
			}

			if guildDoc.Greetings.Autorole.Delay > 0 {
				go func() {
					time.Sleep(time.Duration(guildDoc.Greetings.Autorole.Delay) * time.Millisecond)
					applyRole()
				}()
			} else {
				applyRole()
			}
		}
	} else {
		// Fallback to default logic if no DB entry exists
		if guild.SystemChannelID != "" {
			welcomeEmbed := &discordgo.MessageEmbed{
				Title:       "¡Bienvenido/a! 🎉",
				Description: fmt.Sprintf("Dale la bienvenida a <@%s>\nAhora somos **%d** miembros.", m.User.ID, guild.MemberCount),
				Color:       0x00ff00,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: m.User.AvatarURL("128"),
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    guild.Name,
					IconURL: guild.IconURL("64"),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.ChannelMessageSendEmbed(guild.SystemChannelID, welcomeEmbed)
		}
	}
}

// onGuildMemberRemove is called when a member leaves the server
func onGuildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	logger.Info(fmt.Sprintf("👋 Adiós: %s#%s salió del servidor %s",
		m.User.Username, m.User.Discriminator, m.GuildID), "Member")

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		return
	}

	// Fetch guild settings from DB
	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"_id": m.GuildID})
	if err == nil && guildDoc != nil {
		if guildDoc.Greetings.Farewell.Enable {
			message := guildDoc.Greetings.Farewell.Message
			if message == "" {
				message = fmt.Sprintf("👋 **{user}** ha salido del servidor.")
			}

			// Replace {user} with username (not mention since they left)
			message = strings.ReplaceAll(message, "{user}", m.User.Username)

			channelID := guildDoc.Greetings.Farewell.Channel
			if channelID == "" {
				channelID = guild.SystemChannelID
			}

			if channelID != "" {
				farewellEmbed := &discordgo.MessageEmbed{
					Description: message,
					Color:       0xe74c3c,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: m.User.AvatarURL("64"),
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("Ahora somos %d miembros", guild.MemberCount),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.ChannelMessageSendEmbed(channelID, farewellEmbed)
			}
		}
	} else {
		// Fallback to default
		if guild.SystemChannelID != "" {
			farewellEmbed := &discordgo.MessageEmbed{
				Description: fmt.Sprintf("👋 **%s#%s** ha salido del servidor.\nAhora somos **%d** miembros.",
					m.User.Username, m.User.Discriminator, guild.MemberCount),
				Color: 0xe74c3c,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: m.User.AvatarURL("64"),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.ChannelMessageSendEmbed(guild.SystemChannelID, farewellEmbed)
		}
	}
}

// onGuildMemberUpdate is called when a member is updated (roles, nickname, etc.)
func onGuildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	if m.BeforeUpdate != nil {
		oldNick := m.BeforeUpdate.Nick
		newNick := m.Nick

		if oldNick != newNick {
			logger.Debug(fmt.Sprintf("✏️ %s cambió nickname: '%s' → '%s'",
				m.User.Username, oldNick, newNick), "Member")
		}

		if len(m.BeforeUpdate.Roles) != len(m.Roles) {
			logger.Debug(fmt.Sprintf("🎭 Roles actualizados para %s", m.User.Username), "Member")
		}
	}
}
