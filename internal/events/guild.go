// Package events provides event handlers for guild (server) events
package events

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
	"github.com/goccy/go-json"
)

// checkGuildBlacklist verifica si un guild está en la blacklist
func checkGuildBlacklist(guildID string) (bool, *models.Blacklist) {
	return database.IsGuildBlacklisted(guildID)
}

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

	// Verificar si el servidor está en la blacklist
	go func() {
		defer errors.RecoverMiddleware()()

		// Importar database aquí para verificar blacklist
		if isBlacklisted, entry := checkGuildBlacklist(g.ID); isBlacklisted {
			logger.Warn(fmt.Sprintf("🚫 Bot agregado a servidor blacklisted: %s (ID: %s). Saliendo...", g.Name, g.ID), "Guild")

			// Intentar notificar al owner
			if g.OwnerID != "" {
				channel, err := s.UserChannelCreate(g.OwnerID)
				if err == nil {
					embed := &discordgo.MessageEmbed{
						Title:       "🚫 Servidor en Blacklist",
						Description: fmt.Sprintf("El servidor **%s** está en la blacklist. El bot no puede permanecer en este servidor.", g.Name),
						Color:       0xFF0000,
						Timestamp:   time.Now().Format(time.RFC3339),
					}

					if entry != nil && entry.Reason != "" {
						embed.Fields = []*discordgo.MessageEmbedField{
							{
								Name:  "Razón",
								Value: entry.Reason,
							},
						}
					}

					s.ChannelMessageSendEmbed(channel.ID, embed)
				}
			}

			// Salir del servidor
			time.Sleep(2 * time.Second)
			if err := s.GuildLeave(g.ID); err != nil {
				logger.Error(fmt.Sprintf("Error saliendo del servidor blacklisted %s: %v", g.ID, err), "Guild")
			}
			return
		}
	}()

	// Enviar mensaje de bienvenida al canal del sistema
	go func() {
		errors.RecoverMiddleware()()
		if g.SystemChannelID != "" {
			welcomeEmbed := &discordgo.MessageEmbed{
				Title:       "¡Gracias por agregarme! 🎉",
				Description: "Hola, soy **PancyBot**. Usa `/utils help` para ver todos mis comandos.",
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
						Value:  "Usa `/utils help` para más información",
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

		webhook := config.Get().GuildsWebhook
		if webhook == "" {
			return
		}
		fechaCreacion, err := discordgo.SnowflakeTimestamp(g.ID)
		if err != nil {
			log.Println("Error obteniendo fecha:", err)
			return
		}

		embed := map[string]interface{}{
			"title":       "➕ Nuevo servidor agregado",
			"description": "El bot ha sido agregado a un nuevo servidor.",
			"color":       0x00ff00,
			"fields": []map[string]string{
				{
					"name":   "Servidor",
					"value":  fmt.Sprintf("%s (%s)", g.Name, g.ID),
					"inline": "true",
				},
				{
					"name":   "Miembros",
					"value":  fmt.Sprintf("%d", g.MemberCount),
					"inline": "true",
				},
				{
					"name":   "Canales",
					"value":  fmt.Sprintf("%d", len(g.Channels)),
					"inline": "true",
				},
				{
					"name":   "Fecha de creación",
					"value":  fechaCreacion.Format(time.RFC850),
					"inline": "true",
				},
			},
			"timestamp": time.Now().Format(time.RFC3339),
			"footer": map[string]string{
				"text": "💫 Developed by PancyStudio | PancyBot Go",
			},
		}

		payload := map[string]interface{}{
			"embeds": []interface{}{embed},
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return
		}

		req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(jsonData))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
	}()
}

// onGuildDelete is called when the bot is removed from a server
func onGuildDelete(s *discordgo.Session, g *discordgo.GuildDelete) {
	logger.Info(fmt.Sprintf("➖ Bot removido del servidor ID: %s", g.ID), "Guild")
}
