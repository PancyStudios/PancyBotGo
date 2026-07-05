package levels

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

var toggleCommand = &discord.Command{
	Name:            "toggle",
	Description:     "⚙️ | Activa o desactiva el sistema de niveles en el servidor",
	UserPermissions: discordgo.PermissionAdministrator,
	Run: func(ctx *discord.CommandContext) error {
		guildID := ctx.Interaction.GuildID
		if guildID == "" {
			return ctx.ReplyEphemeral("❌ Este comando solo se puede usar en un servidor.")
		}

		// Obtener configuración del servidor
		guildData, err := database.GlobalGuildDM.Get(bson.M{"id": guildID})
		if err != nil {
			return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error al obtener la configuración del servidor: %v", err))
		}

		// Alternar estado
		newState := !guildData.Levels.Enable
		guildData.Levels.Enable = newState

		// Guardar configuración
		_, err = database.GlobalGuildDM.Set(bson.M{"id": guildID}, guildData)
		if err != nil {
			return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error al guardar la configuración: %v", err))
		}

		if newState {
			return ctx.Reply("✅ **Sistema de Niveles Activado.**\n\nLos usuarios ahora ganarán experiencia al chatear. Puedes ver la tabla de clasificación con `/levels leaderboard` o en el panel web.")
		}

		return ctx.Reply("✅ **Sistema de Niveles Desactivado.**\n\nLos usuarios ya no ganarán experiencia al chatear. La experiencia acumulada no se perderá.")
	},
}
