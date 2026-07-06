package help

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createCmdsCommand() *discord.Command {
	return discord.NewCommand(
		"cmds",
		"📚 | Muestra la lista de todos los comandos disponibles",
		"help",
		cmdsHandler,
	)
}

func cmdsHandler(ctx *discord.CommandContext) error {
	commands := ctx.Client.CommandHandler.GetRegisteredCommands()

	if len(commands) == 0 {
		return ctx.ReplyEphemeral("❌ No hay comandos registrados actualmente.")
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📚 Lista de Comandos de PancyBot",
		Description: "Aquí tienes todos los comandos disponibles, categorizados automáticamente.",
		Color:       0x3498DB, // Blue
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ctx.Client.Session.State.User.AvatarURL("128"),
		},
		Fields: make([]*discordgo.MessageEmbedField, 0),
	}

	for _, cmd := range commands {
		// Construir una lista de subcomandos para este comando base
		var subcommands []string

		for _, opt := range cmd.Options {
			if opt.Type == discordgo.ApplicationCommandOptionSubCommand {
				subcommands = append(subcommands, fmt.Sprintf("`/%s %s` - %s", cmd.Name, opt.Name, cleanDesc(opt.Description)))
			} else if opt.Type == discordgo.ApplicationCommandOptionSubCommandGroup {
				for _, subOpt := range opt.Options {
					if subOpt.Type == discordgo.ApplicationCommandOptionSubCommand {
						subcommands = append(subcommands, fmt.Sprintf("`/%s %s %s` - %s", cmd.Name, opt.Name, subOpt.Name, cleanDesc(subOpt.Description)))
					}
				}
			}
		}

		if len(subcommands) > 0 {
			// Es un comando grupo (ej. /utils, /mod)
			fieldValue := strings.Join(subcommands, "\n")
			// Truncar si es muy largo (límite de Discord: 1024 caracteres por field value)
			if len(fieldValue) > 1024 {
				fieldValue = fieldValue[:1021] + "..."
			}

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("💠 Comando Base: `/%s`", cmd.Name),
				Value: fieldValue,
			})
		} else {
			// Es un comando simple raíz
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("💠 `/%s`", cmd.Name),
				Value: cleanDesc(cmd.Description),
			})
		}
	}

	// Si hay muchos fields, Discord permite máximo 25 fields por embed
	if len(embed.Fields) > 25 {
		embed.Fields = embed.Fields[:25]
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: "Se muestran hasta 25 grupos. Usa comandos específicos para más información.",
		}
	}

	return ctx.ReplyEmbed(embed)
}

// cleanDesc elimina los prefijos emoji si los hay para que la lista se vea más limpia
func cleanDesc(desc string) string {
	parts := strings.SplitN(desc, " | ", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return desc
}
