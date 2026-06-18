package embeds

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createEmbedCreateCommand() *discord.Command {
	return &discord.Command{
		Name:        "create",
		Description: "Abre el constructor interactivo de embeds",
		Run: func(ctx *discord.CommandContext) error {
			user := ctx.User()
			
			embedState := getBuilderState(user.ID)
			saveBuilderState(user.ID, embedState)

			// Create the select menu
			menu := discordgo.SelectMenu{
				CustomID:    "embed_builder_menu",
				Placeholder: "Selecciona el campo a editar...",
				Options: []discordgo.SelectMenuOption{
					{Label: "Título", Value: "title", Description: "Edita el título del embed", Emoji: &discordgo.ComponentEmoji{Name: "📝"}},
					{Label: "Descripción", Value: "description", Description: "Edita la descripción principal", Emoji: &discordgo.ComponentEmoji{Name: "📄"}},
					{Label: "Color", Value: "color", Description: "Edita el color lateral (Hex)", Emoji: &discordgo.ComponentEmoji{Name: "🎨"}},
					{Label: "Autor (Nombre)", Value: "author_name", Description: "Edita el nombre del autor", Emoji: &discordgo.ComponentEmoji{Name: "👤"}},
					{Label: "Autor (Icono)", Value: "author_icon", Description: "URL de la imagen del autor", Emoji: &discordgo.ComponentEmoji{Name: "🖼️"}},
					{Label: "Footer (Texto)", Value: "footer_text", Description: "Edita el texto al pie", Emoji: &discordgo.ComponentEmoji{Name: "🔻"}},
					{Label: "Footer (Icono)", Value: "footer_icon", Description: "URL de la imagen del pie", Emoji: &discordgo.ComponentEmoji{Name: "🖼️"}},
					{Label: "Imagen Principal", Value: "image", Description: "URL de la imagen grande", Emoji: &discordgo.ComponentEmoji{Name: "🏞️"}},
					{Label: "Thumbnail", Value: "thumbnail", Description: "URL de la imagen pequeña", Emoji: &discordgo.ComponentEmoji{Name: "🖼️"}},
				},
			}

			// Save/Cancel buttons
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

			return ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "🛠️ **Constructor de Embeds**\nUtiliza el menú de abajo para modificar las propiedades. Presiona `Finalizar` cuando termines.",
					Embeds:  []*discordgo.MessageEmbed{embedState},
					Components: []discordgo.MessageComponent{
						menuRow,
						buttonsRow,
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
		},
	}
}
