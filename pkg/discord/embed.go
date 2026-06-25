package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// Brand colors
const (
	ColorSuccess = 0x2ECC71 // Green
	ColorError   = 0xE74C3C // Red
	ColorWarning = 0xF1C40F // Yellow/Gold
	ColorInfo    = 0x5865F2 // Blurple (Default Discord)
	ColorBot     = 0x2b2d31 // Dark (Discord Theme)
)

// EmbedBuilder provides methods to construct standard stylized embeds
type EmbedBuilder struct {
	embed *discordgo.MessageEmbed
}

// NewEmbed creates a new base embed builder with common styling
func NewEmbed() *EmbedBuilder {
	return &EmbedBuilder{
		embed: &discordgo.MessageEmbed{
			Color:     ColorInfo,
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "PancyBot Go",
			},
		},
	}
}

// SetTitle sets the embed title
func (b *EmbedBuilder) SetTitle(title string) *EmbedBuilder {
	b.embed.Title = title
	return b
}

// SetDescription sets the embed description
func (b *EmbedBuilder) SetDescription(desc string) *EmbedBuilder {
	b.embed.Description = desc
	return b
}

// SetColor sets the embed color
func (b *EmbedBuilder) SetColor(color int) *EmbedBuilder {
	b.embed.Color = color
	return b
}

// AddField adds a field to the embed
func (b *EmbedBuilder) AddField(name, value string, inline bool) *EmbedBuilder {
	b.embed.Fields = append(b.embed.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
	return b
}

// SetThumbnail sets the embed thumbnail
func (b *EmbedBuilder) SetThumbnail(url string) *EmbedBuilder {
	b.embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: url}
	return b
}

// SetImage sets the embed image
func (b *EmbedBuilder) SetImage(url string) *EmbedBuilder {
	b.embed.Image = &discordgo.MessageEmbedImage{URL: url}
	return b
}

// SetAuthor sets the embed author
func (b *EmbedBuilder) SetAuthor(name, iconURL string) *EmbedBuilder {
	b.embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    name,
		IconURL: iconURL,
	}
	return b
}

// SetFooter sets the embed footer
func (b *EmbedBuilder) SetFooter(text, iconURL string) *EmbedBuilder {
	b.embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    text,
		IconURL: iconURL,
	}
	return b
}

// Build returns the constructed MessageEmbed
func (b *EmbedBuilder) Build() *discordgo.MessageEmbed {
	return b.embed
}

// --- Predefined Embed Helpers ---

// NewSuccessEmbed creates a standard success embed
func NewSuccessEmbed(title, description string) *discordgo.MessageEmbed {
	return NewEmbed().
		SetTitle(title).
		SetDescription(description).
		SetColor(ColorSuccess).
		Build()
}

// NewErrorEmbed creates a standard error embed
func NewErrorEmbed(title, description string) *discordgo.MessageEmbed {
	return NewEmbed().
		SetTitle(title).
		SetDescription(description).
		SetColor(ColorError).
		Build()
}

// NewWarningEmbed creates a standard warning embed
func NewWarningEmbed(title, description string) *discordgo.MessageEmbed {
	return NewEmbed().
		SetTitle(title).
		SetDescription(description).
		SetColor(ColorWarning).
		Build()
}

// NewInfoEmbed creates a standard info embed
func NewInfoEmbed(title, description string) *discordgo.MessageEmbed {
	return NewEmbed().
		SetTitle(title).
		SetDescription(description).
		SetColor(ColorInfo).
		Build()
}

// SimpleEmbed parses a simple text message and returns an appropriate embed based on prefixes
func SimpleEmbed(content string) *discordgo.MessageEmbed {
	// If it starts with an error emoji, make it an Error embed
	if len(content) > 3 && (content[:3] == "❌" || content[:4] == "❌ ") {
		return NewErrorEmbed("Error", content)
	}
	// If it starts with success emoji
	if len(content) > 3 && (content[:3] == "✅" || content[:4] == "✅ ") {
		return NewSuccessEmbed("Éxito", content)
	}
	// Default to Info embed
	return NewInfoEmbed("Información", content)
}
