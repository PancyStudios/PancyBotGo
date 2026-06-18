package ia

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

func createGetImageCommand() *discord.Command {
	return discord.NewCommand(
		"getimage",
		"Pide una imagen generada anteriormente",
		"ia",
		getImageHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "id",
			Description: "ID de la imagen",
			Required:    true,
		},
	)
}

func getImageHandler(ctx *discord.CommandContext) error {
	id := ctx.GetStringOption("id")

			dbUrl := os.Getenv("IMAGE_DB_URL")
			if dbUrl == "" {
				sendError(ctx, "IMAGE_DB_URL no está configurado.")
				return nil
			}

			// Format: IMAGE_DB_URL + "image/craiyon/craiyon" + id + ".png"
			targetUrl := fmt.Sprintf("%simage/craiyon/craiyon%s.png", dbUrl, id)

			// Fast check if it exists
			resp, err := http.Head(targetUrl)
			if err != nil || (resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 304) {
				// The TS version also checked local files, but since Go is stateless in this context
				// we will just return that it doesn't exist.
				sendError(ctx, "No existe esta imagen en la base de datos de imágenes.")
				return nil
			}

			embed := &discordgo.MessageEmbed{
				Title:     "Imagen solicitada",
				Color:     0xff0000,
				Image:     &discordgo.MessageEmbedImage{URL: targetUrl},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			err = ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			if err != nil {
				logger.Error(fmt.Sprintf("Error enviando getimage: %v", err), "IA")
			}
			return nil
}

func sendError(ctx *discord.CommandContext, msg string) {
	err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error enviando error en getimage: %v", err), "IA")
	}
}
