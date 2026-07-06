package fun

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

type DogResponse struct {
	URL string `json:"url"`
}

func dogCommand(ctx *messagecommands.MessageContext) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://nekos.life/api/v2/img/woof")
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ No se pudo conectar a la API.")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		_, err = ctx.ReplyError("Error", "❌ La API devolvió un error.")
		return err
	}

	body, _ := io.ReadAll(resp.Body)
	var dogResp DogResponse
	if err := json.Unmarshal(body, &dogResp); err != nil {
		_, err = ctx.ReplyError("Error", "❌ Error al procesar la imagen.")
		return err
	}

	embed := &discordgo.MessageEmbed{
		Title: "🐕",
		Color: 0xF8A269,
		Image: &discordgo.MessageEmbedImage{
			URL: dogResp.URL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    ctx.Message.Author.Username,
			IconURL: ctx.Message.Author.AvatarURL(""),
		},
	}

	_, err = ctx.ReplyEmbed(embed)
	return err
}
