package help

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func cmdsCommand(ctx *messagecommands.MessageContext) error {
	commands := messagecommands.GetRegisteredCommands()

	if len(commands) == 0 {
		_, err := ctx.ReplyError("Error", "❌ No hay comandos registrados actualmente.")
		return err
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📚 Lista de Comandos de PancyBot",
		Description: "Aquí tienes todos los comandos de prefijo disponibles agrupados por categoría.",
		Color:       0x3498DB,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ctx.Session.State.User.AvatarURL("128"),
		},
		Fields: []*discordgo.MessageEmbedField{},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Usa el prefijo de tu servidor o menciona al bot antes de cada comando.",
		},
	}

	categories := make(map[string][]string)
	for _, cmd := range commands {
		cat := cmd.Category
		if cat == "" {
			cat = "General"
		}
		categories[cat] = append(categories[cat], fmt.Sprintf("`%s` - %s", cmd.Usage, cmd.Description))
	}

	menu := createPrefixMenu()

	msgSend := &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: menu,
	}

	_, err := ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, msgSend)
	return err
}

func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if i.Type != discordgo.InteractionMessageComponent {
		return false
	}
	
	data := i.MessageComponentData()
	if data.CustomID != "help_cmds_menu" || len(data.Values) == 0 {
		return false
	}

	val := data.Values[0]
	if val == "none" {
		embed := &discordgo.MessageEmbed{
			Title:       "📚 Lista de Comandos de PancyBot",
			Description: "Aquí tienes todos los comandos de prefijo disponibles agrupados por categoría. ¡Selecciona una en el menú de abajo!",
			Color:       0x3498DB,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: s.State.User.AvatarURL("128"),
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Usa el prefijo de tu servidor o menciona al bot antes de cada comando.",
			},
		}

		menu := createPrefixMenu()

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

	if !strings.HasPrefix(val, "help_cat_") {
		return false
	}

	categoryName := strings.TrimPrefix(val, "help_cat_")

	commands := messagecommands.GetRegisteredCommands()
	var categoryCommands []string

	for _, cmd := range commands {
		cat := cmd.Category
		if cat == "" {
			cat = "General"
		}
		if cat == categoryName {
			categoryCommands = append(categoryCommands, fmt.Sprintf("`%s` - %s", cmd.Usage, cmd.Description))
		}
	}

	valStr := strings.Join(categoryCommands, "\n")
	if len(valStr) > 4000 {
		valStr = valStr[:3997] + "..."
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📚 Comandos: " + categoryName,
		Description: valStr,
		Color:       0x3498DB,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: s.State.User.AvatarURL("128"),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Usa el prefijo de tu servidor o menciona al bot antes de cada comando.",
		},
	}

	menu := createPrefixMenu()

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

func createPrefixMenu() []discordgo.MessageComponent {
	commands := messagecommands.GetRegisteredCommands()
	categories := make(map[string][]string)
	for _, cmd := range commands {
		cat := cmd.Category
		if cat == "" {
			cat = "General"
		}
		categories[cat] = append(categories[cat], cmd.Name)
	}

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

	for cat := range categories {
		options = append(options, discordgo.SelectMenuOption{
			Label:       cat,
			Value:       "help_cat_" + cat,
			Description: "Ver comandos de " + cat,
			Emoji: &discordgo.ComponentEmoji{
				Name: "💠",
			},
		})
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "help_cmds_menu",
					Placeholder: "Selecciona una categoría...",
					Options:     options,
				},
			},
		},
	}
}
