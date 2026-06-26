package levels

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
)

var rewardsCommand = &discord.Command{
	Name:            "rewards",
	Description:     "🎁 | Gestiona las recompensas (roles) del sistema de niveles",
	UserPermissions: discordgo.PermissionAdministrator,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "➕ | Agrega un rol como recompensa por alcanzar un nivel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "nivel",
					Description: "⭐ | Nivel requerido para obtener el rol",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "rol",
					Description: "🎭 | Rol a entregar",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "remove",
			Description: "➖ | Elimina una recompensa de un nivel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "nivel",
					Description: "⭐ | Nivel del que quieres eliminar la recompensa",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "📋 | Muestra todas las recompensas configuradas",
		},
	},
	Run: func(ctx *discord.CommandContext) error {
		guildID := ctx.Interaction.GuildID
		if guildID == "" {
			return ctx.ReplyEphemeral("❌ Este comando solo se puede usar en un servidor.")
		}

		subcommand := ctx.Interaction.ApplicationCommandData().Options[0]
		guildData, err := database.GlobalGuildDM.Get(map[string]interface{}{"_id": guildID})
		if err != nil {
			return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error al obtener la configuración del servidor: %v", err))
		}

		switch subcommand.Name {
		case "add":
			level := subcommand.Options[0].IntValue()
			roleID := subcommand.Options[1].RoleValue(nil, "").ID

			// Verificar si ya existe una recompensa para este nivel
			exists := false
			for i, r := range guildData.Levels.Rewards {
				if r.Level == level {
					guildData.Levels.Rewards[i].RoleID = roleID
					exists = true
					break
				}
			}

			if !exists {
				guildData.Levels.Rewards = append(guildData.Levels.Rewards, models.LevelReward{
					Level:  level,
					RoleID: roleID,
				})
			}

			_, err = database.GlobalGuildDM.Set(map[string]interface{}{"_id": guildID}, guildData)
			if err != nil {
				return ctx.ReplyEphemeral("❌ Error al guardar la configuración.")
			}

			return ctx.Reply(fmt.Sprintf("✅ Recompensa configurada: Al llegar al nivel **%d**, se entregará el rol <@&%s>.", level, roleID))

		case "remove":
			level := subcommand.Options[0].IntValue()

			found := false
			var newRewards []models.LevelReward
			for _, r := range guildData.Levels.Rewards {
				if r.Level == level {
					found = true
				} else {
					newRewards = append(newRewards, r)
				}
			}

			if !found {
				return ctx.ReplyEphemeral(fmt.Sprintf("❌ No hay ninguna recompensa configurada para el nivel %d.", level))
			}

			guildData.Levels.Rewards = newRewards
			_, err = database.GlobalGuildDM.Set(map[string]interface{}{"_id": guildID}, guildData)
			if err != nil {
				return ctx.ReplyEphemeral("❌ Error al guardar la configuración.")
			}

			return ctx.Reply(fmt.Sprintf("✅ Recompensa del nivel **%d** eliminada.", level))

		case "list":
			if len(guildData.Levels.Rewards) == 0 {
				return ctx.ReplyEphemeral("📉 No hay recompensas de nivel configuradas en este servidor.")
			}

			listDesc := ""
			for _, r := range guildData.Levels.Rewards {
				listDesc += fmt.Sprintf("⭐ **Nivel %d**: <@&%s>\n", r.Level, r.RoleID)
			}

			embed := discord.NewEmbed().
				SetTitle("🎁 Recompensas de Nivel").
				SetDescription(listDesc).
				SetColor(0x22d3ee)

			return ctx.ReplyEmbed(embed.Build())
		}

		return nil
	},
}
