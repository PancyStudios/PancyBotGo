package database

import (
	"context"
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find().
		SetSort(bson.D{{Key: "xp", Value: -1}}).
		SetLimit(limit)

	cursor, err := LocalLevelsDM.collection.Find(ctx, query, findOptions)
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
