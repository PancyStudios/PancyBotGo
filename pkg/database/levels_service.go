package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var levelsMu sync.Mutex

type LevelMessageResult struct {
	Profile *models.UserLevelProfile
	LevelUp bool
	Skipped bool
}

// GetLocalLevelProfile retrieves the user's level profile or creates a new one
func GetLocalLevelProfile(guildID, userID string) (*models.UserLevelProfile, error) {
	id := fmt.Sprintf("%s_%s", guildID, userID)
	query := bson.M{"_id": id}

	profile, err := LocalLevelsDM.Get(query)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		now := time.Now()
		newProfile := &models.UserLevelProfile{
			ID:              id,
			GuildID:         guildID,
			UserID:          userID,
			XP:              0,
			Level:           0,
			TotalMessages:   0,
			LastMessageTime: time.Time{},
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		profile, err = LocalLevelsDM.Set(query, newProfile)
		if err != nil {
			return nil, err
		}
	}

	return profile, nil
}

// GetTopLevels retrieves the top users by XP in a guild
func GetTopLevels(guildID string, limit int64) ([]*models.UserLevelProfile, error) {
	query := bson.M{"guild_id": guildID}
	collection := LocalLevelsDM.collection()
	if collection == nil {
		return nil, ErrDatabaseUnavailable
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find().
		SetSort(bson.D{{Key: "xp", Value: -1}}).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*models.UserLevelProfile
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// ProcessLevelMessage applies spam protection and XP as one serialized profile update.
func ProcessLevelMessage(guildID, userID string, now time.Time) (*LevelMessageResult, error) {
	levelsMu.Lock()
	defer levelsMu.Unlock()

	profile, err := GetLocalLevelProfile(guildID, userID)
	if err != nil {
		return nil, err
	}
	result := &LevelMessageResult{Profile: profile}
	if now.Before(profile.CooldownUntil) {
		result.Skipped = true
		return result, nil
	}

	if now.Sub(profile.SpamWindowStart) > 3*time.Second {
		profile.SpamWindowStart = now
		profile.SpamCount = 1
	} else {
		profile.SpamCount++
		if profile.SpamCount >= 4 {
			profile.CooldownUntil = now.Add(5 * time.Second)
			profile.SpamCount = 0
			_, err = LocalLevelsDM.Set(bson.M{"_id": profile.ID}, profile)
			result.Skipped = true
			return result, err
		}
	}

	profile.XP += 1 + (now.UnixNano() % 15)
	profile.TotalMessages++
	profile.LastMessageTime = now
	profile.UpdatedAt = now

	nextLevel := profile.Level + 1
	requiredXP := nextLevel * nextLevel * 100
	if profile.XP >= requiredXP {
		profile.Level = nextLevel
		result.LevelUp = true
	}

	_, err = LocalLevelsDM.Set(bson.M{"_id": profile.ID}, profile)
	return result, err
}
