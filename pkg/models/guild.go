package models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// GuildDocument represents the main document for a guild in MongoDB
type GuildDocument struct {
	ID            string             `bson:"_id" json:"id"`
	Configuration GuildConfiguration `bson:"configuration" json:"configuration"`
	Greetings     Greetings          `bson:"greetings" json:"greetings"`
	Moderation    ModeratorData      `bson:"moderation" json:"moderation"`
	Protection    ProtectionConfig   `bson:"protection" json:"protection"`
	Levels        LevelsConfig       `bson:"levels" json:"levels"`
	Embeds        []CustomEmbed      `bson:"embeds" json:"embeds"`
}

// CustomEmbed represents a user-created embed
type CustomEmbed struct {
	ID          string `bson:"id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Color       int    `bson:"color" json:"color"`
	Thumbnail   string `bson:"thumbnail" json:"thumbnail"`
	Image       string `bson:"image" json:"image"`
	FooterText  string `bson:"footerText" json:"footerText"`
	FooterIcon  string `bson:"footerIcon" json:"footerIcon"`
	AuthorName  string `bson:"authorName" json:"authorName"`
	AuthorIcon  string `bson:"authorIcon" json:"authorIcon"`
}

// LevelReward represents a role given at a specific level
type LevelReward struct {
	Level  int64  `bson:"level" json:"level"`
	RoleID string `bson:"roleId" json:"roleId"`
}

// LevelsConfig holds user level system settings
type LevelsConfig struct {
	Enable         bool          `bson:"enable" json:"enable"`
	LevelUpChannel string        `bson:"levelUpChannel" json:"levelUpChannel"` // Empty for same channel
	LevelUpMessage string        `bson:"levelUpMessage" json:"levelUpMessage"` // Available placeholders: {user}, {level}
	Rewards        []LevelReward `bson:"rewards" json:"rewards"`
}

// ProtectionConfig holds security settings like antibots and antiraid
type ProtectionConfig struct {
	Antibots             AntibotsConfig       `bson:"antibots" json:"antibots"`
	AntiRaid             AntiRaidConfig       `bson:"antiraid" json:"antiraid"` // TS uses antiraid lowercase
	AntiTokens           AntiTokensConfig     `bson:"antitokens" json:"antitokens"`
	AntiJoins            AntiJoinsConfig      `bson:"antijoins" json:"antijoins"`
	MarkMalicious        MarkMaliciousConfig  `bson:"markMalicious" json:"markMalicious"`
	WarnEntry            bool                 `bson:"warnEntry" json:"warnEntry"`
	KickMalicious        KickMaliciousConfig  `bson:"kickMalicious" json:"kickMalicious"`
	OwnSystem            OwnSystemConfig      `bson:"ownSystem" json:"ownSystem"`
	Verification         VerificationConfig   `bson:"verification" json:"verification"`
	CannotEnterTwice     CannotEnterTwiceConf `bson:"cannotEnterTwice" json:"cannotEnterTwice"`
	PurgeWebhooksAttacks PurgeWebhooksConfig  `bson:"purgeWebhooksAttacks" json:"purgeWebhooksAttacks"`
	IntelligentSOS       IntelligentSOSConfig `bson:"intelligentSOS" json:"intelligentSOS"`
	IntelligentAntiflood bool                 `bson:"intelligentAntiflood" json:"intelligentAntiflood"`
	Antiflood            bool                 `bson:"antiflood" json:"antiflood"`
	BloqEntritiesByName  BloqEntritiesConfig  `bson:"bloqEntritiesByName" json:"bloqEntritiesByName"`
	BloqNewCreatedUsers  BloqNewCreatedConfig `bson:"bloqNewCreatedUsers" json:"bloqNewCreatedUsers"`
	Raidmode             RaidmodeConfig       `bson:"raidmode" json:"raidmode"`
}

type AntibotsConfig struct {
	Enable bool   `bson:"enable" json:"enable"`
	Type   string `bson:"_type" json:"_type"`
}

// UnmarshalBSONValue handles decoding when antibots is a string in legacy data
func (a *AntibotsConfig) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t == bsontype.String {
		var s string
		if err := bson.UnmarshalValue(t, data, &s); err != nil {
			return err
		}
		a.Enable = (s == "enable" || s == "true")
		return nil
	}

	if t == bsontype.EmbeddedDocument {
		type Alias AntibotsConfig
		var alias Alias
		if err := bson.UnmarshalValue(t, data, &alias); err != nil {
			return err
		}
		*a = AntibotsConfig(alias)
		return nil
	}

	return fmt.Errorf("cannot decode %v into AntibotsConfig", t)
}

type AntiTokensConfig struct {
	Enable         bool     `bson:"enable" json:"enable"`
	UsersEntrities []string `bson:"usersEntrities" json:"usersEntrities"`
	EntritiesCount int      `bson:"entritiesCount" json:"entritiesCount"`
}

type AntiJoinsConfig struct {
	Enable            bool     `bson:"enable" json:"enable"`
	RememberEntrities []string `bson:"rememberEntrities" json:"rememberEntrities"`
}

type MarkMaliciousConfig struct {
	Enable            bool     `bson:"enable" json:"enable"`
	Type              string   `bson:"_type" json:"_type"`
	RememberEntrities []string `bson:"rememberEntrities" json:"rememberEntrities"`
}

type KickMaliciousConfig struct {
	Enable            bool     `bson:"enable" json:"enable"`
	RememberEntrities []string `bson:"rememberEntrities" json:"rememberEntrities"`
}

type OwnSystemConfig struct {
	Enable bool `bson:"enable" json:"enable"`
	Events struct {
		MessageCreate     []string `bson:"messageCreate" json:"messageCreate"`
		MessageDelete     []string `bson:"messageDelete" json:"messageDelete"`
		MessageUpdate     []string `bson:"messageUpdate" json:"messageUpdate"`
		ChannelCreate     []string `bson:"channelCreate" json:"channelCreate"`
		ChannelDelete     []string `bson:"channelDelete" json:"channelDelete"`
		ChannelUpdate     []string `bson:"channelUpdate" json:"channelUpdate"`
		RoleCreate        []string `bson:"roleCreate" json:"roleCreate"`
		RoleDelete        []string `bson:"roleDelete" json:"roleDelete"`
		RoleUpdate        []string `bson:"roleUpdate" json:"roleUpdate"`
		EmojiCreate       []string `bson:"emojiCreate" json:"emojiCreate"`
		EmojiDelete       []string `bson:"emojiDelete" json:"emojiDelete"`
		EmojiUpdate       []string `bson:"emojiUpdate" json:"emojiUpdate"`
		StickerCreate     []string `bson:"stickerCreate" json:"stickerCreate"`
		StickerDelete     []string `bson:"stickerDelete" json:"stickerDelete"`
		StickerUpdate     []string `bson:"stickerUpdate" json:"stickerUpdate"`
		GuildMemberAdd    []string `bson:"guildMemberAdd" json:"guildMemberAdd"`
		GuildMemberRemove []string `bson:"guildMemberRemove" json:"guildMemberRemove"`
		GuildMemberUpdate []string `bson:"guildMemberUpdate" json:"guildMemberUpdate"`
		GuildBanAdd       []string `bson:"guildBanAdd" json:"guildBanAdd"`
		GuildBanRemove    []string `bson:"guildBanRemove" json:"guildBanRemove"`
		InviteCreate      []string `bson:"inviteCreate" json:"inviteCreate"`
		InviteDelete      []string `bson:"inviteDelete" json:"inviteDelete"`
		ThreadCreate      []string `bson:"threadCreate" json:"threadCreate"`
		ThreadDelete      []string `bson:"threadDelete" json:"threadDelete"`
	} `bson:"events" json:"events"`
}

type VerificationConfig struct {
	Enable            bool   `bson:"enable" json:"enable"`
	Type              string `bson:"_type" json:"_type"`
	Channel           string `bson:"channel" json:"channel"`
	Role              string `bson:"role" json:"role"`
	MinAccountAgeDays int    `bson:"minAccountAgeDays" json:"minAccountAgeDays"`
}

type CannotEnterTwiceConf struct {
	Enable bool     `bson:"enable" json:"enable"`
	Users  []string `bson:"users" json:"users"`
}

type PurgeWebhooksConfig struct {
	Enable         bool   `bson:"enable" json:"enable"`
	Amount         int    `bson:"amount" json:"amount"`
	RememberOwners string `bson:"rememberOwners" json:"rememberOwners"`
}

type IntelligentSOSConfig struct {
	Enable   bool `bson:"enable" json:"enable"`
	Cooldown bool `bson:"cooldown" json:"cooldown"`
}

type BloqEntritiesConfig struct {
	Enable bool     `bson:"enable" json:"enable"`
	Names  []string `bson:"names" json:"names"`
}

type BloqNewCreatedConfig struct {
	Time string `bson:"time" json:"time"`
}

type RaidmodeConfig struct {
	Enable        bool   `bson:"enable" json:"enable"`
	TimeToDisable string `bson:"timeToDisable" json:"timeToDisable"`
	Password      string `bson:"password" json:"password"`
	ActivedDate   int    `bson:"activedDate" json:"activedDate"`
}

// AntiRaidConfig holds the configuration for raid protection
type AntiRaidConfig struct {
	Enable            bool `bson:"enable" json:"enable"`
	Amount            int  `bson:"amount" json:"amount"`
	SaveBotsEntrities struct {
		AuthorOfEntry string `bson:"authorOfEntry" json:"authorOfEntry"`
		Bot           string `bson:"_bot" json:"_bot"`
	} `bson:"saveBotsEntrities" json:"saveBotsEntrities"`
	// Additional fields added by the Go bot:
	Action            string `bson:"action,omitempty" json:"action,omitempty"`
	MinAccountAgeDays int    `bson:"minAccountAgeDays,omitempty" json:"minAccountAgeDays,omitempty"`
	JoinLimit         int    `bson:"joinLimit,omitempty" json:"joinLimit,omitempty"`
	TimeWindow        int    `bson:"timeWindow,omitempty" json:"timeWindow,omitempty"`
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
	EmbedID string `bson:"embedId" json:"embedId"`
	IsDM    bool   `bson:"isDM" json:"isDM"`
}

// FarewellConfig holds farewell message settings
type FarewellConfig struct {
	Enable  bool   `bson:"enable" json:"enable"`
	Channel string `bson:"channel" json:"channel"`
	Message string `bson:"message" json:"message"`
	EmbedID string `bson:"embedId" json:"embedId"`
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

// NewDefaultGuildDocument creates a new GuildDocument with all default values initialized
func NewDefaultGuildDocument(guildID string) *GuildDocument {
	return &GuildDocument{
		ID: guildID,
		Configuration: GuildConfiguration{
			Version:        "1.0.0",
			Prefix:         "pan!",
			Language:       "es",
			Whitelist:      []string{},
			Logs:           []string{},
			LogsChannel:    "",
			IgnoreChannels: []string{},
			Password: PasswordConfig{
				Enable:          false,
				Password:        "",
				UsersWithAccess: []string{},
			},
			SubData: SubDataConfig{
				ShowDetailsInCmdsCommand:         "lessDetails",
				PingMessage:                      "allDetails",
				DontRepeatTheAutomoderatorAction: false,
			},
		},
		Greetings: Greetings{
			Welcome: WelcomeConfig{
				Enable:  false,
				Channel: "",
				Message: "",
				IsDM:    false,
			},
			Farewell: FarewellConfig{
				Enable:  false,
				Channel: "",
				Message: "",
			},
			Autorole: AutoroleConfig{
				Enable: false,
				Roles:  []string{},
				Delay:  0,
			},
		},
		Moderation: ModeratorData{
			Logs: ModLogsConfig{
				Warns: LogChannelConfig{Enable: false, Channel: ""},
				Mutes: LogChannelConfig{Enable: false, Channel: ""},
				Kicks: LogChannelConfig{Enable: false, Channel: ""},
				Bans:  LogChannelConfig{Enable: false, Channel: ""},
			},
			DataModeration: DataModerationConfig{
				MuteRole:     "",
				ForceReasons: []string{},
				Timers:       []interface{}{},
				BadWords:     []string{},
				Events: ModEventsConfig{
					ManyPings:      false,
					CapitalLetters: false,
					ManyEmojis:     false,
					ManyWords:      false,
					LinkDetect:     false,
					Ghostping:      false,
					NsfwFilter:     false,
					IpLoggerFilter: false,
				},
			},
			Automoderator: AutomoderatorConfig{
				Enable: false,
				Actions: AutomoderatorActions{
					Warns:       []int{},
					MuteTime:    []int{},
					Action:      "",
					FloodDetect: 0,
					ManyEmojis:  0,
					ManyPings:   0,
					ManyWords:   0,
				},
				Events: ModEventsConfig{
					ManyPings:      false,
					CapitalLetters: false,
					ManyEmojis:     false,
					ManyWords:      false,
					LinkDetect:     false,
					Ghostping:      false,
					NsfwFilter:     false,
					IpLoggerFilter: false,
				},
			},
		},
		Protection: ProtectionConfig{
			AntiRaid: AntiRaidConfig{
				Enable: false,
				Amount: 0,
				SaveBotsEntrities: struct {
					AuthorOfEntry string `bson:"authorOfEntry" json:"authorOfEntry"`
					Bot           string `bson:"_bot" json:"_bot"`
				}{"", ""},
			},
			Antibots: AntibotsConfig{Enable: false, Type: "all"},
			AntiTokens: AntiTokensConfig{
				Enable:         false,
				UsersEntrities: []string{},
				EntritiesCount: 0,
			},
			AntiJoins: AntiJoinsConfig{Enable: false, RememberEntrities: []string{}},
			MarkMalicious: MarkMaliciousConfig{
				Enable:            true,
				Type:              "changeNickname",
				RememberEntrities: []string{},
			},
			WarnEntry:     true,
			KickMalicious: KickMaliciousConfig{Enable: false, RememberEntrities: []string{}},
			OwnSystem: OwnSystemConfig{
				Enable: false,
				Events: struct {
					MessageCreate     []string `bson:"messageCreate" json:"messageCreate"`
					MessageDelete     []string `bson:"messageDelete" json:"messageDelete"`
					MessageUpdate     []string `bson:"messageUpdate" json:"messageUpdate"`
					ChannelCreate     []string `bson:"channelCreate" json:"channelCreate"`
					ChannelDelete     []string `bson:"channelDelete" json:"channelDelete"`
					ChannelUpdate     []string `bson:"channelUpdate" json:"channelUpdate"`
					RoleCreate        []string `bson:"roleCreate" json:"roleCreate"`
					RoleDelete        []string `bson:"roleDelete" json:"roleDelete"`
					RoleUpdate        []string `bson:"roleUpdate" json:"roleUpdate"`
					EmojiCreate       []string `bson:"emojiCreate" json:"emojiCreate"`
					EmojiDelete       []string `bson:"emojiDelete" json:"emojiDelete"`
					EmojiUpdate       []string `bson:"emojiUpdate" json:"emojiUpdate"`
					StickerCreate     []string `bson:"stickerCreate" json:"stickerCreate"`
					StickerDelete     []string `bson:"stickerDelete" json:"stickerDelete"`
					StickerUpdate     []string `bson:"stickerUpdate" json:"stickerUpdate"`
					GuildMemberAdd    []string `bson:"guildMemberAdd" json:"guildMemberAdd"`
					GuildMemberRemove []string `bson:"guildMemberRemove" json:"guildMemberRemove"`
					GuildMemberUpdate []string `bson:"guildMemberUpdate" json:"guildMemberUpdate"`
					GuildBanAdd       []string `bson:"guildBanAdd" json:"guildBanAdd"`
					GuildBanRemove    []string `bson:"guildBanRemove" json:"guildBanRemove"`
					InviteCreate      []string `bson:"inviteCreate" json:"inviteCreate"`
					InviteDelete      []string `bson:"inviteDelete" json:"inviteDelete"`
					ThreadCreate      []string `bson:"threadCreate" json:"threadCreate"`
					ThreadDelete      []string `bson:"threadDelete" json:"threadDelete"`
				}{
					[]string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{},
				},
			},
			Verification:         VerificationConfig{Enable: false, Type: "button", Channel: "", Role: "", MinAccountAgeDays: 0},
			CannotEnterTwice:     CannotEnterTwiceConf{Enable: false, Users: []string{}},
			PurgeWebhooksAttacks: PurgeWebhooksConfig{Enable: false, Amount: 0, RememberOwners: "Nadie"},
			IntelligentSOS:       IntelligentSOSConfig{Enable: false, Cooldown: false},
			IntelligentAntiflood: false,
			Antiflood:            true,
			BloqEntritiesByName:  BloqEntritiesConfig{Enable: false, Names: []string{"raider", "doxer", "hacker", "infecter"}},
			BloqNewCreatedUsers:  BloqNewCreatedConfig{Time: "1h"},
			Raidmode:             RaidmodeConfig{Enable: false, TimeToDisable: "1d", Password: "Nothing", ActivedDate: 0},
		},
		Levels: LevelsConfig{
			Enable:         true,
			LevelUpChannel: "",
			LevelUpMessage: "",
		},
	}
}
