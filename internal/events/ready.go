// Package events provides event handlers for the bot
package events

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// RegisterReadyEvent registers the ready event handler
func RegisterReadyEvent(client *discord.ExtendedClient) {
	client.Session.AddHandler(onReady)
	client.Session.AddHandler(onDebug)
}

// onReady is called when the bot successfully connects to Discord
func onReady(s *discordgo.Session, r *discordgo.Ready) {
	logger.Success(fmt.Sprintf("✅ Bot conectado: %s#%s", r.User.Username, r.User.Discriminator), "Ready")
	logger.Info(fmt.Sprintf("📊 Conectado a %d servidores", len(r.Guilds)), "Ready")

	// Establecer estado del bot
	err := s.UpdateGameStatus(0, "🎵 Música con /play")
	if err != nil {
		logger.Error(fmt.Sprintf("Error estableciendo estado: %v", err), "Ready")
		return
	}

	logger.Debug("Estado del bot establecido correctamente", "Ready")
}

func onDebug(s *discordgo.Session, log string) {
	logger.Debug(log, "DiscordGO")
}
