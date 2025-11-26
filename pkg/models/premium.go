package models

import "time"

// UserPremium represents premium data for a user account
type UserPremium struct {
	UserID    string `bson:"user" json:"user"`
	Permanent bool   `bson:"permanent" json:"permanent"`
	ExpiresAt int64  `bson:"expira" json:"expira"`
}

// GuildPremium represents premium data for a guild (server)
type GuildPremium struct {
	GuildID   string `bson:"guild" json:"guild"`
	Permanent bool   `bson:"permanent" json:"permanent"`
	ExpiresAt int64  `bson:"expira" json:"expira"`
}

// PremiumCodeType represents the type of premium code
type PremiumCodeType string

const (
	PremiumCodeTypeUser  PremiumCodeType = "user"
	PremiumCodeTypeGuild PremiumCodeType = "guild"
)

// PremiumCode represents a premium code that can be redeemed
type PremiumCode struct {
	Code         string          `bson:"_id"`                  // El código es la llave primaria
	Type         PremiumCodeType `bson:"type"`                 // "user" o "guild"
	DurationDays int             `bson:"duration_days"`        // 30, 365, etc. (0 para permanente)
	Permanent    bool            `bson:"permanent"`            // Si es permanente
	IsClaimed    bool            `bson:"is_claimed"`           // Estado
	CreatedAt    time.Time       `bson:"created_at"`           // Cuándo se creó
	ClaimedBy    string          `bson:"claimed_by,omitempty"` // ID de usuario que lo reclamó
	ClaimedAt    time.Time       `bson:"claimed_at,omitempty"` // Cuándo se reclamó
	CreatedBy    string          `bson:"created_by,omitempty"` // ID del admin que lo creó
}
