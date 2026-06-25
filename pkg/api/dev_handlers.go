package api

import (
	"errors"
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/PancyStudios/PancyBotGo/pkg/mqtt"
	"github.com/bwmarrin/discordgo"
)

// RegisterDevHandlers registra los handlers MQTT exclusivos para el panel de developer.
// Todos los tópicos usan el prefijo "dev-" para distinguirlos de los handlers normales.
func RegisterDevHandlers(mc *mqtt.MqttCommunicator, discordClient *discord.ExtendedClient) {

	// ─────────────────────────────────────────────────────────────────────────
	// get-all-bot-guilds
	// Devuelve la lista completa de servidores donde está el bot (sin filtros).
	// ─────────────────────────────────────────────────────────────────────────
	mc.On("get-all-bot-guilds", func(payload map[string]interface{}) (interface{}, error) {
		if discordClient == nil || discordClient.Session == nil || discordClient.Session.State == nil {
			return []interface{}{}, nil
		}

		discordClient.Session.State.RLock()
		defer discordClient.Session.State.RUnlock()

		type GuildSummary struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Icon        string `json:"icon"`
			MemberCount int    `json:"memberCount"`
		}

		result := make([]GuildSummary, 0, len(discordClient.Session.State.Guilds))
		for _, g := range discordClient.Session.State.Guilds {
			result = append(result, GuildSummary{
				ID:          g.ID,
				Name:        g.Name,
				Icon:        g.Icon,
				MemberCount: g.MemberCount,
			})
		}

		logger.Info(fmt.Sprintf("[Dev] get-all-bot-guilds: devolviendo %d servidores", len(result)), "DevMQTT")
		return result, nil
	})

	// ─────────────────────────────────────────────────────────────────────────
	// dev-leave-guild
	// Hace que el bot abandone un servidor específico.
	// Payload: { "guildId": "..." }
	// ─────────────────────────────────────────────────────────────────────────
	mc.On("dev-leave-guild", func(payload map[string]interface{}) (interface{}, error) {
		if discordClient == nil || discordClient.Session == nil {
			return nil, fmt.Errorf("discord client not ready")
		}

		guildID, ok := payload["guildId"].(string)
		if !ok || guildID == "" {
			return nil, fmt.Errorf("missing or invalid guildId")
		}

		if err := discordClient.Session.GuildLeave(guildID); err != nil {
			logger.Error(fmt.Sprintf("[Dev] Error saliendo del servidor %s: %v", guildID, err), "DevMQTT")
			return nil, fmt.Errorf("error al salir del servidor: %w", err)
		}

		logger.Info(fmt.Sprintf("[Dev] Bot salió del servidor %s por orden del developer", guildID), "DevMQTT")
		return map[string]interface{}{"success": true, "guildId": guildID}, nil
	})

	// ─────────────────────────────────────────────────────────────────────────
	// dev-blacklist-guild
	// Añade un servidor a la blacklist y hace que el bot lo abandone.
	// Payload: { "guildId": "...", "reason": "..." }
	// ─────────────────────────────────────────────────────────────────────────
	mc.On("dev-blacklist-guild", func(payload map[string]interface{}) (interface{}, error) {
		if discordClient == nil || discordClient.Session == nil {
			return nil, fmt.Errorf("discord client not ready")
		}

		guildID, ok := payload["guildId"].(string)
		if !ok || guildID == "" {
			return nil, fmt.Errorf("missing or invalid guildId")
		}

		reason := "Sin razón especificada."
		if r, ok := payload["reason"].(string); ok && r != "" {
			reason = r
		}

		// Añadir a la blacklist del bot (database + cache)
		_, err := database.AddToBlacklist(guildID, models.BlacklistTypeGuild, reason, "developer-panel")
		if err != nil && !errors.Is(err, database.ErrAlreadyBlacklisted) {
			logger.Error(fmt.Sprintf("[Dev] Error añadiendo %s a blacklist: %v", guildID, err), "DevMQTT")
			return nil, fmt.Errorf("error al añadir a blacklist: %w", err)
		}

		// Notificar en el servidor antes de salir (best-effort)
		guild, stateErr := discordClient.Session.State.Guild(guildID)
		if stateErr == nil {
			// Buscar el canal de sistema del servidor para enviar aviso
			if guild.SystemChannelID != "" {
				embed := &discordgo.MessageEmbed{
					Title:       "🚫 Servidor en Blacklist",
					Description: "Este servidor ha sido añadido a la blacklist por los desarrolladores. El bot se retirará ahora.",
					Color:       0xFF0000,
					Fields: []*discordgo.MessageEmbedField{
						{Name: "Razón", Value: reason},
					},
				}
				discordClient.Session.ChannelMessageSendEmbed(guild.SystemChannelID, embed) //nolint:errcheck
			}
		}

		// Salir del servidor
		if leaveErr := discordClient.Session.GuildLeave(guildID); leaveErr != nil {
			logger.Warn(fmt.Sprintf("[Dev] No se pudo salir del servidor blacklisted %s: %v", guildID, leaveErr), "DevMQTT")
			// No retornamos error — el blacklist ya fue guardado
		}

		logger.Info(fmt.Sprintf("[Dev] Servidor %s blacklisted y abandonado. Razón: %s", guildID, reason), "DevMQTT")
		return map[string]interface{}{"success": true, "guildId": guildID}, nil
	})

	// ─────────────────────────────────────────────────────────────────────────
	// dev-unblacklist-guild
	// Remueve un servidor de la blacklist.
	// Payload: { "guildId": "..." }
	// ─────────────────────────────────────────────────────────────────────────
	mc.On("dev-unblacklist-guild", func(payload map[string]interface{}) (interface{}, error) {
		guildID, ok := payload["guildId"].(string)
		if !ok || guildID == "" {
			return nil, fmt.Errorf("missing or invalid guildId")
		}

		if err := database.RemoveFromBlacklist(guildID); err != nil {
			if errors.Is(err, database.ErrBlacklistNotFound) {
				// No está en la blacklist del bot — OK, puede haber sido solo en la DB del API
				logger.Info(fmt.Sprintf("[Dev] Servidor %s no estaba en la blacklist del bot, ignorando.", guildID), "DevMQTT")
				return map[string]interface{}{"success": true, "guildId": guildID, "note": "not in bot blacklist"}, nil
			}
			logger.Error(fmt.Sprintf("[Dev] Error removiendo %s de blacklist: %v", guildID, err), "DevMQTT")
			return nil, fmt.Errorf("error al remover de blacklist: %w", err)
		}

		logger.Info(fmt.Sprintf("[Dev] Servidor %s removido de la blacklist del bot", guildID), "DevMQTT")
		return map[string]interface{}{"success": true, "guildId": guildID}, nil
	})

	// ─────────────────────────────────────────────────────────────────────────
	// dev-send-message
	// Envía un mensaje de texto a un canal específico de un servidor.
	// Payload: { "guildId": "...", "channelId": "...", "content": "..." }
	// ─────────────────────────────────────────────────────────────────────────
	mc.On("dev-send-message", func(payload map[string]interface{}) (interface{}, error) {
		if discordClient == nil || discordClient.Session == nil {
			return nil, fmt.Errorf("discord client not ready")
		}

		guildID, ok := payload["guildId"].(string)
		if !ok || guildID == "" {
			return nil, fmt.Errorf("missing or invalid guildId")
		}

		channelID, ok := payload["channelId"].(string)
		if !ok || channelID == "" {
			return nil, fmt.Errorf("missing or invalid channelId")
		}

		content, ok := payload["content"].(string)
		if !ok || content == "" {
			return nil, fmt.Errorf("missing or invalid content")
		}

		// Verificar que el canal pertenece al servidor esperado (seguridad básica)
		channel, err := discordClient.Session.Channel(channelID)
		if err != nil {
			return nil, fmt.Errorf("canal no encontrado: %w", err)
		}
		if channel.GuildID != guildID {
			return nil, fmt.Errorf("el canal no pertenece al servidor indicado")
		}

		msg, err := discordClient.Session.ChannelMessageSend(channelID, content)
		if err != nil {
			logger.Error(fmt.Sprintf("[Dev] Error enviando mensaje al canal %s: %v", channelID, err), "DevMQTT")
			return nil, fmt.Errorf("error al enviar mensaje: %w", err)
		}

		logger.Info(fmt.Sprintf("[Dev] Mensaje enviado al canal %s del servidor %s (msgID: %s)", channelID, guildID, msg.ID), "DevMQTT")
		return map[string]interface{}{
			"success":   true,
			"messageId": msg.ID,
			"channelId": channelID,
			"guildId":   guildID,
		}, nil
	})
}
