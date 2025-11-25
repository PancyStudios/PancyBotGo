// Package events - EXAMPLE FILE
// This is a template for creating new event handlers
// Copy this file and rename it to your event category (e.g., moderation.go, logging.go, etc.)

package events

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// STEP 1: Create a registration function
// This function will be called from register.go
func RegisterExampleEvents(client *discord.ExtendedClient) {
	// Register your event handlers here using AddHandler directly:

	// For any Discord event, use client.Session.AddHandler:
	client.Session.AddHandler(onExampleEvent)
	client.Session.AddHandler(onChannelCreate)
	client.Session.AddHandler(onRoleCreate)
	client.Session.AddHandler(onBanRemove)
	client.Session.AddHandler(onReactionAdd)
	// client.Session.AddHandler(onPresenceUpdate)  // Careful: fires very frequently!

	logger.Debug("Example events registered", "Events")
}

// STEP 2: Create your event handler functions
// Each handler receives the Discord session and event data

// Example: Guild Ban Add event
func onExampleEvent(s *discordgo.Session, b *discordgo.GuildBanAdd) {
	logger.Info(fmt.Sprintf("🔨 User banned: %s#%s", b.User.Username, b.User.Discriminator), "Example")

	// Get guild info
	guild, err := s.Guild(b.GuildID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting guild: %v", err), "Example")
		return
	}

	// Send a message to a log channel (example)
	logChannelID := "YOUR_LOG_CHANNEL_ID" // Replace with actual channel ID

	embed := &discordgo.MessageEmbed{
		Title:       "🔨 User Banned",
		Description: fmt.Sprintf("**%s#%s** was banned from the server", b.User.Username, b.User.Discriminator),
		Color:       0xff0000,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User ID",
				Value:  b.User.ID,
				Inline: true,
			},
			{
				Name:   "Guild",
				Value:  guild.Name,
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: b.User.AvatarURL("128"),
		},
	}

	_, err = s.ChannelMessageSendEmbed(logChannelID, embed)
	if err != nil {
		logger.Error(fmt.Sprintf("Error sending ban notification: %v", err), "Example")
	}
}

// STEP 3: Add your registration function to register.go
// In internal/events/register.go, add:
/*
func RegisterAll(client *discord.ExtendedClient) {
	logger.System("📋 Registrando eventos del bot...", "Events")

	// ... existing registrations ...

	// Add your new registration here:
	RegisterExampleEvents(client)

	logger.Success("✅ Todos los eventos registrados correctamente", "Events")
}
*/

// ==============================================================================
// MORE EXAMPLES OF EVENT HANDLERS
// ==============================================================================

// Example: Channel Create
func onChannelCreate(s *discordgo.Session, c *discordgo.ChannelCreate) {
	logger.Info(fmt.Sprintf("📝 New channel created: %s", c.Name), "Example")
}

// Example: Role Create
func onRoleCreate(s *discordgo.Session, r *discordgo.GuildRoleCreate) {
	logger.Info(fmt.Sprintf("🎭 New role created: %s", r.Role.Name), "Example")
}

// Example: Ban Remove
func onBanRemove(s *discordgo.Session, b *discordgo.GuildBanRemove) {
	logger.Info(fmt.Sprintf("✅ Ban removed for: %s", b.User.Username), "Example")
}

// Example: Message Reaction Add
func onReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	logger.Debug(fmt.Sprintf("👍 Reaction added: %s", r.Emoji.Name), "Example")

	// Example: Auto-role on reaction
	if r.Emoji.Name == "✅" {
		// Add a role to the user
		roleID := "YOUR_ROLE_ID"
		err := s.GuildMemberRoleAdd(r.GuildID, r.UserID, roleID)
		if err != nil {
			logger.Error(fmt.Sprintf("Error adding role: %v", err), "Example")
		}
	}
}

// Example: Presence Update (user status change)
func onPresenceUpdate(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	// Note: This event fires VERY frequently, use with caution
	// Only log in debug mode or for specific users

	if p.Status == discordgo.StatusOnline {
		logger.Debug(fmt.Sprintf("👤 %s is now online", p.User.Username), "Example")
	}
}

// ==============================================================================
// AVAILABLE DISCORD EVENTS
// ==============================================================================
/*
You can handle any of these events using client.EventHandler.RegisterEvent():

Bot Events:
- *discordgo.Ready                    - Bot is ready
- *discordgo.Resumed                  - Bot resumed connection

Guild Events:
- *discordgo.GuildCreate              - Bot joined guild
- *discordgo.GuildUpdate              - Guild updated
- *discordgo.GuildDelete              - Bot left guild
- *discordgo.GuildBanAdd              - User banned
- *discordgo.GuildBanRemove           - Ban removed
- *discordgo.GuildMemberAdd           - Member joined
- *discordgo.GuildMemberUpdate        - Member updated
- *discordgo.GuildMemberRemove        - Member left
- *discordgo.GuildRoleCreate          - Role created
- *discordgo.GuildRoleUpdate          - Role updated
- *discordgo.GuildRoleDelete          - Role deleted
- *discordgo.GuildEmojisUpdate        - Emojis updated

Channel Events:
- *discordgo.ChannelCreate            - Channel created
- *discordgo.ChannelUpdate            - Channel updated
- *discordgo.ChannelDelete            - Channel deleted
- *discordgo.ChannelPinsUpdate        - Pinned messages updated

Message Events:
- *discordgo.MessageCreate            - Message sent
- *discordgo.MessageUpdate            - Message edited
- *discordgo.MessageDelete            - Message deleted
- *discordgo.MessageDeleteBulk        - Multiple messages deleted
- *discordgo.MessageReactionAdd       - Reaction added
- *discordgo.MessageReactionRemove    - Reaction removed
- *discordgo.MessageReactionRemoveAll - All reactions removed

Voice Events:
- *discordgo.VoiceStateUpdate         - Voice state changed
- *discordgo.VoiceServerUpdate        - Voice server updated

Interaction Events:
- *discordgo.InteractionCreate        - Interaction created (buttons, menus, modals)

Other Events:
- *discordgo.PresenceUpdate           - User presence changed
- *discordgo.TypingStart              - User started typing
- *discordgo.UserUpdate               - User updated
- *discordgo.InviteCreate             - Invite created
- *discordgo.InviteDelete             - Invite deleted
*/
