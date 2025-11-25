// Package events provides event handlers for voice events
package events

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// RegisterVoiceEvents registers all voice-related event handlers
func RegisterVoiceEvents(client *discord.ExtendedClient) {
	client.Session.AddHandler(onVoiceStateUpdate)
}

// onVoiceStateUpdate is called when a user's voice state changes
func onVoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	// Usuario se unió a un canal de voz
	if v.ChannelID != "" && (v.BeforeUpdate == nil || v.BeforeUpdate.ChannelID == "") {
		channel, err := s.Channel(v.ChannelID)
		if err == nil {
			user, _ := s.User(v.UserID)
			if user != nil {
				logger.Debug(fmt.Sprintf("🎤 %s se unió a: %s", user.Username, channel.Name), "Voice")
			}
		}
		return
	}

	// Usuario salió de un canal de voz
	if v.ChannelID == "" && v.BeforeUpdate != nil && v.BeforeUpdate.ChannelID != "" {
		user, _ := s.User(v.UserID)
		if user != nil {
			logger.Debug(fmt.Sprintf("🔇 %s salió del canal de voz", user.Username), "Voice")
		}
		return
	}

	// Usuario cambió de canal de voz
	if v.ChannelID != "" && v.BeforeUpdate != nil &&
		v.BeforeUpdate.ChannelID != "" && v.ChannelID != v.BeforeUpdate.ChannelID {
		oldChannel, _ := s.Channel(v.BeforeUpdate.ChannelID)
		newChannel, _ := s.Channel(v.ChannelID)
		user, _ := s.User(v.UserID)

		if oldChannel != nil && newChannel != nil && user != nil {
			logger.Debug(fmt.Sprintf("🔄 %s: %s → %s",
				user.Username, oldChannel.Name, newChannel.Name), "Voice")
		}
		return
	}

	// Usuario silenciado/desilenciado
	if v.BeforeUpdate != nil {
		user, _ := s.User(v.UserID)
		if user != nil {
			if v.Mute && !v.BeforeUpdate.Mute {
				logger.Debug(fmt.Sprintf("🔇 %s fue silenciado", user.Username), "Voice")
			} else if !v.Mute && v.BeforeUpdate.Mute {
				logger.Debug(fmt.Sprintf("🔊 %s fue desilenciado", user.Username), "Voice")
			}

			if v.Deaf && !v.BeforeUpdate.Deaf {
				logger.Debug(fmt.Sprintf("🔇 %s fue ensorDecido", user.Username), "Voice")
			} else if !v.Deaf && v.BeforeUpdate.Deaf {
				logger.Debug(fmt.Sprintf("🔊 %s dejó de estar ensordecido", user.Username), "Voice")
			}
		}
	}
}
