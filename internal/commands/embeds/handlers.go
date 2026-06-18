package embeds

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HandleInteraction processes interactions for the embed builder
// Returns true if the interaction was handled by this module
func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if i.Type == discordgo.InteractionMessageComponent {
		customID := i.MessageComponentData().CustomID

		if customID == "embed_builder_menu" {
			handleMenuSelect(s, i)
			return true
		} else if customID == "embed_builder_cancel" {
			
			user := i.Member.User
			if user == nil && i.User != nil {
			    user = i.User
			}
			
			clearBuilderState(user.ID)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content:    "🗑️ Creación de embed descartada.",
					Embeds:     nil,
					Components: []discordgo.MessageComponent{},
				},
			})
			return true
		} else if customID == "embed_builder_save" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: "💾 ¡Embed guardado temporalmente en tu sesión!\nUsa `/embed send` para enviarlo a un canal.",
					Components: []discordgo.MessageComponent{},
				},
			})
			return true
		}
	} else if i.Type == discordgo.InteractionModalSubmit {
		customID := i.ModalSubmitData().CustomID
		if strings.HasPrefix(customID, "embed_modal_") {
			handleModalSubmit(s, i)
			return true
		}
	}
	return false
}

func handleMenuSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	if len(data.Values) == 0 {
		return
	}
	selected := data.Values[0]

	user := i.Member.User
	if user == nil && i.User != nil {
	    user = i.User
	}
	
	embedState := getBuilderState(user.ID)

	var title string
	var value string
	var style discordgo.TextInputStyle = discordgo.TextInputShort

	switch selected {
	case "title":
		title = "Editar Título"
		value = embedState.Title
	case "description":
		title = "Editar Descripción"
		value = embedState.Description
		style = discordgo.TextInputParagraph
	case "color":
		title = "Editar Color (Hexadecimal)"
		value = fmt.Sprintf("%06X", embedState.Color)
	case "author_name":
		title = "Editar Nombre del Autor"
		if embedState.Author != nil {
			value = embedState.Author.Name
		}
	case "author_icon":
		title = "Editar Icono del Autor (URL)"
		if embedState.Author != nil {
			value = embedState.Author.IconURL
		}
	case "footer_text":
		title = "Editar Texto del Pie"
		if embedState.Footer != nil {
			value = embedState.Footer.Text
		}
	case "footer_icon":
		title = "Editar Icono del Pie (URL)"
		if embedState.Footer != nil {
			value = embedState.Footer.IconURL
		}
	case "image":
		title = "Editar Imagen Principal (URL)"
		if embedState.Image != nil {
			value = embedState.Image.URL
		}
	case "thumbnail":
		title = "Editar Miniatura (URL)"
		if embedState.Thumbnail != nil {
			value = embedState.Thumbnail.URL
		}
	}

	textInput := discordgo.TextInput{
		CustomID:    "value",
		Label:       "Nuevo valor",
		Style:       style,
		Placeholder: "Escribe aquí...",
		Value:       value,
		Required:    false,
	}

	modal := &discordgo.InteractionResponseData{
		CustomID: "embed_modal_" + selected,
		Title:    title,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{textInput},
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: modal,
	})
}

func handleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	field := strings.TrimPrefix(data.CustomID, "embed_modal_")

	user := i.Member.User
	if user == nil && i.User != nil {
	    user = i.User
	}
	
	embedState := getBuilderState(user.ID)

	newValue := ""
	for _, comp := range data.Components {
		if row, ok := comp.(*discordgo.ActionsRow); ok {
			for _, c := range row.Components {
				if ti, ok := c.(*discordgo.TextInput); ok {
					newValue = ti.Value
					break
				}
			}
		}
	}

	switch field {
	case "title":
		embedState.Title = newValue
	case "description":
		embedState.Description = newValue
	case "color":
		newValue = strings.TrimPrefix(newValue, "#")
		if colorInt, err := strconv.ParseInt(newValue, 16, 64); err == nil {
			embedState.Color = int(colorInt)
		}
	case "author_name":
		if embedState.Author == nil {
			embedState.Author = &discordgo.MessageEmbedAuthor{}
		}
		embedState.Author.Name = newValue
	case "author_icon":
		if embedState.Author == nil {
			embedState.Author = &discordgo.MessageEmbedAuthor{}
		}
		embedState.Author.IconURL = newValue
	case "footer_text":
		if embedState.Footer == nil {
			embedState.Footer = &discordgo.MessageEmbedFooter{}
		}
		embedState.Footer.Text = newValue
	case "footer_icon":
		if embedState.Footer == nil {
			embedState.Footer = &discordgo.MessageEmbedFooter{}
		}
		embedState.Footer.IconURL = newValue
	case "image":
		if newValue == "" {
			embedState.Image = nil
		} else {
			if embedState.Image == nil {
				embedState.Image = &discordgo.MessageEmbedImage{}
			}
			embedState.Image.URL = newValue
		}
	case "thumbnail":
		if newValue == "" {
			embedState.Thumbnail = nil
		} else {
			if embedState.Thumbnail == nil {
				embedState.Thumbnail = &discordgo.MessageEmbedThumbnail{}
			}
			embedState.Thumbnail.URL = newValue
		}
	}

	saveBuilderState(user.ID, embedState)

	// Update the message with the new embed
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embedState},
		},
	})
}
