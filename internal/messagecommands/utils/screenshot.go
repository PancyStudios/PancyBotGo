package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func screenshotCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar la URL de la página web.\nUso: `pan!screenshot <url>`")
		return err
	}

	urlParam := ctx.Args[0]

	channel, err := ctx.Session.Channel(ctx.Message.ChannelID)
	isNsfw := false
	if err == nil && channel != nil {
		isNsfw = channel.NSFW
	}

	baseURL := os.Getenv("SCREENSHOT_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	urlApi := baseURL + "/api/private/screenshot/sfw"
	if isNsfw {
		urlApi = baseURL + "/api/private/screenshot/nsfw"
	}

	payload, _ := json.Marshal(map[string]string{"url": urlParam})
	req, err := http.NewRequest("POST", urlApi, bytes.NewBuffer(payload))
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("Error interno: %v", err))
		return err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("authScreenshots"))
	req.Header.Set("Accept", "image/png")
	req.Header.Set("Content-Type", "application/json")

	user := ctx.Message.Author
	req.Header.Set("X-From-ID", user.ID)

	clientHTTP := &http.Client{Timeout: 64 * time.Second}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("Error conectando con la API: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 520 {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("La API devolvió un código de error: %d", resp.StatusCode))
		return err
	}

	imageBuffer, err := io.ReadAll(resp.Body)
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("Error leyendo imagen: %v", err))
		return err
	}

	file := &discordgo.File{
		Name:        "screenshot.png",
		ContentType: "image/png",
		Reader:      bytes.NewReader(imageBuffer),
	}

	embed := &discordgo.MessageEmbed{
		Title: "Captura de pantalla",
		Color: 0x3498DB,
		Image: &discordgo.MessageEmbedImage{
			URL: "attachment://screenshot.png",
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Solicitado por %s | HTTP: %d", user.String(), resp.StatusCode),
			IconURL: user.AvatarURL(""),
		},
	}

	_, err = ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
		Files:  []*discordgo.File{file},
	})
	return err
}
