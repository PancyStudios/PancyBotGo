package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

func createScreenshotCommand() *discord.Command {
	return &discord.Command{
		Name:        "screenshot",
		Description: "🧰 | Toma una captura de pantalla de una página web",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "🧰 | La url de la página web",
				Required:    true,
			},
		},
		Run: func(ctx *discord.CommandContext) error {
			urlParam := ctx.GetStringOption("url")

			// Defer response
			err := ctx.Defer()
			if err != nil {
				logger.Error(fmt.Sprintf("Error deferring response: %v", err), "Screenshot")
				return err
			}

			// Get channel to check NSFW
			channel := ctx.Channel()

			isNsfw := channel != nil && channel.NSFW
			
			// Use local screenshot API or an env variable
			baseURL := os.Getenv("SCREENSHOT_API_URL")
			if baseURL == "" {
				baseURL = "http://localhost:3000"
			}
			
			urlApi := baseURL + "/api/private/screenshot/sfw"
			if isNsfw {
				urlApi = baseURL + "/api/private/screenshot/nsfw"
			}

			// Prepare request payload
			payload, _ := json.Marshal(map[string]string{"url": urlParam})
			req, err := http.NewRequest("POST", urlApi, bytes.NewBuffer(payload))
			if err != nil {
				return sendError(ctx, fmt.Sprintf("Error interno: %v", err))
			}

			req.Header.Set("Authorization", "Bearer "+os.Getenv("authScreenshots"))
			req.Header.Set("Accept", "image/png")
			req.Header.Set("Content-Type", "application/json")
			
			user := ctx.User()
			req.Header.Set("X-From-ID", user.ID)

			clientHTTP := &http.Client{Timeout: 64 * time.Second}
			resp, err := clientHTTP.Do(req)
			if err != nil {
				return sendError(ctx, fmt.Sprintf("Error conectando con la API: %v", err))
			}
			defer resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode >= 520 {
				return sendError(ctx, fmt.Sprintf("La API devolvió un código de error: %d", resp.StatusCode))
			}

			imageBuffer, err := io.ReadAll(resp.Body)
			if err != nil {
				return sendError(ctx, fmt.Sprintf("Error leyendo imagen: %v", err))
			}

			file := &discordgo.File{
				Name:        "screenshot.png",
				ContentType: "image/png",
				Reader:      bytes.NewReader(imageBuffer),
			}

			embed := &discordgo.MessageEmbed{
				Title: "Captura de pantalla",
				Color: 0x3498DB, // Blue
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://screenshot.png",
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    fmt.Sprintf("Solicitado por %s | HTTP: %d", user.String(), resp.StatusCode),
					IconURL: user.AvatarURL(""),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err = ctx.Session.InteractionResponseEdit(ctx.Interaction.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{embed},
				Files:  []*discordgo.File{file},
			})
			if err != nil {
				logger.Error(fmt.Sprintf("Error editing response: %v", err), "Screenshot")
			}
			return err
		},
	}
}

func sendError(ctx *discord.CommandContext, errorMsg string) error {
	embed := &discordgo.MessageEmbed{
		Title:       "Error",
		Description: "🧰 | No se pudo tomar la captura de pantalla, verifica la url o intenta nuevamente.\n\n" + errorMsg,
		Color:       0xFF0000, // Red
	}
	
	user := ctx.User()
	
	if user != nil {
	    embed.Footer = &discordgo.MessageEmbedFooter{
		    Text:    fmt.Sprintf("Solicitado por %s", user.String()),
		    IconURL: user.AvatarURL(""),
	    }
	}
	
	embed.Timestamp = time.Now().Format(time.RFC3339)

	return ctx.EditReplyEmbed(embed)
}
