package embeds

import (
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func createEmbedCommand(ctx *messagecommands.MessageContext) error {
	user := ctx.Message.Author

	embedState := getBuilderState(user.ID)
	saveBuilderState(user.ID, embedState)

	menu := discordgo.SelectMenu{
		CustomID:    "embed_builder_menu",
		Placeholder: "Selecciona el campo a editar...",
		Options: []discordgo.SelectMenuOption{
			{Label: "Título", Value: "title", Description: "📝 | Edita el título del embed", Emoji: &discordgo.ComponentEmoji{Name: "📝"}},
			{Label: "Descripción", Value: "description", Description: "📝 | Edita la descripción principal", Emoji: &discordgo.ComponentEmoji{Name: "📄"}},
			{Label: "Color", Value: "color", Description: "📝 | Edita el color lateral (Hex)", Emoji: &discordgo.ComponentEmoji{Name: "🎨"}},
			{Label: "Autor (Nombre)", Value: "author_name", Description: "📝 | Edita el nombre del autor", Emoji: &discordgo.ComponentEmoji{Name: "👤"}},
			{Label: "Autor (Icono)", Value: "author_icon", Description: "📝 | URL de la imagen del autor", Emoji: &discordgo.ComponentEmoji{Name: "🖼️"}},
			{Label: "Footer (Texto)", Value: "footer_text", Description: "📝 | Edita el texto al pie", Emoji: &discordgo.ComponentEmoji{Name: "🔻"}},
			{Label: "Footer (Icono)", Value: "footer_icon", Description: "📝 | URL de la imagen del pie", Emoji: &discordgo.ComponentEmoji{Name: "🖼️"}},
			{Label: "Imagen Principal", Value: "image", Description: "📝 | URL de la imagen grande", Emoji: &discordgo.ComponentEmoji{Name: "🏞️"}},
			{Label: "Thumbnail", Value: "thumbnail", Description: "📝 | URL de la imagen pequeña", Emoji: &discordgo.ComponentEmoji{Name: "🖼️"}},
		},
	}

	buttonsRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Finalizar y Guardar",
				Style:    discordgo.SuccessButton,
				CustomID: "embed_builder_save",
				Emoji:    &discordgo.ComponentEmoji{Name: "💾"},
			},
			discordgo.Button{
				Label:    "Descartar",
				Style:    discordgo.DangerButton,
				CustomID: "embed_builder_cancel",
				Emoji:    &discordgo.ComponentEmoji{Name: "🗑️"},
			},
		},
	}

	menuRow := discordgo.ActionsRow{Components: []discordgo.MessageComponent{menu}}

	_, err := ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Content: "🛠️ **Constructor de Embeds**\nUtiliza el menú de abajo para modificar las propiedades. Presiona `Finalizar` cuando termines.",
		Embeds:  []*discordgo.MessageEmbed{embedState},
		Components: []discordgo.MessageComponent{
			menuRow,
			buttonsRow,
		},
	})
	return err
}
