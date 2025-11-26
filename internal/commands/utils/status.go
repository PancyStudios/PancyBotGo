// filepath: /home/turbis/GolandProjects/PancyBotGo/internal/commands/utils/status.go
package utils

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
)

// createStatusCommand creates the /utils status subcommand
func createStatusCommand() *discord.Command {
	return discord.NewCommand(
		"status",
		"Muestra el estado del bot",
		"utils",
		statusHandler,
	)
}

// statusHandler handles the /utils status command
func statusHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		db := database.Get()
		dbStatus, _ := db.GetStatus()

		ctx.Reply(fmt.Sprintf(
			"📊 **Estado del Bot**\n"+
				"• Bot: 🟢 Online\n"+
				"• Base de datos: %s\n"+
				"• Servidores: %d",
			dbStatus,
			ctx.Client.GuildCount(),
		))
	}()
	return nil
}
