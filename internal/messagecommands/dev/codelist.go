package dev

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/bwmarrin/discordgo"
)

func codelistCommand(ctx *messagecommands.MessageContext) error {
	if !isDev(ctx.Message.Author.ID) {
		_, err := ctx.ReplyError("Acceso Denegado", "Este comando es solo para la desarrolladora.")
		return err
	}

	codes, err := database.GetAllPremiumCodes()
	if err != nil {
		_, err = ctx.ReplyError("Error", "No se pudo obtener la lista de códigos.")
		return err
	}

	if len(codes) == 0 {
		_, err = ctx.ReplySuccess("Códigos Premium", "No hay códigos premium generados.")
		return err
	}

	var descBuilder strings.Builder
	for _, code := range codes {
		status := "🟢 Disponible"
		if code.IsClaimed {
			status = "🔴 Canjeado"
		}
		
		descBuilder.WriteString(fmt.Sprintf("**%s** - %s - %s\n", code.Code, string(code.Type), status))
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📋 Lista de Códigos Premium",
		Description: descBuilder.String(),
		Color:       0x00FF00,
	}

	_, err = ctx.ReplyEmbed(embed)
	return err
}
