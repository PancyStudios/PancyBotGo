package ia

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

func createCreateImageCommand() *discord.Command {
	return discord.NewCommand(
		"createimage",
		"🎨 | Genera una imagen con IA",
		"ia",
		createImageHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "prompt",
			Description: "🤖 | Descripción de la imagen a generar",
			Required:    true,
		},
	)
}

func createImageHandler(ctx *discord.CommandContext) error {
	prompt := ctx.GetStringOption("prompt")

	// Initial reply
	err := ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Generando... ⏳",
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error sending initial response: %v", err), "IA")
		return err
	}

	start := time.Now()

	// Use local API
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
		sendErrorEdit(ctx, fmt.Sprintf("Error interno: %v", err))
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("authScreenshots"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 4 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		sendErrorEdit(ctx, fmt.Sprintf("Error de red: %v", err))
		return nil
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		sendErrorEdit(ctx, fmt.Sprintf("Error leyendo respuesta: %v", err))
		return nil
	}

	if resp.StatusCode != 200 {
		sendErrorEdit(ctx, fmt.Sprintf("Error de la IA (Estado: %d): %s", resp.StatusCode, string(bodyBytes)))
		return nil
	}

	var jsonRes map[string]string
	if err := json.Unmarshal(bodyBytes, &jsonRes); err != nil {
		sendErrorEdit(ctx, "Respuesta inválida del microservicio")
		return nil
	}

	imageUrl := jsonRes["image_url"]
	if imageUrl == "" {
		sendErrorEdit(ctx, "No se generó ninguna imagen")
		return nil
	}

	embed := &discordgo.MessageEmbed{
		Color: 0xff0000,
		Title: "🖼️ Imagen Generada",
		Description: fmt.Sprintf("**Prompt:** %s", prompt),
		Image: &discordgo.MessageEmbedImage{
			URL: imageUrl,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Tiempo de generación: %v | Craiyon v3", time.Since(start).Round(time.Millisecond)),
			IconURL: ctx.User().AvatarURL(""),
		},
	}

	content := "Imagen generada exitosamente ✨"
	_, err = ctx.Session.InteractionResponseEdit(ctx.Interaction.Interaction, &discordgo.WebhookEdit{
		Content: &content,
		Embeds:  &[]*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error editing response: %v", err), "IA")
	}
	return nil
}

func sendErrorEdit(ctx *discord.CommandContext, errMsg string) {
	embed := &discordgo.MessageEmbed{
		Title:       "Craiyon Error",
		Description: fmt.Sprintf("Error: %s", errMsg),
		Color:       0xff0000,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	content := ""
	_, err := ctx.Session.InteractionResponseEdit(ctx.Interaction.Interaction, &discordgo.WebhookEdit{
		Content: &content,
		Embeds:  &[]*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error editing error response: %v", err), "IA")
	}
}
