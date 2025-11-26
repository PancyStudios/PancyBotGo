// filepath: /home/turbis/GolandProjects/PancyBotGo/internal/commands/utils/ping.go
package utils

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
)

// createPingCommand creates the /utils ping subcommand
func createPingCommand() *discord.Command {
	return discord.NewCommand(
		"ping",
		"Comprueba la latencia del bot",
		"utils",
		pingHandler,
	)
}

// pingHandler handles the /utils ping command
func pingHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		latency := ctx.Client.Session.HeartbeatLatency().Milliseconds()
		ctx.Reply(fmt.Sprintf("🏓 Pong! Latencia: %dms", latency))
	}()
	return nil
}
