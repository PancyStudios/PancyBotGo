package reaction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

type gifResponse struct {
	URL string `json:"url"`
}

func fetchGifImage(category string) (string, error) {
	url := fmt.Sprintf("https://api.otakugifs.xyz/gif?reaction=%s&format=gif", category)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var data gifResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data.URL, nil
}

func createReactionCommand(name string, actionText string, singleText string, requiresTarget bool) messagecommands.CommandRunFunc {
	return func(ctx *messagecommands.MessageContext) error {
		if requiresTarget && len(ctx.Args) == 0 {
			_, err := ctx.ReplyError("Uso Incorrecto", fmt.Sprintf("Uso: `pan!%s <usuario>`", name))
			return err
		}

		loadingMsg, err := ctx.Reply("Cargando imagen... ⏳")
		if err != nil {
			return err
		}

		imageURL, err := fetchGifImage(name)
		if err != nil {
			ctx.Session.ChannelMessageDelete(loadingMsg.ChannelID, loadingMsg.ID)
			_, err = ctx.ReplyError("Error", "❌ Ocurrió un error al contactar la API de GIFs.")
			return err
		}

		user := ctx.Message.Author
		embed := &discordgo.MessageEmbed{
			Color:     0xFF69B4,
			Image:     &discordgo.MessageEmbedImage{URL: imageURL},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		if requiresTarget {
			parsedID := ctx.ParseUser(0)
			var target *discordgo.User
			if parsedID != "" {
				member, err := ctx.Session.GuildMember(ctx.Message.GuildID, parsedID)
				if err == nil {
					target = member.User
				}
			}

			if target != nil {
				if target.ID == user.ID {
					embed.Description = fmt.Sprintf("❤️ **%s** se %s a sí mismo/a... ¿Todo bien en casa?", user.Username, actionText)
				} else if target.ID == ctx.Session.State.User.ID {
					embed.Description = fmt.Sprintf("❤️ **%s** me %s a mí. ¡Gracias! ✨", user.Username, actionText)
				} else {
					embed.Description = fmt.Sprintf("❤️ **%s** %s a **%s**", user.Username, actionText, target.Username)
				}
			} else {
				embed.Description = fmt.Sprintf("❤️ **%s** %s", user.Username, actionText)
			}
		} else {
			embed.Description = fmt.Sprintf("❤️ **%s** %s", user.Username, singleText)
		}

		content := ""
		_, err = ctx.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:      loadingMsg.ID,
			Channel: loadingMsg.ChannelID,
			Content: &content,
			Embeds:  &[]*discordgo.MessageEmbed{embed},
		})
		return err
	}
}
