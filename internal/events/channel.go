package events

import (
	"fmt"
	"sync"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	nukeCache = make(map[string][]time.Time)
	nukeMutex sync.Mutex
)

// RegisterChannelEvents registers channel-related events like Anti-Nuke
func RegisterChannelEvents(client *discord.ExtendedClient) {
	client.Session.AddHandler(onChannelDelete)
}

func onChannelDelete(s *discordgo.Session, c *discordgo.ChannelDelete) {
	if c.GuildID == "" {
		return
	}

	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"_id": c.GuildID})
	if err != nil || guildDoc == nil {
		return
	}

	limit := 4
	window := 10 * time.Second

	nukeMutex.Lock()
	deletions := nukeCache[c.GuildID]
	now := time.Now()

	var validDeletions []time.Time
	for _, d := range deletions {
		if now.Sub(d) <= window {
			validDeletions = append(validDeletions, d)
		}
	}
	validDeletions = append(validDeletions, now)
	nukeCache[c.GuildID] = validDeletions
	nukeCount := len(validDeletions)
	nukeMutex.Unlock()

	// If nuke is detected
	if nukeCount >= limit {
		// Prevent repeated triggers in the same window
		nukeMutex.Lock()
		delete(nukeCache, c.GuildID)
		nukeMutex.Unlock()

		logger.Warn(fmt.Sprintf("⚠️ ANTI-NUKE TRIGGERED in guild %s! %d channels deleted.", c.GuildID, nukeCount), "AntiNuke")

		// Activate Panic Mode
		antiRaid := &guildDoc.Protection.AntiRaid
		antiRaid.Enable = true
		database.GlobalGuildDM.Set(bson.M{"_id": c.GuildID}, guildDoc)

		// Find the culprit via Audit Logs
		auditLog, err := s.GuildAuditLog(c.GuildID, "", "", int(discordgo.AuditLogActionChannelDelete), 10)
		culpritID := ""
		if err == nil && len(auditLog.AuditLogEntries) > 0 {
			entry := auditLog.AuditLogEntries[0]
			culpritID = entry.UserID
		}

		guild, err := s.Guild(c.GuildID)
		
		var culpritName string
		if culpritID != "" {
			user, err := s.User(culpritID)
			if err == nil {
				culpritName = user.Username
			} else {
				culpritName = culpritID
			}

			if culpritID != s.State.User.ID {
				// Put in Quarantine (Timeout for 28 days)
				until := time.Now().Add(28 * 24 * time.Hour)
				errTimeout := s.GuildMemberTimeout(c.GuildID, culpritID, &until)
				
				if errTimeout != nil {
					// Fallback: Strip all roles
					member, err := s.GuildMember(c.GuildID, culpritID)
					if err == nil {
						for _, roleID := range member.Roles {
							s.GuildMemberRoleRemove(c.GuildID, culpritID, roleID)
						}
					}
				}
			}
		}

		// Send Alerts
		alertChannel := guildDoc.Configuration.LogsChannel
		if alertChannel == "" && guild != nil {
			alertChannel = guild.SystemChannelID
		}

		desc := "Se detectó un borrado masivo de canales.\n\nEl **Modo Pánico** ha sido activado automáticamente."
		if culpritID != "" {
			desc += fmt.Sprintf("\n\nEl presunto atacante **%s** ha sido puesto en **Cuarentena** (Timeout/Roles removidos) para evitar más daños.", culpritName)
		}

		embedAlert := discord.NewEmbed().
			SetColor(0xFF0000).
			SetTitle("☢️ ¡ALERTA ANTI-NUKE!").
			SetDescription(desc).
			Build()

		if alertChannel != "" {
			s.ChannelMessageSendEmbed(alertChannel, embedAlert)
		}

		if guild != nil && guild.OwnerID != "" {
			dmChannel, err := s.UserChannelCreate(guild.OwnerID)
			if err == nil {
				s.ChannelMessageSendEmbed(dmChannel.ID, embedAlert)
			}
		}
	}
}
