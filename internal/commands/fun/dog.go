package fun

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

type DogResponse struct {
	URL string `json:"url"`
}

func createDogCommand() *discord.Command {
	return discord.NewCommand(
		"dog",
		"🐶 | Muestra una imagen de un perro",
		"fun",
		dogHandler,
	)
}

func dogHandler(ctx *discord.CommandContext) error {
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Get("https://nekos.life/api/v2/img/woof")
			if err != nil {
				sendErrorDog(ctx, "No se pudo conectar a la API.")
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				sendErrorDog(ctx, "La API devolvió un error.")
				return nil
			}

			body, _ := io.ReadAll(resp.Body)
			var dogResp DogResponse
			if err := json.Unmarshal(body, &dogResp); err != nil {
				sendErrorDog(ctx, "Error al procesar la imagen.")
				return err
			}

			embed := &discordgo.MessageEmbed{
				Title: "🐕",
				Color: 0xF8A269,
				Image: &discordgo.MessageEmbedImage{
					URL: dogResp.URL,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    ctx.User().Username,
					IconURL: ctx.User().AvatarURL(""),
				},
			}

			err = ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando dog: %v", err), "Fun")
			}
			return nil
}

func sendErrorDog(ctx *discord.CommandContext, msg string) {
	err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("❌ %s", msg),
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error respondiendo a dog: %v", err), "Fun")
	}
}
