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

// exampleWarningsHandlerDirect shows direct ephemeral embed response (for quick/sync operations)
func exampleWarningsHandlerDirect(ctx *discord.CommandContext) error {
	embed := &discordgo.MessageEmbed{
		Title:       "✅ Advertencias",
		Description: "Aquí están las advertencias del usuario",
		Color:       0x00FF00,
	}
	return ctx.ReplyEphemeralEmbed(embed)
}

// exampleWarningsHandlerAsync shows async processing with ephemeral embeds
// Note: For async operations, send ephemeral reply first, then edit it after processing
func exampleWarningsHandlerAsync(ctx *discord.CommandContext) error {
	// Send initial ephemeral loading message
	loadingEmbed := &discordgo.MessageEmbed{
		Title:       "⏳ Cargando...",
		Description: "Obteniendo información de advertencias...",
		Color:       0xFFFF00,
	}
	
	if err := ctx.ReplyEphemeralEmbed(loadingEmbed); err != nil {
		return err
	}
	
	// Process in background
	go func() {
		// ... do some async processing ...
		
		// Update with final result
		// Note: EditReplyEmbed maintains the ephemeral flag from the initial reply
		resultEmbed := &discordgo.MessageEmbed{
			Title:       "✅ Advertencias",
			Description: "Aquí están las advertencias del usuario",
			Color:       0x00FF00,
		}
		ctx.EditReplyEmbed(resultEmbed)
	}()
	
	return nil
}
