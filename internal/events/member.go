// Package events provides event handlers for member events
package events

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// RegisterMemberEvents registers all member-related event handlers
func RegisterMemberEvents(client *discord.ExtendedClient) {
	client.Session.AddHandler(onGuildMemberAdd)
	client.Session.AddHandler(onGuildMemberRemove)
	client.Session.AddHandler(onGuildMemberUpdate)
}

// onGuildMemberAdd is called when a new member joins the server
func onGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	logger.Info(fmt.Sprintf("👋 Nuevo miembro: %s#%s en servidor %s",
		m.User.Username, m.User.Discriminator, m.GuildID), "Member")

	// Obtener información del servidor
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error obteniendo servidor: %v", err), "Member")
		return
	}

	// Enviar mensaje de bienvenida al canal del sistema
	if guild.SystemChannelID != "" {
		welcomeEmbed := &discordgo.MessageEmbed{
			Title:       "¡Bienvenido/a! 🎉",
			Description: fmt.Sprintf("Dale la bienvenida a <@%s>\nAhora somos **%d** miembros.", m.User.ID, guild.MemberCount),
			Color:       0x00ff00,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: m.User.AvatarURL("128"),
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    guild.Name,
				IconURL: guild.IconURL("64"),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		_, err := s.ChannelMessageSendEmbed(guild.SystemChannelID, welcomeEmbed)
		if err != nil {
			logger.Error(fmt.Sprintf("Error enviando mensaje de bienvenida: %v", err), "Member")
		}
	}

	// Opcional: Enviar DM de bienvenida
	channel, err := s.UserChannelCreate(m.User.ID)
	if err == nil {
		dmEmbed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("¡Bienvenido/a a %s!", guild.Name),
			Description: "Esperamos que disfrutes tu estancia. Si necesitas ayuda, no dudes en preguntar a los administradores.",
			Color:       0x3498db,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: guild.IconURL("256"),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "📋 Reglas",
					Value: "Asegúrate de leer las reglas del servidor",
				},
				{
					Name:  "💬 Preséntate",
					Value: "¡No olvides presentarte a la comunidad!",
				},
			},
		}

		_, dmErr := s.ChannelMessageSendEmbed(channel.ID, dmEmbed)
		if dmErr != nil {
			logger.Debug("No se pudo enviar DM de bienvenida (DMs cerrados)", "Member")
		}
	}

	// Opcional: Asignar rol automático
	// autoRoleID := "TU_ROL_ID_AQUI"
	// err = s.GuildMemberRoleAdd(m.GuildID, m.User.ID, autoRoleID)
	// if err != nil {
	//     logger.Error(fmt.Sprintf("Error asignando rol: %v", err), "Member")
	// }
}

// onGuildMemberRemove is called when a member leaves the server
func onGuildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	logger.Info(fmt.Sprintf("👋 Adiós: %s#%s salió del servidor %s",
		m.User.Username, m.User.Discriminator, m.GuildID), "Member")

	// Enviar mensaje de despedida
	guild, err := s.Guild(m.GuildID)
	if err == nil && guild.SystemChannelID != "" {
		farewellEmbed := &discordgo.MessageEmbed{
			Description: fmt.Sprintf("👋 **%s#%s** ha salido del servidor.\nAhora somos **%d** miembros.",
				m.User.Username, m.User.Discriminator, guild.MemberCount),
			Color: 0xe74c3c,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: m.User.AvatarURL("64"),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		_, sendErr := s.ChannelMessageSendEmbed(guild.SystemChannelID, farewellEmbed)
		if sendErr != nil {
			logger.Error(fmt.Sprintf("Error enviando mensaje de despedida: %v", sendErr), "Member")
		}
	}
}

// onGuildMemberUpdate is called when a member is updated (roles, nickname, etc.)
func onGuildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	// Solo loguear si hay cambios significativos
	if m.BeforeUpdate != nil {
		// Detectar cambio de nickname
		oldNick := m.BeforeUpdate.Nick
		newNick := m.Nick

		if oldNick != newNick {
			logger.Debug(fmt.Sprintf("✏️ %s cambió nickname: '%s' → '%s'",
				m.User.Username, oldNick, newNick), "Member")
		}

		// Detectar cambio de roles
		if len(m.BeforeUpdate.Roles) != len(m.Roles) {
			logger.Debug(fmt.Sprintf("🎭 Roles actualizados para %s", m.User.Username), "Member")
		}
	}
}
