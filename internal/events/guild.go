// Package events provides event handlers for guild (server) events
package events

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// RegisterGuildEvents registers all guild-related event handlers
func RegisterGuildEvents(client *discord.ExtendedClient) {
	client.Session.AddHandler(onGuildCreate)
	client.Session.AddHandler(onGuildDelete)
}

// onGuildCreate is called when the bot joins a server
func onGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {

	Join := g.JoinedAt
	if Join.Compare(time.Now().Add(-10*time.Second)) < 0 {
		return
	}

	logger.Info(fmt.Sprintf("➕ Bot agregado a servidor: %s (ID: %s)", g.Name, g.ID), "Guild")
	logger.Debug(fmt.Sprintf("   Miembros: %d | Canales: %d", g.MemberCount, len(g.Channels)), "Guild")

	// Enviar mensaje de bienvenida al canal del sistema

	if g.SystemChannelID != "" {
		welcomeEmbed := &discordgo.MessageEmbed{
			Title:       "¡Gracias por agregarme! 🎉",
			Description: "Hola, soy **PancyBot**. Usa `/help` para ver todos mis comandos.",
			Color:       0x00ff00,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "🎵 Música",
					Value:  "Reproduce música con `/play`",
					Inline: true,
				},
				{
					Name:   "🔧 Moderación",
					Value:  "Usa `/mod` para moderar",
					Inline: true,
				},
				{
					Name:   "❓ Ayuda",
					Value:  "Usa `/help` para más información",
					Inline: true,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "¡Disfruta de PancyBot!",
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		_, err := s.ChannelMessageSendEmbed(g.SystemChannelID, welcomeEmbed)
		if err != nil {
			logger.Error(fmt.Sprintf("Error enviando mensaje de bienvenida: %v", err), "Guild")
		}
	}
}

// onGuildDelete is called when the bot is removed from a server
func onGuildDelete(s *discordgo.Session, g *discordgo.GuildDelete) {
	logger.Info(fmt.Sprintf("➖ Bot removido del servidor ID: %s", g.ID), "Guild")
}
