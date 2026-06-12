package lavalink

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/pkg/mqtt"
)

// RegisterMusicHandlers registers all MQTT endpoints for music control
func RegisterMusicHandlers(mc *mqtt.MqttCommunicator, llClient *LavalinkClient) {

	// PLAY / RESUME
	mc.On("music/+/play", func(payload map[string]interface{}) (interface{}, error) {
		actualTopic := payload["_topic"].(string)
		parts := strings.Split(actualTopic, "/")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid topic structure")
		}
		guildID := parts[1]

		player := llClient.GetPlayer(guildID)
		player.Mu.RLock()
		isPaused := player.IsPaused
		player.Mu.RUnlock()

		if isPaused {
			if err := llClient.Pause(guildID, false); err != nil {
				return nil, err
			}
		}

		response := map[string]interface{}{
			"success": true,
			"message": "Reproducción reanudada",
		}

		player.Mu.RLock()
		if player.CurrentTrack != nil {
			response["track"] = player.CurrentTrack.Info.Title
		}
		player.Mu.RUnlock()

		return response, nil
	})

	// PAUSE
	mc.On("music/+/pause", func(payload map[string]interface{}) (interface{}, error) {
		actualTopic := payload["_topic"].(string)
		parts := strings.Split(actualTopic, "/")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid topic structure")
		}
		guildID := parts[1]

		if err := llClient.Pause(guildID, true); err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"success": true,
			"message": "Reproducción pausada",
		}, nil
	})

	// SKIP (Next or Previous)
	mc.On("music/+/skip", func(payload map[string]interface{}) (interface{}, error) {
		actualTopic := payload["_topic"].(string)
		parts := strings.Split(actualTopic, "/")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid topic structure")
		}
		guildID := parts[1]

		direction := "next"
		if d, ok := payload["direction"].(string); ok {
			direction = d
		}

		if direction == "next" {
			if err := llClient.Skip(guildID); err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"success": true,
				"message": "Saltado a la siguiente canción",
			}, nil
		}

		return nil, fmt.Errorf("skip to previous not fully implemented yet")
	})

	// SKIP TO INDEX
	mc.On("music/+/skip/+", func(payload map[string]interface{}) (interface{}, error) {
		actualTopic := payload["_topic"].(string)
		parts := strings.Split(actualTopic, "/")
		if len(parts) < 4 {
			return nil, fmt.Errorf("invalid topic structure")
		}
		guildID := parts[1]

		var index float64 // json decodes numbers as float64
		if idx, ok := payload["index"].(float64); ok {
			index = idx
		}

		player := llClient.GetPlayer(guildID)
		player.Mu.Lock()
		defer player.Mu.Unlock()

		if int(index) < 0 || int(index) >= len(player.Queue) {
			return nil, fmt.Errorf("index out of bounds")
		}

		player.Queue = player.Queue[int(index):]

		player.Mu.Unlock()
		err := llClient.Skip(guildID)
		player.Mu.Lock()

		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Saltado al índice %d", int(index)),
		}, nil
	})
}
