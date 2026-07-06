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

	menu := createSlashMenu(ctx.Client)

	return ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: menu,
		},
	})
}

func HandleSlashInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, client *discord.ExtendedClient) bool {
	if i.Type != discordgo.InteractionMessageComponent {
		return false
	}
	
	data := i.MessageComponentData()
	if data.CustomID != "slash_help_cmds_menu" || len(data.Values) == 0 {
		return false
	}

	val := data.Values[0]
	if val == "none" {
		embed := &discordgo.MessageEmbed{
			Title:       "📚 Lista de Comandos de PancyBot",
			Description: "Aquí tienes todos los comandos disponibles, categorizados automáticamente. ¡Selecciona una en el menú de abajo!",
			Color:       0x3498DB,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: s.State.User.AvatarURL("128"),
			},
		}

		menu := createSlashMenu(client)

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Components: menu,
			},
		})
		if err != nil {
			fmt.Println("Error in InteractionRespond:", err)
		}
		return true
	}

	if !strings.HasPrefix(val, "slash_help_cat_") {
		return false
	}

	categoryName := strings.TrimPrefix(val, "slash_help_cat_")

	commands := client.CommandHandler.GetRegisteredCommands()
	var subcommands []string
	var cmdDesc string

	for _, cmd := range commands {
		if cmd.Name == categoryName {
			cmdDesc = cleanDesc(cmd.Description)
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
			break
		}
	}

	if len(subcommands) == 0 {
		subcommands = append(subcommands, "Este comando no tiene subcomandos.")
	}

	valStr := strings.Join(subcommands, "\n")
	if len(valStr) > 4000 {
		valStr = valStr[:3997] + "..."
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📚 Comando: /" + categoryName,
		Description: "**Descripción:** " + cmdDesc + "\n\n**Subcomandos:**\n" + valStr,
		Color:       0x3498DB,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: s.State.User.AvatarURL("128"),
		},
	}

	menu := createSlashMenu(client)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: menu,
		},
	})
	if err != nil {
		fmt.Println("Error in InteractionRespond:", err)
	}
	return true
}

// cleanDesc elimina los prefijos emoji si los hay para que la lista se vea más limpia
func cleanDesc(desc string) string {
	parts := strings.SplitN(desc, " | ", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return desc
}

func createSlashMenu(client *discord.ExtendedClient) []discordgo.MessageComponent {
	commands := client.CommandHandler.GetRegisteredCommands()

	options := []discordgo.SelectMenuOption{
		{
			Label:       "Selecciona una categoría",
			Value:       "none",
			Description: "Muestra las categorías disponibles",
			Emoji: &discordgo.ComponentEmoji{
				Name: "📚",
			},
			Default: true,
		},
	}

	for _, cmd := range commands {
		options = append(options, discordgo.SelectMenuOption{
			Label:       "/" + cmd.Name,
			Value:       "slash_help_cat_" + cmd.Name,
			Description: cleanDesc(cmd.Description),
			Emoji: &discordgo.ComponentEmoji{
				Name: "💠",
			},
		})
	}

	// Discord allows max 25 options per select menu
	if len(options) > 25 {
		options = options[:25]
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "slash_help_cmds_menu",
					Placeholder: "Selecciona una categoría...",
					Options:     options,
				},
			},
		},
	}
}
