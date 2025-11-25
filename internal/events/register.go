// Package events provides a registry for organizing bot events.
// Events are organized by category (guild, member, message, voice, etc.)
package events

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
)

// RegisterAll registers all events with the Discord client
// Add your event registration calls here
func RegisterAll(client *discord.ExtendedClient) {
	logger.System("📋 Registrando eventos del bot...", "Events")

	// Ready event (bot startup)
	RegisterReadyEvent(client)

	// Guild events (server join/leave)
	RegisterGuildEvents(client)

	// Member events (join/leave/update)
	RegisterMemberEvents(client)

	// Message events (create/update/delete)
	RegisterMessageEvents(client)

	// Voice events (join/leave/move)
	RegisterVoiceEvents(client)

	// Add more categories here as needed:
	// RegisterModerationEvents(client)
	// RegisterInteractionEvents(client)

	logger.Success("✅ Todos los eventos registrados correctamente", "Events")
}
