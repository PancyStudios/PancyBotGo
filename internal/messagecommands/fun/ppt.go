package fun

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func pptCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes elegir piedra, papel o tijera.\nUso: `pan!ppt <piedra|papel|tijera>`")
		return err
	}

	action := strings.ToLower(ctx.Args[0])
	moves := map[string]int{"piedra": 0, "papel": 1, "tijera": 2}
	moveVals := []string{"piedra", "papel", "tijera"}

	if _, ok := moves[action]; !ok {
		action = "piedra"
	}

	rand.Seed(time.Now().UnixNano())
	machineInput := moveVals[rand.Intn(3)]

	i := moves[action]
	j := moves[machineInput]

	// determine winner: 0 = tie, 1 = user wins, 2 = machine wins
	winner := ((i-j)%3 + 3) % 3

	embed := &discordgo.MessageEmbed{
		Title: "Piedra, papel o tijera",
		Color: 0x5865F2,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("%s eligió", ctx.Message.Author.Username),
				Value:  strings.Title(action),
				Inline: true,
			},
			{
				Name:   "Computadora eligió",
				Value:  strings.Title(machineInput),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    ctx.Message.Author.Username,
			IconURL: ctx.Message.Author.AvatarURL(""),
		},
	}

	if winner == 0 {
		embed.Description = "¡Vaya, hubo un empate!"
	} else if winner == 1 {
		embed.Description = "¡Has ganado, felicidades!"
	} else if winner == 2 {
		embed.Description = "¡La computadora ha ganado, suerte para la próxima!"
	}

	_, err := ctx.ReplyEmbed(embed)
	return err
}
