package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TempBan represents a temporary ban in the database
type TempBan struct {
	GuildID   string    `bson:"guildId"`
	UserID    string    `bson:"userId"`
	ExpiresAt time.Time `bson:"expiresAt"`
}

var client *discord.ExtendedClient

// StartTempBanScheduler starts a background goroutine to check for expired tempbans
func StartTempBanScheduler(c *discord.ExtendedClient) {
	client = c
	go func() {
		for {
			checkExpiredBans()
			time.Sleep(1 * time.Minute)
		}
	}()
}

func checkExpiredBans() {
	db := database.Get()
	if db == nil || !db.Connected() {
		return
	}

	col := db.GetCollection("tempbans")
	if col == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{"expiresAt": bson.M{"$lte": now}}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding expired tempbans: "+err.Error(), "Scheduler")
		return
	}
	defer cursor.Close(ctx)

	var expiredBans []TempBan
	if err := cursor.All(ctx, &expiredBans); err != nil {
		logger.Error("Error decoding expired tempbans: "+err.Error(), "Scheduler")
		return
	}

	for _, ban := range expiredBans {
		// Attempt to unban
		err := client.Session.GuildBanDelete(ban.GuildID, ban.UserID)
		if err != nil {
			logger.Warn(fmt.Sprintf("Could not unban user %s in guild %s: %v", ban.UserID, ban.GuildID, err), "Scheduler")
		} else {
			logger.Info(fmt.Sprintf("Tempban expired for user %s in guild %s", ban.UserID, ban.GuildID), "Scheduler")
		}

		// Delete from DB regardless of success (maybe user already unbanned manually)
		_, _ = col.DeleteOne(ctx, bson.M{"guildId": ban.GuildID, "userId": ban.UserID})
	}
}

// AddTempBan adds a new temporary ban to the database
func AddTempBan(guildID, userID string, duration time.Duration) error {
	db := database.Get()
	if db == nil || !db.Connected() {
		return fmt.Errorf("base de datos no conectada")
	}

	col := db.GetCollection("tempbans")
	if col == nil {
		return fmt.Errorf("no se pudo obtener la colección tempbans")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ban := TempBan{
		GuildID:   guildID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
	}

	opts := options.Update().SetUpsert(true)
	_, err := col.UpdateOne(ctx, bson.M{"guildId": guildID, "userId": userID}, bson.M{"$set": ban}, opts)
	return err
}
