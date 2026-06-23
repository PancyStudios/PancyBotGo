package reaction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

type waifuResponse struct {
	URL string `json:"url"`
}

func fetchWaifuImage(category string) (string, error) {
	url := fmt.Sprintf("https://api.waifu.pics/sfw/%s", category)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var data waifuResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data.URL, nil
}

func createReactionCommand(name string, description string, actionText string, singleText string, requiresTarget bool) *discord.Command {
	cmd := &discord.Command{
		Name:        name,
		Description: description,
	}

	if requiresTarget {
		cmd.Options = []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "usuario",
				Description: "🎭 | El usuario con el que vas a interactuar",
				Required:    true,
			},
		}
	}

	cmd.Run = func(ctx *discord.CommandContext) error {
		ctx.Defer()

		imageURL, err := fetchWaifuImage(name)
		if err != nil {
			return ctx.EditReply("❌ Ocurrió un error al contactar la API de GIFs.")
		}

		user := ctx.User()
		embed := &discordgo.MessageEmbed{
			Color:     0xFF69B4, // Hot Pink
			Image:     &discordgo.MessageEmbedImage{URL: imageURL},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		if requiresTarget {
			target := ctx.GetUserOption("usuario")
			if target != nil {
				if target.ID == user.ID {
					embed.Description = fmt.Sprintf("❤️ **%s** se %s a sí mismo/a... ¿Todo bien en casa?", user.Username, actionText)
				} else if target.ID == ctx.Session.State.User.ID {
					embed.Description = fmt.Sprintf("❤️ **%s** me %s a mí. ¡Gracias! ✨", user.Username, actionText)
				} else {
					embed.Description = fmt.Sprintf("❤️ **%s** %s a **%s**", user.Username, actionText, target.Username)
				}
			}
		} else {
			embed.Description = fmt.Sprintf("❤️ **%s** %s", user.Username, singleText)
		}

		return ctx.EditReplyEmbed(embed)
	}

	return cmd
}
