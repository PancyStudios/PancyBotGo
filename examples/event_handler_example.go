// Package examples demonstrates how to use the Event Handler
// Copy this code to your main.go to add event handling
//
// NOTE: This example uses the helper methods (OnReady, OnGuildCreate, etc.)
// In the new modular system, use client.Session.AddHandler(handlerFunc) directly
// See internal/events/ for the correct implementation
package examples

import (
	"fmt"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// setupEvents configures all event handlers for the bot
func setupEvents(client *discord.ExtendedClient) {
	logger.System("Configurando eventos del bot...", "Events")

	// ============================================
	// Evento: Bot Ready
	// ============================================
	client.EventHandler.OnReady(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Success(fmt.Sprintf("✅ Bot conectado: %s#%s", r.User.Username, r.User.Discriminator), "Events")
		logger.Info(fmt.Sprintf("📊 Conectado a %d servidores", len(r.Guilds)), "Events")

		// Opcional: Establecer estado del bot
		err := s.UpdateGameStatus(0, "🎵 Música con /play")
		if err != nil {
			logger.Error(fmt.Sprintf("Error estableciendo estado: %v", err), "Events")
		}
	})

	// ============================================
	// Evento: Nuevo Servidor (Guild Create)
	// ============================================
	client.EventHandler.OnGuildCreate(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		logger.Info(fmt.Sprintf("➕ Bot agregado a servidor: %s (ID: %s)", g.Name, g.ID), "Events")
		logger.Debug(fmt.Sprintf("   Miembros: %d | Canales: %d", g.MemberCount, len(g.Channels)), "Events")

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
						Name:   "⚙️ Configuración",
						Value:  "Configura el bot con `/settings`",
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "¡Disfruta!",
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			_, err := s.ChannelMessageSendEmbed(g.SystemChannelID, welcomeEmbed)
			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando mensaje de bienvenida: %v", err), "Events")
			}
		}
	})

	// ============================================
	// Evento: Servidor Removido (Guild Delete)
	// ============================================
	client.EventHandler.OnGuildDelete(func(s *discordgo.Session, g *discordgo.GuildDelete) {
		logger.Info(fmt.Sprintf("➖ Bot removido del servidor ID: %s", g.ID), "Events")
	})

	// ============================================
	// Evento: Nuevo Miembro (Guild Member Add)
	// ============================================
	client.EventHandler.OnGuildMemberAdd(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		logger.Info(fmt.Sprintf("👋 Nuevo miembro: %s#%s en servidor %s",
			m.User.Username, m.User.Discriminator, m.GuildID), "Events")

		// Obtener información del servidor
		guild, err := s.Guild(m.GuildID)
		if err != nil {
			logger.Error(fmt.Sprintf("Error obteniendo servidor: %v", err), "Events")
			return
		}

		// Enviar mensaje de bienvenida al canal del sistema
		if guild.SystemChannelID != "" {
			welcomeMsg := fmt.Sprintf("¡Bienvenido/a <@%s> a **%s**! 🎉\nAhora somos **%d** miembros.",
				m.User.ID, guild.Name, guild.MemberCount)

			_, err := s.ChannelMessageSend(guild.SystemChannelID, welcomeMsg)
			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando mensaje de bienvenida: %v", err), "Events")
			}
		}

		// Opcional: Enviar DM de bienvenida
		channel, err := s.UserChannelCreate(m.User.ID)
		if err == nil {
			welcomeEmbed := &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("¡Bienvenido/a a %s!", guild.Name),
				Description: "Esperamos que disfrutes tu estancia. Si necesitas ayuda, no dudes en preguntar.",
				Color:       0x00ff00,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: guild.IconURL("256"),
				},
			}
			_, dmErr := s.ChannelMessageSendEmbed(channel.ID, welcomeEmbed)
			if dmErr != nil {
				logger.Debug("No se pudo enviar DM de bienvenida (DMs cerrados)", "Events")
			}
		}

		// Opcional: Asignar rol automático
		// autoRoleID := "TU_ROL_ID_AQUI"
		// err = s.GuildMemberRoleAdd(m.GuildID, m.User.ID, autoRoleID)
		// if err != nil {
		//     logger.Error(fmt.Sprintf("Error asignando rol: %v", err), "Events")
		// }
	})

	// ============================================
	// Evento: Miembro Sale (Guild Member Remove)
	// ============================================
	client.EventHandler.OnGuildMemberRemove(func(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
		logger.Info(fmt.Sprintf("👋 Adiós: %s#%s salió del servidor %s",
			m.User.Username, m.User.Discriminator, m.GuildID), "Events")

		// Enviar mensaje de despedida
		guild, err := s.Guild(m.GuildID)
		if err == nil && guild.SystemChannelID != "" {
			farewellMsg := fmt.Sprintf("👋 **%s#%s** ha salido del servidor. ¡Hasta pronto!",
				m.User.Username, m.User.Discriminator)
			s.ChannelMessageSend(guild.SystemChannelID, farewellMsg)
		}
	})

	// ============================================
	// Evento: Estado de Voz (Voice State Update)
	// ============================================
	client.EventHandler.OnVoiceStateUpdate(func(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
		// Usuario se unió a un canal de voz
		if v.ChannelID != "" && (v.BeforeUpdate == nil || v.BeforeUpdate.ChannelID == "") {
			channel, err := s.Channel(v.ChannelID)
			if err == nil {
				logger.Debug(fmt.Sprintf("🎤 Usuario %s se unió a: %s", v.UserID, channel.Name), "Events")
			}
		}

		// Usuario salió de un canal de voz
		if v.ChannelID == "" && v.BeforeUpdate != nil && v.BeforeUpdate.ChannelID != "" {
			logger.Debug(fmt.Sprintf("🔇 Usuario %s salió del canal de voz", v.UserID), "Events")
		}

		// Usuario cambió de canal de voz
		if v.ChannelID != "" && v.BeforeUpdate != nil &&
			v.BeforeUpdate.ChannelID != "" && v.ChannelID != v.BeforeUpdate.ChannelID {
			oldChannel, _ := s.Channel(v.BeforeUpdate.ChannelID)
			newChannel, _ := s.Channel(v.ChannelID)
			if oldChannel != nil && newChannel != nil {
				logger.Debug(fmt.Sprintf("🔄 Usuario %s: %s → %s",
					v.UserID, oldChannel.Name, newChannel.Name), "Events")
			}
		}
	})

	// ============================================
	// Evento: Mensaje Creado (Message Create)
	// ============================================
	client.EventHandler.OnMessageCreate(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignorar mensajes de bots
		if m.Author.Bot {
			return
		}

		// Log de mensajes (solo para debug, puede ser spam)
		// logger.Debug(fmt.Sprintf("💬 %s: %s", m.Author.Username, m.Content), "Events")

		// Responder a menciones
		for _, mention := range m.Mentions {
			if mention.ID == s.State.User.ID {
				embed := &discordgo.MessageEmbed{
					Description: "¡Hola! 👋 Usa comandos **slash (/)** para interactuar conmigo.\nEscribe `/help` para ver todos los comandos disponibles.",
					Color:       0x3498db,
				}
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
				break
			}
		}

		// Respuestas automáticas (ejemplos)
		content := strings.ToLower(m.Content)

		if strings.Contains(content, "hola bot") || strings.Contains(content, "hola pancybot") {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("¡Hola <@%s>! 👋", m.Author.ID))
		}

		if strings.Contains(content, "buenas noches bot") {
			s.ChannelMessageSend(m.ChannelID, "¡Buenas noches! 🌙 Que descanses.")
		}

		// Easter egg: reaccionar a palabras específicas
		if strings.Contains(content, "🎵") || strings.Contains(content, "música") {
			s.MessageReactionAdd(m.ChannelID, m.ID, "🎵")
		}
	})

	// ============================================
	// Evento: Mensaje Editado (Message Update)
	// ============================================
	client.EventHandler.OnMessageUpdate(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		if m.Author != nil && !m.Author.Bot {
			logger.Debug(fmt.Sprintf("✏️ Mensaje editado por %s en canal %s",
				m.Author.Username, m.ChannelID), "Events")
		}
	})

	// ============================================
	// Evento: Mensaje Eliminado (Message Delete)
	// ============================================
	client.EventHandler.OnMessageDelete(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		logger.Debug(fmt.Sprintf("🗑️ Mensaje eliminado: ID %s en canal %s",
			m.ID, m.ChannelID), "Events")
	})

	// ============================================
	// Evento: Miembro Actualizado (Guild Member Update)
	// ============================================
	client.EventHandler.OnGuildMemberUpdate(func(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
		// Detectar cambio de nickname
		if m.BeforeUpdate != nil {
			oldNick := m.BeforeUpdate.Nick
			newNick := m.Nick

			if oldNick != newNick {
				logger.Debug(fmt.Sprintf("✏️ %s cambió nickname: '%s' → '%s'",
					m.User.Username, oldNick, newNick), "Events")
			}

			// Detectar cambio de roles
			if len(m.BeforeUpdate.Roles) != len(m.Roles) {
				logger.Debug(fmt.Sprintf("🎭 Roles actualizados para %s", m.User.Username), "Events")
			}
		}
	})

	// ============================================
	// Evento: Interacción (Interaction Create)
	// ============================================
	// Nota: El CommandHandler ya maneja slash commands
	// Este es útil para botones, menús, modales, etc.
	client.EventHandler.OnInteractionCreate(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Manejar componentes de mensaje (botones, menús)
		if i.Type == discordgo.InteractionMessageComponent {
			customID := i.MessageComponentData().CustomID
			logger.Debug(fmt.Sprintf("🔘 Componente clickeado: %s", customID), "Events")

			// Ejemplo: Manejar botones
			switch customID {
			case "button_example":
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "¡Botón clickeado! ✅",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}
		}
	})

	logger.Success("✅ Eventos configurados correctamente", "Events")
}

