package messagecommands

import (
	"strings"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

// MessageContext holds information about a message command execution
type MessageContext struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
	Args    []string
}

// Reply sends a simple text response
func (ctx *MessageContext) Reply(content string) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, content)
}

// ReplyEmbed sends an embed response
func (ctx *MessageContext) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
}

// ReplyError sends an error embed response
func (ctx *MessageContext) ReplyError(title, content string) (*discordgo.Message, error) {
	return ctx.ReplyEmbed(discord.NewErrorEmbed(title, content))
}

// ReplySuccess sends a success embed response
func (ctx *MessageContext) ReplySuccess(title, content string) (*discordgo.Message, error) {
	return ctx.ReplyEmbed(discord.NewSuccessEmbed(title, content))
}

// CleanMention removes <@>, <@!>, <@&>, <#> from a mention string and returns just the ID
func CleanMention(mention string) string {
	mention = strings.TrimPrefix(mention, "<@!")
	mention = strings.TrimPrefix(mention, "<@&")
	mention = strings.TrimPrefix(mention, "<@")
	mention = strings.TrimPrefix(mention, "<#")
	mention = strings.TrimSuffix(mention, ">")
	return mention
}

// ParseUser tries to extract a User ID from arguments at the given index
func (ctx *MessageContext) ParseUser(index int) string {
	if index >= len(ctx.Args) {
		return ""
	}
	return CleanMention(ctx.Args[index])
}

// ParseRole gets a role ID from the argument at the given index (handles mentions or IDs)
func (ctx *MessageContext) ParseRole(index int) string {
	if index < 0 || index >= len(ctx.Args) {
		return ""
	}
	arg := ctx.Args[index]
	roleID := strings.TrimPrefix(arg, "<@&")
	roleID = strings.TrimSuffix(roleID, ">")
	return roleID
}

// ParseChannel gets a channel ID from the argument at the given index (handles mentions or IDs)
func (ctx *MessageContext) ParseChannel(index int) string {
	if index < 0 || index >= len(ctx.Args) {
		return ""
	}
	arg := ctx.Args[index]
	channelID := strings.TrimPrefix(arg, "<#")
	channelID = strings.TrimSuffix(channelID, ">")
	return channelID
}

// HasPermission checks if the user has a specific permission in the channel
func (ctx *MessageContext) HasPermission(permission int64) bool {
	perms, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.Message.ChannelID)
	if err != nil {
		return false
	}
	return perms&permission == permission
}

// CommandRunFunc represents a message command handler function
type CommandRunFunc func(ctx *MessageContext) error

var registry = make(map[string]CommandRunFunc)

// RegisterCommand adds a command to the prefix router
func RegisterCommand(name string, handler CommandRunFunc) {
	registry[strings.ToLower(name)] = handler
}

// GetRegisteredCommands returns a list of all registered command names
func GetRegisteredCommands() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// Handle routes an incoming message to the correct command handler
func Handle(s *discordgo.Session, m *discordgo.MessageCreate, commandName string, args []string) error {
	handler, exists := registry[strings.ToLower(commandName)]
	if !exists {
		// Command not found
		return nil
	}

	ctx := &MessageContext{
		Session: s,
		Message: m,
		Args:    args,
	}

	return handler(ctx)
}
