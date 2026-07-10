package utils

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createAvatarCommand() *discord.Command {
	return &discord.Command{
		Name:        "avatar",
		Description: "🧰 | Muestra el avatar de un usuario",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Usuario del cual mostrar el avatar",
				Required:    false,
			},
		},
		Run: func(ctx *discord.CommandContext) error {
			user := ctx.GetUserOption("user")
			if user == nil {
				user = ctx.User()
			}

			text := fmt.Sprintf("Avatar de: %s\n[Descargar Avatar](%s)", user.Username, user.AvatarURL(""))

			embed := &discordgo.MessageEmbed{
				Description: text,
				Image: &discordgo.MessageEmbedImage{
					URL: user.AvatarURL(""),
				},
				Color: 0x00FF00, // Green
			}

			return ctx.ReplyEmbed(embed)
		},
	}
}
