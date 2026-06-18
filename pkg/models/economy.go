package models

import (
	"time"
)

// ItemType defines the category of an item
type ItemType string

const (
	ItemTypeConsumable  ItemType = "consumable"
	ItemTypeCollectible ItemType = "collectible"
	ItemTypeTool        ItemType = "tool"
	ItemTypeRole        ItemType = "role"
	ItemTypeBadge       ItemType = "badge"
)

// Item represents a buyable/usable item in the economy
type Item struct {
	ID          string   `bson:"_id" json:"id"`
	GuildID     string   `bson:"guild_id" json:"guild_id"` // Empty if global
	Name        string   `bson:"name" json:"name"`
	Description string   `bson:"description" json:"description"`
	Price       int64    `bson:"price" json:"price"`
	SellPrice   int64    `bson:"sell_price" json:"sell_price"`
	Type        ItemType `bson:"type" json:"type"`
	Emoji       string   `bson:"emoji" json:"emoji"`
	Stock       int      `bson:"stock" json:"stock"` // -1 for infinite
	RoleID      string   `bson:"role_id,omitempty" json:"role_id,omitempty"` // If type is role
}

// GlobalEconomyProfile represents a user's global economy (Stars)
type GlobalEconomyProfile struct {
	UserID       string               `bson:"_id" json:"user_id"`
	StarsWallet  int64                `bson:"stars_wallet" json:"stars_wallet"`
	StarsBank    int64                `bson:"stars_bank" json:"stars_bank"`
	BankCapacity int64                `bson:"bank_capacity" json:"bank_capacity"`
	Inventory    map[string]int       `bson:"inventory" json:"inventory"` // ItemID -> Quantity
	Cooldowns    map[string]time.Time `bson:"cooldowns" json:"cooldowns"`
	CreatedAt    time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at" json:"updated_at"`
}

// LocalEconomyProfile represents a user's local economy in a specific server
type LocalEconomyProfile struct {
	ID           string               `bson:"_id" json:"id"` // Format: GuildID_UserID
	GuildID      string               `bson:"guild_id" json:"guild_id"`
	UserID       string               `bson:"user_id" json:"user_id"`
	Wallet       int64                `bson:"wallet" json:"wallet"`
	Bank         int64                `bson:"bank" json:"bank"`
	BankCapacity int64                `bson:"bank_capacity" json:"bank_capacity"`
	Inventory    map[string]int       `bson:"inventory" json:"inventory"` // ItemID -> Quantity
	Cooldowns    map[string]time.Time `bson:"cooldowns" json:"cooldowns"`
	CreatedAt    time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at" json:"updated_at"`
}
