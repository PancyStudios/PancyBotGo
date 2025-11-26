// filepath: /home/turbis/GolandProjects/PancyBotGo/internal/commands/utils/help.go
package utils

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
)

// createHelpCommand creates the /utils help subcommand
func createHelpCommand() *discord.Command {
	return discord.NewCommand(
		"help",
		"Muestra información de ayuda",
		"utils",
		helpHandler,
	)
}

// helpHandler handles the /utils help command
func helpHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		ctx.Reply(
			"📖 **Ayuda de PancyBot Go**\n\n" +
				"**Comandos disponibles:**\n" +
				"• `/utils ping` - Comprueba la latencia\n" +
				"• `/utils status` - Estado del bot\n" +
				"• `/utils stats` - Estadísticas del bot\n" +
				"• `/play <query>` - Reproduce música\n" +
				"• `/pause` - Pausa/resume la música\n" +
				"• `/skip` - Salta la canción actual\n" +
				"• `/stop` - Detiene la música\n" +
				"• `/queue` - Muestra la cola\n" +
				"• `/volume <0-100>` - Ajusta el volumen\n" +
				"• `/mod ban <usuario> <razón>` - Banea a un usuario\n" +
				"• `/mod kick <usuario> <razón>` - Expulsa a un usuario\n" +
				"• `/mod warn <usuario> <razón>` - Advierte a un usuario\n" +
				"• `/mod mute <usuario> <duración> <razón>` - Mutea a un usuario\n" +
				"• `/mod warnings <usuario>` - Lista las advertencias\n" +
				"• `/mod removewarn <usuario> <id>` - Elimina una advertencia",
		)
	}()
	return nil
}
