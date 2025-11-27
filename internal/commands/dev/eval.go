package dev

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// CreateEvalCommand crea el comando /dev eval
func CreateEvalCommand() *discord.Command {
	return discord.NewCommand(
		"eval",
		"Evalúa código Go y muestra estructuras internas (Peligroso)",
		"dev",
		evalHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "codigo",
			Description: "Código o expresión Go a evaluar",
			Required:    true,
		},
	)
}

func evalHandler(ctx *discord.CommandContext) error {
	go func() {
		errors.RecoverMiddleware()()
		start := time.Now()
		// 1. Seguridad: Validación estricta de ID (Ximena)
		// Usamos el mismo ID que tienes en codegen.go
		if !isDev(ctx.User().ID) {
			ctx.ReplyEphemeral("❌ **Acceso Denegado:** Este comando es solo para la desarrolladora.")
			return
		}

		// Deferimos la respuesta porque compilar el script puede tomar unos milisegundos
		ctx.Defer()

		// 2. Limpieza del código de entrada
		code := ctx.GetStringOption("codigo")
		// Eliminar bloques de código de markdown si existen (```go ... ```)
		code = strings.TrimPrefix(code, "```go")
		code = strings.TrimPrefix(code, "```")
		code = strings.TrimSuffix(code, "```")
		code = strings.TrimSpace(code)

		// 3. Inicializar el intérprete Yaegi
		i := interp.New(interp.Options{})

		// Cargar librería estándar de Go (fmt, strings, os, etc.)
		if err := i.Use(stdlib.Symbols); err != nil {
			ctx.EditReply(fmt.Sprintf("❌ Error cargando stdlib: %v", err))
			return
		}

		// 4. Inyección de variables del Bot usando Exports
		// Esto hace que puedas usar 'DB', 'Ctx', 'Bot' directamente en tu código
		botExports := map[string]reflect.Value{
			"Ctx":     reflect.ValueOf(ctx),
			"Bot":     reflect.ValueOf(ctx.Client),
			"Session": reflect.ValueOf(ctx.Session),
			"DB":      reflect.ValueOf(database.Get()),
			"Config":  reflect.ValueOf(config.Get()),
		}

		// Registrar las variables como símbolos globales
		if err := i.Use(interp.Exports{
			"github.com/PancyStudios/PancyBotGo/internal/commands/dev/dev": botExports,
		}); err != nil {
			ctx.EditReply(fmt.Sprintf("❌ Error registrando variables: %v", err))
			return
		}

		// Importar el paquete automáticamente
		_, err := i.Eval(`import . "github.com/PancyStudios/PancyBotGo/internal/commands/dev"`)
		if err != nil {
			ctx.EditReply(fmt.Sprintf("❌ Error importando variables: %v", err))
			return
		}

		// 5. Ejecutar el código del usuario
		res, err := i.Eval(code)

		// 6. Formatear la respuesta
		var output string
		if err != nil {
			output = fmt.Sprintf("❌ **Error de Ejecución:**\n```go\n%v\n```", err)
		} else {
			// Intentamos formatear el resultado de forma bonita
			var resStr string
			if res.IsValid() {
				// Usamos %#v para ver la estructura interna completa (campos, punteros, etc)
				resStr = fmt.Sprintf("%#v", res.Interface())
			} else {
				resStr = "nil"
			}
			if len(resStr) > 1900 {
				resStr = resStr[:1900] + "... (truncado)"
			}

			output = fmt.Sprintf("✅ **Resultado:**\n```go\n%s\n```", resStr)
		}

		elapsed := time.Since(start)
		logger.Debug(fmt.Sprintf("Eval completado en %s. Limpiando contexto Yaegi...", elapsed), "DevEval")

		ctx.EditReply(output)
		return
	}()
	return nil
}

// Helper para verificar ID (Hardcoded por seguridad como en tu codegen.go)
func isDev(userID string) bool {
	return userID == "852683369899622430"
}
