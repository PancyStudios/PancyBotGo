package models

import "time"

// BlacklistType represents the type of blacklist entry
type BlacklistType string

const (
	BlacklistTypeUser  BlacklistType = "user"
	BlacklistTypeGuild BlacklistType = "guild"
)

// BlacklistEntry represents a blacklisted user or guild
type BlacklistEntry struct {
	ID        string        `bson:"_id"`        // User ID or Guild ID
	Type      BlacklistType `bson:"type"`       // "user" o "guild"
	Reason    string        `bson:"reason"`     // Razón del bloqueo
	CreatedAt time.Time     `bson:"created_at"` // Cuándo se creó
	CreatedBy string        `bson:"created_by"` // ID del desarrollador que lo bloqueó
}
