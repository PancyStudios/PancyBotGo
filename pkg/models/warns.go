package models

// Warn representa una advertencia individual
type Warn struct {
	Reason    string `bson:"reason" json:"reason"`
	Moderator string `bson:"moderator" json:"moderator"`
	ID        string `bson:"id" json:"id"`
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
}

// WarnsDocument representa el documento completo en la colección "Warns"
// Coincide con tu esquema de Mongoose: guildId, userId, warns[]
type WarnsDocument struct {
	GuildID string `bson:"guildId" json:"guildId"`
	UserID  string `bson:"userId" json:"userId"`
	Warns   []Warn `bson:"warns" json:"warns"`
}