// setupAuditLog configura un sistema de logs de auditoría
// Requiere un canal de logs en el servidor
func setupAuditLog(client *discord.ExtendedClient, logChannelID string) {
	if logChannelID == "" {
		logger.Warn("No se configuró canal de audit log", "Events")
		return
	}

	logger.System("Configurando sistema de audit log...", "Events")

	// Log: Nuevo miembro
	client.EventHandler.OnGuildMemberAdd(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		embed := &discordgo.MessageEmbed{
			Title:       "👤 Nuevo Miembro",
			Description: fmt.Sprintf("<@%s> se unió al servidor", m.User.ID),
			Color:       0x00ff00,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Usuario",
					Value:  fmt.Sprintf("%s#%s", m.User.Username, m.User.Discriminator),
					Inline: true,
				},
				{
					Name:   "ID",
					Value:  m.User.ID,
					Inline: true,
				},
				{
					Name:   "Cuenta creada",
					Value:  m.User.ID, // Puedes calcular la fecha desde el ID
					Inline: true,
				},
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: m.User.AvatarURL("128"),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(logChannelID, embed)
	})

	// Log: Miembro sale
	client.EventHandler.OnGuildMemberRemove(func(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
		embed := &discordgo.MessageEmbed{
			Title:       "👋 Miembro Salió",
			Description: fmt.Sprintf("%s#%s salió del servidor", m.User.Username, m.User.Discriminator),
			Color:       0xff0000,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "ID",
					Value:  m.User.ID,
					Inline: true,
				},
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: m.User.AvatarURL("128"),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(logChannelID, embed)
	})

	// Log: Mensaje eliminado
	client.EventHandler.OnMessageDelete(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		embed := &discordgo.MessageEmbed{
			Title: "🗑️ Mensaje Eliminado",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Canal",
					Value:  fmt.Sprintf("<#%s>", m.ChannelID),
					Inline: true,
				},
				{
					Name:   "ID del Mensaje",
					Value:  m.ID,
					Inline: true,
				},
			},
			Color:     0xffa500,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(logChannelID, embed)
	})

	logger.Success("✅ Audit log configurado", "Events")
}

// Ejemplo de uso en main.go:
/*
func main() {
	// ... inicialización del bot ...

	// Configurar eventos ANTES de iniciar el bot
	setupEvents(discordClient)

	// Opcional: Configurar audit log
	// setupAuditLog(discordClient, "ID_DEL_CANAL_DE_LOGS")

	// Iniciar bot
	if err := discordClient.Start(); err != nil {
		logger.Critical(fmt.Sprintf("Error: %v", err), "Main")
		os.Exit(1)
	}

	// ... resto del código ...
}
*/
