package api

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/mqtt"
	"go.mongodb.org/mongo-driver/bson"
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

	// update-guild-cache
	mc.On("update-guild-cache", func(payload map[string]interface{}) (interface{}, error) {
		guildID, _ := payload["guildId"].(string)
		if guildID != "" {
			database.GlobalGuildDM.ClearCache()
		}
		return map[string]interface{}{"success": true}, nil
	})

	// verify-user-web
	mc.On("verify-user-web", func(payload map[string]interface{}) (interface{}, error) {
		guildID, _ := payload["guildId"].(string)
		userID, _ := payload["userId"].(string)

		if guildID == "" || userID == "" {
			return nil, fmt.Errorf("missing guildId or userId")
		}

		guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": guildID})
		if err != nil || guildDoc == nil || !guildDoc.Protection.Verification.Enable || guildDoc.Protection.Verification.Role == "" {
			return nil, fmt.Errorf("verification disabled or role not configured")
		}

		err = discordClient.Session.GuildMemberRoleAdd(guildID, userID, guildDoc.Protection.Verification.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to add role: %v", err)
		}

		return map[string]interface{}{"success": true}, nil
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

		var roles []map[string]interface{}
		for _, r := range guild.Roles {
			roles = append(roles, map[string]interface{}{
				"id":    r.ID,
				"name":  r.Name,
				"color": r.Color,
			})
		}

		var channels []map[string]interface{}
		for _, c := range guild.Channels {
			// Solamente canales de texto (Type 0)
			if c.Type == 0 {
				channels = append(channels, map[string]interface{}{
					"id":   c.ID,
					"name": c.Name,
				})
			}
		}

		return map[string]interface{}{
			"id":       guild.ID,
			"name":     guild.Name,
			"icon":     guild.Icon,
			"roles":    roles,
			"channels": channels,
		}, nil
	})

	// get-levels-leaderboard
	mc.On("get-levels-leaderboard", func(payload map[string]interface{}) (interface{}, error) {
		if discordClient == nil || discordClient.Session == nil {
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

		limit := int64(10)
		if limitInter, ok := payload["limit"]; ok {
			if l, ok := limitInter.(float64); ok {
				limit = int64(l)
			}
		}

		profiles, err := database.GetTopLevels(guildID, limit)
		if err != nil {
			return nil, fmt.Errorf("error fetching leaderboard: %w", err)
		}

		type LeaderboardEntry struct {
			UserID        string `json:"userId"`
			Username      string `json:"username"`
			AvatarURL     string `json:"avatarUrl"`
			Level         int64  `json:"level"`
			XP            int64  `json:"xp"`
			TotalMessages int64  `json:"totalMessages"`
		}

		var result []LeaderboardEntry
		for _, p := range profiles {
			username := "Usuario Desconocido"
			avatar := ""

			// Try to get user info from state, fallback to API
			member, err := discordClient.Session.State.Member(guildID, p.UserID)
			if err == nil && member != nil && member.User != nil {
				username = member.User.Username
				avatar = member.User.AvatarURL("")
			} else {
				user, err := discordClient.Session.User(p.UserID)
				if err == nil && user != nil {
					username = user.Username
					avatar = user.AvatarURL("")
				}
			}

			result = append(result, LeaderboardEntry{
				UserID:        p.UserID,
				Username:      username,
				AvatarURL:     avatar,
				Level:         p.Level,
				XP:            p.XP,
				TotalMessages: p.TotalMessages,
			})
		}

		return result, nil
	})

	// get-user-level
	mc.On("get-user-level", func(payload map[string]interface{}) (interface{}, error) {
		if discordClient == nil || discordClient.Session == nil {
			return nil, fmt.Errorf("discord client not ready")
		}

		guildIDInter, ok := payload["guildId"]
		if !ok {
			return nil, fmt.Errorf("missing guildId")
		}

		userIDInter, ok := payload["userId"]
		if !ok {
			return nil, fmt.Errorf("missing userId")
		}

		guildID, ok := guildIDInter.(string)
		if !ok {
			return nil, fmt.Errorf("guildId must be a string")
		}

		userID, ok := userIDInter.(string)
		if !ok {
			return nil, fmt.Errorf("userId must be a string")
		}

		profile, err := database.GetLocalLevelProfile(guildID, userID)
		if err != nil {
			return nil, fmt.Errorf("error fetching user level: %w", err)
		}

		nextLevel := profile.Level + 1
		requiredXP := nextLevel * nextLevel * 100

		return map[string]interface{}{
			"userId":        profile.UserID,
			"level":         profile.Level,
			"xp":            profile.XP,
			"requiredXp":    requiredXP,
			"totalMessages": profile.TotalMessages,
		}, nil
	})

	// get-stats
	mc.On("get-stats", func(payload map[string]interface{}) (interface{}, error) {
		if discordClient == nil || discordClient.Session == nil || discordClient.Session.State == nil {
			return nil, fmt.Errorf("discord client not ready")
		}

		discordClient.Session.State.RLock()
		guilds := len(discordClient.Session.State.Guilds)
		users := 0
		channels := 0
		for _, g := range discordClient.Session.State.Guilds {
			users += g.MemberCount
			channels += len(g.Channels)
		}
		discordClient.Session.State.RUnlock()

		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		ping := discordClient.Session.HeartbeatLatency()
		uptime := time.Since(discordClient.StartTime).Milliseconds()

		commandsSize := 0
		if discordClient.Commands != nil {
			commandsSize = discordClient.Commands.Size()
		}

		dbStatus := "Desconectado"
		if database.Get() != nil && database.Get().Connected() {
			dbStatus = "Conectado"
		}

		stats := map[string]interface{}{
			"guilds":      guilds,
			"users":       users,
			"channels":    channels,
			"uptime":      uptime,
			"commands":    commandsSize,
			"nodeVersion": runtime.Version(),
			"botVersion":  "1.0.0",
			"ping":        fmt.Sprintf("%dms", ping.Milliseconds()),
			"memory":      float64(memStats.Alloc) / 1024 / 1024,
			"cpu":         0,
			"platform":    runtime.GOOS,
			"arch":        runtime.GOARCH,
			"release":     "PancyBotGo",
			"database":    dbStatus,
		}

		bytes, err := json.Marshal(stats)
		if err != nil {
			return nil, err
		}

		return string(bytes), nil
	})
}
