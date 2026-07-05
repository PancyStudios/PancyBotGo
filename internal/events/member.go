package events

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	raidCache = make(map[string][]time.Time)
	raidMutex sync.Mutex
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
	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": m.GuildID})
	if err == nil && guildDoc != nil {

		// Anti-Raid Logic
		antiRaid := &guildDoc.Protection.AntiRaid

		// 1. Min Account Age
		if antiRaid.MinAccountAgeDays > 0 {
			idInt, err := strconv.ParseInt(m.User.ID, 10, 64)
			if err == nil {
				createdAt := time.UnixMilli((idInt >> 22) + 1420070400000)
				ageDays := int(time.Since(createdAt).Hours() / 24)

				if ageDays < antiRaid.MinAccountAgeDays {
					reason := fmt.Sprintf("Anti-Raid: Cuenta muy reciente (%d días < %d días)", ageDays, antiRaid.MinAccountAgeDays)
					if antiRaid.Action == "ban" {
						s.GuildBanCreateWithReason(m.GuildID, m.User.ID, reason, 0)
					} else {
						s.GuildMemberDeleteWithReason(m.GuildID, m.User.ID, reason)
					}
					logger.Info(fmt.Sprintf("🛡️ %s expulsado/baneado por Anti-Raid (Edad de cuenta: %d días)", m.User.Username, ageDays), "AntiRaid")
					return
				}
			}
		}

		// 2. Active Panic Mode
		if antiRaid.Enable {
			reason := "Anti-Raid: Modo pánico activado"
			if antiRaid.Action == "ban" {
				s.GuildBanCreateWithReason(m.GuildID, m.User.ID, reason, 0)
			} else {
				s.GuildMemberDeleteWithReason(m.GuildID, m.User.ID, reason)
			}
			logger.Info(fmt.Sprintf("🛡️ %s expulsado/baneado por Modo Pánico Anti-Raid", m.User.Username), "AntiRaid")
			return
		}

		// 3. Raid Detection
		if antiRaid.JoinLimit > 0 && antiRaid.TimeWindow > 0 {
			raidMutex.Lock()
			joins := raidCache[m.GuildID]
			now := time.Now()
			window := time.Duration(antiRaid.TimeWindow) * time.Second

			var validJoins []time.Time
			for _, j := range joins {
				if now.Sub(j) <= window {
					validJoins = append(validJoins, j)
				}
			}
			validJoins = append(validJoins, now)
			raidCache[m.GuildID] = validJoins
			raidCount := len(validJoins)
			raidMutex.Unlock()

			if raidCount >= antiRaid.JoinLimit {
				// Trigger Panic Mode!
				antiRaid.Enable = true
				database.GlobalGuildDM.Set(bson.M{"id": m.GuildID}, guildDoc)

				logger.Warn(fmt.Sprintf("🚨 POSIBLE RAID DETECTADO en %s! Modo Pánico activado automáticamente.", m.GuildID), "AntiRaid")

				// Kick current user
				reason := "Anti-Raid: Límite de uniones superado (Modo Pánico Auto-Activado)"
				if antiRaid.Action == "ban" {
					s.GuildBanCreateWithReason(m.GuildID, m.User.ID, reason, 0)
				} else {
					s.GuildMemberDeleteWithReason(m.GuildID, m.User.ID, reason)
				}

				// Try to alert in system channel or logs channel
				alertChannel := guildDoc.Configuration.LogsChannel
				if alertChannel == "" {
					alertChannel = guild.SystemChannelID
				}

				embedAlert := discord.NewEmbed().
					SetColor(0xFF0000). // Red
					SetTitle("🚨 ¡ALERTA DE RAID MASIVO!").
					SetDescription("Se detectó un pico inusual de nuevas cuentas uniéndose al servidor.\n\nEl **Modo Pánico Anti-Raid** ha sido activado automáticamente y todas las nuevas uniones serán bloqueadas.\n\n*Un administrador debe desactivarlo con `/security antiraid toggle` cuando sea seguro.*").
					Build()

				if alertChannel != "" {
					s.ChannelMessageSendEmbed(alertChannel, embedAlert)
				}

				// Alert the owner via DM
				dmChannel, err := s.UserChannelCreate(guild.OwnerID)
				if err == nil {
					s.ChannelMessageSendEmbed(dmChannel.ID, embedAlert)
				}
				return
			}
		}

		// Antibots logic
		if m.User.Bot && guildDoc.Protection.Antibots.Enable {
			if guildDoc.Protection.Antibots.Type == "all" {
				s.GuildMemberDeleteWithReason(m.GuildID, m.User.ID, "Anti-Bots: Todos los bots están bloqueados")
				logger.Info(fmt.Sprintf("🤖 Bot %s expulsado por Anti-Bots (all)", m.User.Username), "Member")
				return
			} else if guildDoc.Protection.Antibots.Type == "only_nv" {
				if m.User.PublicFlags&discordgo.UserFlagVerifiedBot == 0 {
					s.GuildMemberDeleteWithReason(m.GuildID, m.User.ID, "Anti-Bots: Bot no verificado")
					logger.Info(fmt.Sprintf("🤖 Bot %s expulsado por Anti-Bots (only_nv)", m.User.Username), "Member")
					return
				}
			} else if guildDoc.Protection.Antibots.Type == "only_v" {
				if m.User.PublicFlags&discordgo.UserFlagVerifiedBot != 0 {
					s.GuildMemberDeleteWithReason(m.GuildID, m.User.ID, "Anti-Bots: Bot verificado no permitido")
					logger.Info(fmt.Sprintf("🤖 Bot %s expulsado por Anti-Bots (only_v)", m.User.Username), "Member")
					return
				}
			}
		}

		// Welcome logic
		if guildDoc.Greetings.Welcome.Enable {
			messageText := guildDoc.Greetings.Welcome.Message
			if messageText != "" {
				messageText = strings.ReplaceAll(messageText, "{user}", fmt.Sprintf("<@%s>", m.User.ID))
				messageText = strings.ReplaceAll(messageText, "{server}", guild.Name)
			}

			channelID := guildDoc.Greetings.Welcome.Channel
			if channelID == "" {
				channelID = guild.SystemChannelID
			}

			var welcomeEmbed *discordgo.MessageEmbed

			// Revisar si usa Custom Embed
			if guildDoc.Greetings.Welcome.EmbedID != "" {
				for _, ce := range guildDoc.Embeds {
					if ce.ID == guildDoc.Greetings.Welcome.EmbedID {
						welcomeEmbed = buildCustomEmbed(ce, m.User, guild)
						break
					}
				}
			}

			// Fallback si no hay custom embed pero tampoco hay texto, enviamos un embed por defecto
			if welcomeEmbed == nil && messageText == "" {
				welcomeEmbed = &discordgo.MessageEmbed{
					Title:       "¡Bienvenido/a! 🎉",
					Description: fmt.Sprintf("¡Bienvenido/a <@%s> a **%s**!", m.User.ID, guild.Name),
					Color:       0x00ff00,
					Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: m.User.AvatarURL("128")},
					Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Ahora somos %d miembros", guild.MemberCount), IconURL: guild.IconURL("64")},
					Timestamp:   time.Now().Format(time.RFC3339),
				}
			}

			sendData := &discordgo.MessageSend{
				Content: messageText,
			}
			if welcomeEmbed != nil {
				sendData.Embeds = []*discordgo.MessageEmbed{welcomeEmbed}
			}

			if guildDoc.Greetings.Welcome.IsDM {
				channel, err := s.UserChannelCreate(m.User.ID)
				if err == nil {
					s.ChannelMessageSendComplex(channel.ID, sendData)
				}
			} else if channelID != "" {
				s.ChannelMessageSendComplex(channelID, sendData)
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
	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": m.GuildID})
	if err == nil && guildDoc != nil {
		// Farewell logic
		if guildDoc.Greetings.Farewell.Enable {
			messageText := guildDoc.Greetings.Farewell.Message
			if messageText != "" {
				messageText = strings.ReplaceAll(messageText, "{user}", m.User.Username)
				messageText = strings.ReplaceAll(messageText, "{server}", guild.Name)
			}

			channelID := guildDoc.Greetings.Farewell.Channel
			if channelID == "" {
				channelID = guild.SystemChannelID
			}

			var farewellEmbed *discordgo.MessageEmbed

			if guildDoc.Greetings.Farewell.EmbedID != "" {
				for _, ce := range guildDoc.Embeds {
					if ce.ID == guildDoc.Greetings.Farewell.EmbedID {
						farewellEmbed = buildCustomEmbed(ce, m.User, guild)
						break
					}
				}
			}

			if farewellEmbed == nil && messageText == "" {
				farewellEmbed = &discordgo.MessageEmbed{
					Title:       "Despedida 👋",
					Description: fmt.Sprintf("**%s** ha salido del servidor.", m.User.Username),
					Color:       0xff0000,
					Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: m.User.AvatarURL("128")},
				}
			}

			sendData := &discordgo.MessageSend{
				Content: messageText,
			}
			if farewellEmbed != nil {
				sendData.Embeds = []*discordgo.MessageEmbed{farewellEmbed}
			}

			if channelID != "" {
				s.ChannelMessageSendComplex(channelID, sendData)
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

// buildCustomEmbed converte un CustomEmbed del DB a un discordgo.MessageEmbed e interpola variables
func buildCustomEmbed(customEmbed models.CustomEmbed, user *discordgo.User, guild *discordgo.Guild) *discordgo.MessageEmbed {
	replaceVars := func(s string) string {
		s = strings.ReplaceAll(s, "{user}", fmt.Sprintf("<@%s>", user.ID))
		s = strings.ReplaceAll(s, "{user.id}", user.ID)
		s = strings.ReplaceAll(s, "{user.name}", user.Username)
		s = strings.ReplaceAll(s, "{user.avatar}", user.AvatarURL("256"))
		
		s = strings.ReplaceAll(s, "{server}", guild.Name)
		s = strings.ReplaceAll(s, "{server.name}", guild.Name)
		s = strings.ReplaceAll(s, "{server.id}", guild.ID)
		s = strings.ReplaceAll(s, "{server.icon}", guild.IconURL("256"))
		s = strings.ReplaceAll(s, "{server.members}", fmt.Sprintf("%d", guild.MemberCount))
		
		// Retro-compatibilidad
		s = strings.ReplaceAll(s, "{username}", user.Username)
		s = strings.ReplaceAll(s, "{guild.id}", guild.ID)
		s = strings.ReplaceAll(s, "{guild.name}", guild.Name)
		return s
	}

	embed := &discordgo.MessageEmbed{
		Title:       replaceVars(customEmbed.Title),
		Description: replaceVars(customEmbed.Description),
		Color:       customEmbed.Color,
	}
	if customEmbed.Thumbnail != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: replaceVars(customEmbed.Thumbnail)}
	}
	if customEmbed.Image != "" {
		embed.Image = &discordgo.MessageEmbedImage{URL: replaceVars(customEmbed.Image)}
	}
	if customEmbed.AuthorName != "" || customEmbed.AuthorIcon != "" {
		embed.Author = &discordgo.MessageEmbedAuthor{Name: replaceVars(customEmbed.AuthorName), IconURL: replaceVars(customEmbed.AuthorIcon)}
	}
	if customEmbed.FooterText != "" || customEmbed.FooterIcon != "" {
		embed.Footer = &discordgo.MessageEmbedFooter{Text: replaceVars(customEmbed.FooterText), IconURL: replaceVars(customEmbed.FooterIcon)}
	}
	return embed
}
