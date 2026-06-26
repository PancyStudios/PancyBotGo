// Package events provides event handlers for message events
package events

import (
	"fmt"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
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

	// ------------------ SISTEMA DE NIVELES ------------------
	if m.GuildID != "" {
		handleUserLeveling(s, m)
	}
}

func handleUserLeveling(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Verificar si el servidor tiene activado el sistema
	guildData, err := database.GlobalGuildDM.Get(bson.M{"_id": m.GuildID})
	if err != nil || guildData == nil || !guildData.Levels.Enable {
		return
	}

	profile, err := database.GetLocalLevelProfile(m.GuildID, m.Author.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error obteniendo perfil de nivel para %s: %v", m.Author.ID, err), "Levels")
		return
	}

	now := time.Now()

	// Comprobar si está en enfriamiento
	if now.Before(profile.CooldownUntil) {
		return
	}

	// Comprobar ventana de spam (4 mensajes en 3 segundos)
	if now.Sub(profile.SpamWindowStart) > 3*time.Second {
		// Resetear la ventana
		profile.SpamWindowStart = now
		profile.SpamCount = 1
	} else {
		profile.SpamCount++
		if profile.SpamCount >= 4 {
			// Activar cooldown de 5 segundos
			profile.CooldownUntil = now.Add(5 * time.Second)
			profile.SpamCount = 0 // Resetear cuenta para después del cooldown

			// Guardar el perfil para que el cooldown haga efecto, y no dar XP
			_, err = database.LocalLevelsDM.Set(bson.M{"_id": profile.ID}, profile)
			if err != nil {
				logger.Error(fmt.Sprintf("Error guardando cooldown para %s: %v", m.Author.ID, err), "Levels")
			}
			return
		}
	}

	// Añadir XP aleatorio (1 a 15)
	addedXP := int64(1 + (now.UnixNano() % 15)) // simple pseudo-random
	profile.XP += addedXP
	profile.TotalMessages += 1
	profile.LastMessageTime = now

	// Verificar si subió de nivel
	nextLevel := profile.Level + 1
	requiredXP := nextLevel * nextLevel * 100

	levelUp := false
	if profile.XP >= requiredXP {
		profile.Level = nextLevel
		levelUp = true
	}

	_, err = database.LocalLevelsDM.Set(bson.M{"_id": profile.ID}, profile)
	if err != nil {
		logger.Error(fmt.Sprintf("Error guardando perfil de nivel para %s: %v", m.Author.ID, err), "Levels")
		return
	}

	// Enviar mensaje de Level Up si es necesario
	if levelUp {
		// Asignar roles de recompensa
		for _, reward := range guildData.Levels.Rewards {
			if reward.Level == profile.Level {
				err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, reward.RoleID)
				if err != nil {
					logger.Error(fmt.Sprintf("No se pudo asignar rol de nivel %d a %s: %v", profile.Level, m.Author.ID, err), "Levels")
				} else {
					logger.Info(fmt.Sprintf("Rol %s asignado a %s por alcanzar el nivel %d", reward.RoleID, m.Author.ID, profile.Level), "Levels")
				}
			}
		}

		chID := m.ChannelID
		if guildData.Levels.LevelUpChannel != "" {
			chID = guildData.Levels.LevelUpChannel
		}

		msgContent := guildData.Levels.LevelUpMessage
		if msgContent == "" {
			msgContent = "¡Felicidades {user}, has avanzado al **Nivel {level}**! 🎉"
		}

		msgContent = strings.ReplaceAll(msgContent, "{user}", fmt.Sprintf("<@%s>", m.Author.ID))
		msgContent = strings.ReplaceAll(msgContent, "{level}", fmt.Sprintf("%d", profile.Level))

		_, err = s.ChannelMessageSend(chID, msgContent)
		if err != nil {
			logger.Error(fmt.Sprintf("No se pudo enviar mensaje de level up a %s: %v", chID, err), "Levels")
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
