package api

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/mqtt"
)

// RegisterAPIHandlers registers all MQTT endpoints for the REST API
func RegisterAPIHandlers(mc *mqtt.MqttCommunicator, discordClient *discord.ExtendedClient) {
	
	// get-bot-guild-ids
	mc.On("get-bot-guild-ids", func(payload map[string]interface{}) (interface{}, error) {
		if discordClient == nil || discordClient.Session == nil || discordClient.Session.State == nil {
			return []string{}, nil
		}

		discordClient.Session.State.RLock()
		defer discordClient.Session.State.RUnlock()

		ids := make([]string, 0, len(discordClient.Session.State.Guilds))
		for _, g := range discordClient.Session.State.Guilds {
			ids = append(ids, g.ID)
		}
		
		return ids, nil
	})

	// get-guild-info
	mc.On("get-guild-info", func(payload map[string]interface{}) (interface{}, error) {
		if discordClient == nil || discordClient.Session == nil || discordClient.Session.State == nil {
			return nil, fmt.Errorf("discord client not ready")
		}

		guildIDInter, ok := payload["guildId"]
		if !ok {
			return nil, fmt.Errorf("missing guildId")
		}
		
		guildID, ok := guildIDInter.(string)
		if !ok {
			return nil, fmt.Errorf("guildId must be a string")
		}

		guild, err := discordClient.Session.State.Guild(guildID)
		if err != nil {
			return nil, fmt.Errorf("guild not found")
		}

		return map[string]interface{}{
			"id":   guild.ID,
			"name": guild.Name,
			"icon": guild.Icon,
		}, nil
	})
}
