// Package events provides event handlers for message events
package events

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// RegisterMessageEvents registers all message-related event handlers
func RegisterMessageEvents(client *discord.ExtendedClient) {
	client.Session.AddHandler(onMessageCreate)
	client.Session.AddHandler(onMessageUpdate)
	client.Session.AddHandler(onMessageDelete)
}

// onMessageCreate is called when a new message is created
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignorar mensajes de bots
	if m.Author.Bot {
		return
	}

	// Log solo en modo debug (puede ser spam)
	// logger.Debug(fmt.Sprintf("💬 %s: %s", m.Author.Username, m.Content), "Message")

	// Responder a menciones del bot
	for _, mention := range m.Mentions {
		if mention.ID == s.State.User.ID {
			embed := &discordgo.MessageEmbed{
				Title:       "👋 ¡Hola!",
				Description: "Usa comandos **slash (/)** para interactuar conmigo.\nEscribe `/help` para ver todos los comandos disponibles.",
				Color:       0x3498db,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "🎵 Música",
						Value:  "`/play` - Reproduce música",
						Inline: true,
					},
					{
						Name:   "🔧 Moderación",
						Value:  "`/mod` - Comandos de moderación",
						Inline: true,
					},
					{
						Name:   "❓ Ayuda",
						Value:  "`/help` - Ver todos los comandos",
						Inline: true,
					},
				},
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando respuesta: %v", err), "Message")
			}
			break
		}
	}

	// Respuestas automáticas (ejemplos)
	content := strings.ToLower(m.Content)

	if strings.Contains(content, "hola bot") || strings.Contains(content, "hola pancybot") {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("¡Hola <@%s>! 👋 ¿En qué puedo ayudarte?", m.Author.ID))
		if err != nil {
			logger.Error(fmt.Sprintf("Error enviando saludo: %v", err), "Message")
		}
	}

	if strings.Contains(content, "buenas noches bot") {
		_, err := s.ChannelMessageSend(m.ChannelID, "¡Buenas noches! 🌙 Que descanses.")
		if err != nil {
			logger.Error(fmt.Sprintf("Error enviando despedida: %v", err), "Message")
		}
	}

	if strings.Contains(content, "gracias bot") || strings.Contains(content, "gracias pancybot") {
		_, err := s.ChannelMessageSend(m.ChannelID, "¡De nada! 😊 Siempre es un placer ayudar.")
		if err != nil {
			logger.Error(fmt.Sprintf("Error enviando agradecimiento: %v", err), "Message")
		}
	}

	// Easter eggs: reaccionar a palabras específicas
	if strings.Contains(content, "🎵") || strings.Contains(content, "música") || strings.Contains(content, "music") {
		err := s.MessageReactionAdd(m.ChannelID, m.ID, "🎵")
		if err != nil {
			logger.Debug(fmt.Sprintf("Error agregando reacción: %v", err), "Message")
		}
	}

	if strings.Contains(content, "❤️") || strings.Contains(content, "♥️") {
		err := s.MessageReactionAdd(m.ChannelID, m.ID, "❤️")
		if err != nil {
			logger.Debug(fmt.Sprintf("Error agregando reacción: %v", err), "Message")
		}
	}
}

// onMessageUpdate is called when a message is edited
func onMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author != nil && !m.Author.Bot {
		logger.Debug(fmt.Sprintf("✏️ Mensaje editado por %s en canal %s",
			m.Author.Username, m.ChannelID), "Message")
	}
}

// onMessageDelete is called when a message is deleted
func onMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	logger.Debug(fmt.Sprintf("🗑️ Mensaje eliminado: ID %s en canal %s",
		m.ID, m.ChannelID), "Message")
}
