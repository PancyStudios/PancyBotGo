package fun

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

var eightBallResponses = []string{
	"En mi opinión, sí",
	"Es cierto",
	"Es decididamente así",
	"Probablemente",
	"Buen pronóstico",
	"Todo apunta a que sí",
	"Sin duda",
	"Sí",
	"Sí - definitivamente",
	"Debes confiar en ello",
	"Respuesta vaga, vuelve a intentarlo",
	"Pregunta en otro momento",
	"Será mejor que no te lo diga ahora",
	"No puedo predecirlo ahora",
	"Concéntrate y vuelve a preguntar",
	"No cuentes con ello",
	"Mi respuesta es no",
	"Mis fuentes me dicen que no",
	"Las perspectivas no son buenas",
	"Muy dudoso",
}

func eightBallCommand(ctx *messagecommands.MessageContext) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes hacer una pregunta.\nUso: `pan!8ball <pregunta>`")
		return err
	}

	pregunta := strings.Join(ctx.Args, " ")

	rand.Seed(time.Now().UnixNano())
	respuesta := eightBallResponses[rand.Intn(len(eightBallResponses))]

	_, err := ctx.Reply(fmt.Sprintf("🎱 **Pregunta:** %s\n💬 **Respuesta:** %s", pregunta, respuesta))
	return err
}
