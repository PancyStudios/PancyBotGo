package ia

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func getImageCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar el ID de la imagen.\nUso: `pan!getimage <id>`")
		return err
	}

	id := ctx.Args[0]

	dbUrl := os.Getenv("IMAGE_DB_URL")
	if dbUrl == "" {
		_, err := ctx.ReplyError("Error", "IMAGE_DB_URL no está configurado.")
		return err
	}

	targetUrl := fmt.Sprintf("%simage/craiyon/craiyon%s.png", dbUrl, id)

	resp, err := http.Head(targetUrl)
	if err != nil || (resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 304) {
		_, err = ctx.ReplyError("Error", "No existe esta imagen en la base de datos de imágenes.")
		return err
	}

	embed := &discordgo.MessageEmbed{
		Title:     "Imagen solicitada",
		Color:     0xff0000,
		Image:     &discordgo.MessageEmbedImage{URL: targetUrl},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = ctx.ReplyEmbed(embed)
	return err
}
