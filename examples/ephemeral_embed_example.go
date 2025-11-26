// Package examples demonstrates how to use ReplyEphemeralEmbed
// This is useful for sending embeds that are only visible to the user who triggered the command
package examples

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

// ExampleEphemeralEmbedCommand shows how to use ReplyEphemeralEmbed
func ExampleEphemeralEmbedCommand() *discord.Command {
	return discord.NewCommand(
		"example",
		"Ejemplo de respuesta embed efímera",
		"util",
		exampleHandler,
	)
}

// exampleHandler demonstrates sending an ephemeral embed reply
func exampleHandler(ctx *discord.CommandContext) error {
	// Create an embed
	embed := &discordgo.MessageEmbed{
		Title:       "Respuesta Privada",
		Description: "Este mensaje solo es visible para ti 🔒",
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Información",
				Value:  "Los mensajes efímeros son útiles para comandos sensibles",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Solo tú puedes ver esto",
		},
	}

	// Send as ephemeral embed (only visible to the user)
	return ctx.ReplyEphemeralEmbed(embed)
}

// Example of using ReplyEphemeralEmbed in a moderation command
// Note: This example uses Defer + EditReplyEmbed for async operations.
// For immediate ephemeral embeds, use ctx.ReplyEphemeralEmbed directly.
func exampleWarningsHandler(ctx *discord.CommandContext) error {
	// For async operations with embeds, use Defer first (ephemeral by default when first reply was ephemeral)
	// Then do processing in a goroutine and use EditReplyEmbed
	
	// Option 1: Direct ephemeral embed (for quick responses)
	embed := &discordgo.MessageEmbed{
		Title:       "✅ Advertencias",
		Description: "Aquí están las advertencias del usuario",
		Color:       0x00FF00,
	}
	return ctx.ReplyEphemeralEmbed(embed)
	
	// Option 2: For async operations, use Defer + goroutine + EditReplyEmbed
	// Note: Defer doesn't support ephemeral flag directly in discordgo
	// So for ephemeral async responses, send a quick ephemeral reply first
}
