package models

import "time"

// BlacklistType representa el tipo de entrada en la blacklist
type BlacklistType string

const (
	BlacklistTypeUser  BlacklistType = "user"
	BlacklistTypeGuild BlacklistType = "guild"
)

// Blacklist representa una entrada en la blacklist
type Blacklist struct {
	ID      string        `bson:"_id" json:"id"`                            // User ID o Guild ID
	Type    BlacklistType `bson:"type" json:"type"`                         // "user" o "guild"
	Reason  string        `bson:"reason,omitempty" json:"reason,omitempty"` // Razón del blacklist
	AddedBy string        `bson:"added_by" json:"added_by"`                 // ID del admin que lo añadió
	AddedAt time.Time     `bson:"added_at" json:"added_at"`                 // Cuándo se añadió
}
