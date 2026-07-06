package dev

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func isDev(id string) bool {
	return id == "852683369899622430"
}

func evalCommand(ctx *messagecommands.MessageContext) error {
	if !isDev(ctx.Message.Author.ID) {
		_, err := ctx.ReplyError("Acceso Denegado", "Este comando es solo para la desarrolladora.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!eval <codigo>`")
		return err
	}

	code := strings.Join(ctx.Args, " ")
	code = strings.TrimPrefix(code, "```go")
	code = strings.TrimPrefix(code, "```")
	code = strings.TrimSuffix(code, "```")
	code = strings.TrimSpace(code)

	i := interp.New(interp.Options{})

	if err := i.Use(stdlib.Symbols); err != nil {
		_, err := ctx.ReplyError("Error", fmt.Sprintf("Error cargando stdlib: %v", err))
		return err
	}

	botExports := map[string]reflect.Value{
		"Ctx":     reflect.ValueOf(ctx),
		"Session": reflect.ValueOf(ctx.Session),
		"DB":      reflect.ValueOf(database.Get()),
		"Config":  reflect.ValueOf(config.Get()),
	}

	if err := i.Use(interp.Exports{
		"github.com/PancyStudios/PancyBotGo/internal/messagecommands/dev/dev": botExports,
	}); err != nil {
		_, err := ctx.ReplyError("Error", fmt.Sprintf("Error registrando variables: %v", err))
		return err
	}

	_, err := i.Eval(`import . "github.com/PancyStudios/PancyBotGo/internal/messagecommands/dev"`)
	if err != nil {
		_, err := ctx.ReplyError("Error", fmt.Sprintf("Error importando variables: %v", err))
		return err
	}

	start := time.Now()
	res, err := i.Eval(code)
	duration := time.Since(start)

	var output string
	if err != nil {
		output = fmt.Sprintf("❌ **Error de Ejecución:**\n```go\n%v\n```", err)
	} else {
		resStr := fmt.Sprintf("%+v", res)
		if len(resStr) > 1900 {
			resStr = resStr[:1900] + "..."
		}
		output = fmt.Sprintf("✅ **Resultado:**\n```go\n%s\n```\n⏱️ Tiempo de ejecución: %s", resStr, duration)
	}

	_, err = ctx.Reply(output)
	return err
}
