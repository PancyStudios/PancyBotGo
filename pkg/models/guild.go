package models

// GuildDocument represents the main document for a guild in MongoDB
type GuildDocument struct {
	ID            string             `bson:"_id" json:"id"`
	Configuration GuildConfiguration `bson:"configuration" json:"configuration"`
	Greetings     Greetings          `bson:"greetings" json:"greetings"`
	Moderation    ModeratorData      `bson:"moderation" json:"moderation"`
	Protection    ProtectionConfig   `bson:"protection" json:"protection"`
	Levels        LevelsConfig       `bson:"levels" json:"levels"`
}

// LevelsConfig holds user level system settings
type LevelsConfig struct {
	Enable           bool   `bson:"enable" json:"enable"`
	LevelUpChannel   string `bson:"levelUpChannel" json:"levelUpChannel"`     // Empty for same channel
	LevelUpMessage   string `bson:"levelUpMessage" json:"levelUpMessage"`     // Available placeholders: {user}, {level}
}

// ProtectionConfig holds security settings like antibots
type ProtectionConfig struct {
	Antibots string `bson:"antibots" json:"antibots"` // "all", "only_nv", "only_v", or "" (disabled)
}

// GuildConfiguration holds general configuration for the bot in a guild
type GuildConfiguration struct {
	Version        string         `bson:"_version" json:"_version"`
	Prefix         string         `bson:"prefix" json:"prefix"`
	Logs           []string       `bson:"logs" json:"logs"` // AuditLogEvent equivalent
	LogsChannel    string         `bson:"logsChannel" json:"logsChannel"`
	Language       string         `bson:"language" json:"language"`
	IgnoreChannels []string       `bson:"ignoreChannels" json:"ignoreChannels"`
	Password       PasswordConfig `bson:"password" json:"password"`
	Whitelist      []string       `bson:"whitelist" json:"whitelist"`
	SubData        SubDataConfig  `bson:"subData" json:"subData"`
}

// PasswordConfig holds password protection settings
type PasswordConfig struct {
	Enable          bool     `bson:"enable" json:"enable"`
	Password        string   `bson:"_password" json:"_password"`
	UsersWithAccess []string `bson:"usersWithAcces" json:"usersWithAcces"` // Typo from TS kept for compatibility
}

// SubDataConfig holds miscellaneous guild settings
type SubDataConfig struct {
	ShowDetailsInCmdsCommand         string `bson:"showDetailsInCmdsCommand" json:"showDetailsInCmdsCommand"`
	PingMessage                      string `bson:"pingMessage" json:"pingMessage"`
	DontRepeatTheAutomoderatorAction bool   `bson:"dontRepeatTheAutomoderatorAction" json:"dontRepeatTheAutomoderatorAction"`
	SuggestChannel                   string `bson:"suggestChannel" json:"suggestChannel"`
	ConfessionChannel                string `bson:"confessionChannel" json:"confessionChannel"`
	VerifyChannel                    string `bson:"verifyChannel" json:"verifyChannel"`
	VerifyRole                       string `bson:"verifyRole" json:"verifyRole"`
}

// Greetings holds welcome, farewell, and autorole configurations
type Greetings struct {
	Welcome  WelcomeConfig  `bson:"welcome" json:"welcome"`
	Farewell FarewellConfig `bson:"farewell" json:"farewell"`
	Autorole AutoroleConfig `bson:"autorole" json:"autorole"`
}

// WelcomeConfig holds welcome message settings
type WelcomeConfig struct {
	Enable  bool   `bson:"enable" json:"enable"`
	Channel string `bson:"channel" json:"channel"`
	Message string `bson:"message" json:"message"`
	IsDM    bool   `bson:"isDM" json:"isDM"`
}

// FarewellConfig holds farewell message settings
type FarewellConfig struct {
	Enable  bool   `bson:"enable" json:"enable"`
	Channel string `bson:"channel" json:"channel"`
	Message string `bson:"message" json:"message"`
}

// AutoroleConfig holds autorole settings
type AutoroleConfig struct {
	Enable bool     `bson:"enable" json:"enable"`
	Roles  []string `bson:"roles" json:"roles"`
	Delay  int      `bson:"delay" json:"delay"` // Delay in ms
}

// ModeratorData holds moderation and automoderator settings
type ModeratorData struct {
	Logs           ModLogsConfig        `bson:"logs" json:"logs"`
	DataModeration DataModerationConfig `bson:"dataModeration" json:"dataModeration"`
	Automoderator  AutomoderatorConfig  `bson:"automoderator" json:"automoderator"`
}

// ModLogsConfig holds logging channels for moderation actions
type ModLogsConfig struct {
	Warns LogChannelConfig `bson:"warns" json:"warns"`
	Mutes LogChannelConfig `bson:"mutes" json:"mutes"`
	Kicks LogChannelConfig `bson:"kicks" json:"kicks"`
	Bans  LogChannelConfig `bson:"bans" json:"bans"`
}

// LogChannelConfig holds enable status and channel ID for a log
type LogChannelConfig struct {
	Enable  bool   `bson:"enable" json:"enable"`
	Channel string `bson:"channel" json:"channel"`
}

// DataModerationConfig holds manual moderation rules
type DataModerationConfig struct {
	MuteRole     string          `bson:"muterole" json:"muterole"`
	ForceReasons []string        `bson:"forceReasons" json:"forceReasons"`
	Timers       []interface{}   `bson:"timers" json:"timers"`
	BadWords     []string        `bson:"badwords" json:"badwords"`
	Events       ModEventsConfig `bson:"events" json:"events"`
}

// ModEventsConfig holds boolean flags for filtering events
type ModEventsConfig struct {
	ManyPings      bool `bson:"manyPings" json:"manyPings"`
	CapitalLetters bool `bson:"capitalLetters" json:"capitalLetters"`
	ManyEmojis     bool `bson:"manyEmojis" json:"manyEmojis"`
	ManyWords      bool `bson:"manyWords" json:"manyWords"`
	LinkDetect     bool `bson:"linkDetect" json:"linkDetect"`
	Ghostping      bool `bson:"ghostping" json:"ghostping"`
	NsfwFilter     bool `bson:"nsfwFilter" json:"nsfwFilter"`
	IpLoggerFilter bool `bson:"iploggerFilter" json:"iploggerFilter"`
}

// AutomoderatorConfig holds automatic moderation action thresholds
type AutomoderatorConfig struct {
	Enable  bool                 `bson:"enable" json:"enable"`
	Actions AutomoderatorActions `bson:"actions" json:"actions"`
	Events  ModEventsConfig      `bson:"events" json:"events"`
}

// AutomoderatorActions holds specific thresholds for automoderation
type AutomoderatorActions struct {
	Warns       []int  `bson:"warns" json:"warns"`
	MuteTime    []int  `bson:"muteTime" json:"muteTime"`
	Action      string `bson:"action" json:"action"`
	FloodDetect int    `bson:"floodDetect" json:"floodDetect"`
	ManyEmojis  int    `bson:"manyEmojis" json:"manyEmojis"`
	ManyPings   int    `bson:"manyPings" json:"manyPings"`
	ManyWords   int    `bson:"manyWords" json:"manyWords"`
}
