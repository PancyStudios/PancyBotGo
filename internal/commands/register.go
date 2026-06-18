// Package commands provides a registry for organizing bot commands.
// Commands are organized in subdirectories by category (util, music, mod, etc.)
package commands

import (
	"github.com/PancyStudios/PancyBotGo/internal/commands/config"
	"github.com/PancyStudios/PancyBotGo/internal/commands/dev"
	"github.com/PancyStudios/PancyBotGo/internal/commands/economy"
	"github.com/PancyStudios/PancyBotGo/internal/commands/embeds"
	"github.com/PancyStudios/PancyBotGo/internal/commands/fun"
	"github.com/PancyStudios/PancyBotGo/internal/commands/ia"
	"github.com/PancyStudios/PancyBotGo/internal/commands/mod"
	"github.com/PancyStudios/PancyBotGo/internal/commands/premium"
	"github.com/PancyStudios/PancyBotGo/internal/commands/reaction"
	"github.com/PancyStudios/PancyBotGo/internal/commands/security"
	"github.com/PancyStudios/PancyBotGo/internal/commands/utils"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

// RegisterAll registers all commands with the Discord client
// Add your command registration calls here
func RegisterAll(client *discord.ExtendedClient) {
	// Utility commands
	utils.RegisterUtilsCommands(client)

	// Music commands
	RegisterMusicCommands(client)

	// Embeds commands
	embeds.RegisterEmbedsCommands(client)

	// Moderation commands (/mod ban, /mod kick, /mod warn, /mod mute)
	mod.RegisterModCommands(client)

	// Premium commands (/premium redeem)
	premium.Register(client)

	// Config commands (/config ...)
	config.Register(client)

	// Developer commands (/dev codegen)
	dev.Register(client)

	// IA commands (/ia createimage)
	ia.Register(client)

	// Fun commands (/fun 8ball)
	fun.Register(client)

	// Reaction commands (/reaccion hug, kiss)
	reaction.RegisterReactionCommands(client)

	// Economy commands
	economy.Register(client)

	// Security commands (/security antibots)
	security.RegisterSecurityCommands(client)
}
