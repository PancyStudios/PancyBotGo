package ia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func createImageCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar la descripción de la imagen.\nUso: `pan!createimage <prompt>`")
		return err
	}

	prompt := strings.Join(ctx.Args, " ")

	loadingMsg, err := ctx.Reply("Generando... ⏳")
	if err != nil {
		return err
	}

	start := time.Now()

	baseURL := os.Getenv("SCREENSHOT_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	modelUrl := baseURL + "/api/private/ia/fetch"

	reqBody := map[string]interface{}{
		"prompt": prompt,
	}
	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", modelUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.Session.ChannelMessageDelete(loadingMsg.ChannelID, loadingMsg.ID)
		_, err = ctx.ReplyError("Error", fmt.Sprintf("Error interno: %v", err))
		return err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("authScreenshots"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 4 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		ctx.Session.ChannelMessageDelete(loadingMsg.ChannelID, loadingMsg.ID)
		_, err = ctx.ReplyError("Error", fmt.Sprintf("Error de red: %v", err))
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.Session.ChannelMessageDelete(loadingMsg.ChannelID, loadingMsg.ID)
		_, err = ctx.ReplyError("Error", fmt.Sprintf("Error leyendo respuesta: %v", err))
		return err
	}

	if resp.StatusCode != 200 {
		ctx.Session.ChannelMessageDelete(loadingMsg.ChannelID, loadingMsg.ID)
		_, err = ctx.ReplyError("Error", fmt.Sprintf("Error de la IA (Estado: %d): %s", resp.StatusCode, string(bodyBytes)))
		return err
	}

	var jsonRes map[string]string
	if err := json.Unmarshal(bodyBytes, &jsonRes); err != nil {
		ctx.Session.ChannelMessageDelete(loadingMsg.ChannelID, loadingMsg.ID)
		_, err = ctx.ReplyError("Error", "Respuesta inválida del microservicio")
		return err
	}

	imageUrl := jsonRes["image_url"]
	if imageUrl == "" {
		ctx.Session.ChannelMessageDelete(loadingMsg.ChannelID, loadingMsg.ID)
		_, err = ctx.ReplyError("Error", "No se generó ninguna imagen")
		return err
	}

	embed := &discordgo.MessageEmbed{
		Color:       0xff0000,
		Title:       "🖼️ Imagen Generada",
		Description: fmt.Sprintf("**Prompt:** %s", prompt),
		Image: &discordgo.MessageEmbedImage{
			URL: imageUrl,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Tiempo de generación: %v | Craiyon v3", time.Since(start).Round(time.Millisecond)),
			IconURL: ctx.Message.Author.AvatarURL(""),
		},
	}

	embeds := []*discordgo.MessageEmbed{embed}
	content := "Imagen generada exitosamente ✨"
	_, err = ctx.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      loadingMsg.ID,
		Channel: loadingMsg.ChannelID,
		Content: &content,
		Embeds:  &embeds,
	})
	return err
}
