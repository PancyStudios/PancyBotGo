package models

import "time"

// UserLevelProfile represents the level and XP of a user in a specific guild
type UserLevelProfile struct {
	ID              string    `bson:"_id" json:"id"` // Format: GuildID_UserID
	GuildID         string    `bson:"guild_id" json:"guild_id"`
	UserID          string    `bson:"user_id" json:"user_id"`
	XP              int64     `bson:"xp" json:"xp"`
	Level           int64     `bson:"level" json:"level"`
	TotalMessages   int64     `bson:"total_messages" json:"total_messages"`
	LastMessageTime time.Time `bson:"last_message_time" json:"last_message_time"` // For cooldowns
	SpamWindowStart time.Time `bson:"spam_window_start" json:"spam_window_start"`
	SpamCount       int       `bson:"spam_count" json:"spam_count"`
	CooldownUntil   time.Time `bson:"cooldown_until" json:"cooldown_until"`
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at" json:"updated_at"`
}
