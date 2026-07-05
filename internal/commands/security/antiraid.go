package security

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func createAntiraidCommand() *discord.Command {
	return discord.NewCommand(
		"antiraid",
		"🛡️ | Configura el sistema Anti-Raid del servidor",
		"security",
		antiraidHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "toggle",
			Description: "🚨 | Activa o desactiva el modo pánico Anti-Raid",
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "age",
			Description: "⏳ | Configura la edad mínima (en días) de una cuenta para unirse",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "dias",
					Description: "Edad mínima en días (0 para desactivar)",
					Required:    true,
				},
			},
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "limits",
			Description: "📊 | Configura el detector automático de raids",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "uniones",
					Description: "Cantidad de uniones permitidas (0 para desactivar)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "segundos",
					Description: "Ventana de tiempo en segundos",
					Required:    true,
				},
			},
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "action",
			Description: "⚡ | Acción a tomar con los asaltantes",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tipo",
					Description: "Acción",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Expulsar (Kick)", Value: "kick"},
						{Name: "Banear (Ban)", Value: "ban"},
					},
				},
			},
		},
	).WithUserPermissions(discordgo.PermissionAdministrator).
		WithBotPermissions(discordgo.PermissionKickMembers | discordgo.PermissionBanMembers)
}

func antiraidHandler(ctx *discord.CommandContext) error {
	subcommand := ctx.Interaction.ApplicationCommandData().Options[0].Name

	guildData, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Interaction.GuildID})
	if err != nil {
		return ctx.ReplyEphemeral("❌ Ocurrió un error al cargar la configuración del servidor.")
	}
	if guildData == nil {
		return ctx.ReplyEphemeral("❌ No hay datos del servidor. Usa comandos básicos primero.")
	}

	antiRaid := &guildData.Protection.AntiRaid
	if antiRaid.Action == "" {
		antiRaid.Action = "kick" // Default
	}

	var response string

	switch subcommand {
	case "toggle":
		antiRaid.Enable = !antiRaid.Enable
		status := "desactivado"
		if antiRaid.Enable {
			status = "**ACTIVADO**"
		}
		response = fmt.Sprintf("🚨 Modo pánico Anti-Raid %s.", status)

	case "age":
		days := int(ctx.GetIntOption("dias"))
		antiRaid.MinAccountAgeDays = days
		if days > 0 {
			response = fmt.Sprintf("⏳ Las cuentas deberán tener al menos **%d días** de creadas para poder unirse.", days)
		} else {
			response = "⏳ El filtro por edad de cuenta ha sido desactivado."
		}

	case "limits":
		joins := int(ctx.GetIntOption("uniones"))
		seconds := int(ctx.GetIntOption("segundos"))
		if joins < 0 || seconds < 0 {
			return ctx.ReplyEphemeral("❌ Los valores deben ser positivos.")
		}
		antiRaid.JoinLimit = joins
		antiRaid.TimeWindow = seconds

		if joins == 0 || seconds == 0 {
			antiRaid.JoinLimit = 0
			antiRaid.TimeWindow = 0
			response = "📊 Detector automático de raids **desactivado**."
		} else {
			response = fmt.Sprintf("📊 Detector configurado: **%d uniones en %d segundos** activarán el modo pánico.", joins, seconds)
		}

	case "action":
		action := ctx.GetStringOption("tipo")
		antiRaid.Action = action
		actionName := "Expulsar (Kick)"
		if action == "ban" {
			actionName = "Banear (Ban)"
		}
		response = fmt.Sprintf("⚡ Ahora la acción contra raiders será: **%s**", actionName)
	}

	_, err = database.GlobalGuildDM.Set(bson.M{"id": ctx.Interaction.GuildID}, guildData)
	if err != nil {
		return ctx.ReplyEphemeral("❌ Ocurrió un error al guardar la configuración.")
	}

	return ctx.Reply(response)
}
