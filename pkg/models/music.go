package models

// MusicSettings represents the music configuration for a guild
type MusicSettings struct {
	GuildID       string  `bson:"_id" json:"guildId"`   // Use _id as GuildID for easy fetching
	DjRole        *string `bson:"djRole" json:"djRole"` // Can be null
	DefaultVolume int     `bson:"defaultVolume" json:"defaultVolume"`
	StayInVc      bool    `bson:"stayInVc" json:"stayInVc"`
	ChannelID     *string `bson:"channelId" json:"channelId"` // Can be null
}

// NewMusicSettings creates a default MusicSettings instance
func NewMusicSettings(guildID string) *MusicSettings {
	return &MusicSettings{
		GuildID:       guildID,
		DefaultVolume: 100,
		StayInVc:      false,
	}
}
