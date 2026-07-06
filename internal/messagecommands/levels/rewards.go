package levels

import (
	"fmt"
	"strconv"
	"strings"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
)

func rewardsCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionAdministrator) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permiso de Administrador.")
		return err
	}

	guildID := ctx.Message.GuildID
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!rewards <add/remove/list> [nivel] [@rol]`")
		return err
	}

	subcommand := strings.ToLower(ctx.Args[0])
	guildData, err := database.GlobalGuildDM.Get(bson.M{"id": guildID})
	if err != nil {
		_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ Error al obtener la configuración del servidor: %v", err))
		return err
	}

	switch subcommand {
	case "add":
		if len(ctx.Args) < 3 {
			_, err = ctx.ReplyError("Uso Incorrecto", "Uso: `pan!rewards add <nivel> <@rol>`")
			return err
		}
		level, err := strconv.ParseInt(ctx.Args[1], 10, 64)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ El nivel debe ser un número válido.")
			return err
		}
		roleID := ctx.ParseRole(2)
		if roleID == "" {
			_, err = ctx.ReplyError("Error", "❌ Debes especificar un rol válido.")
			return err
		}

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

		_, err = database.GlobalGuildDM.Set(bson.M{"id": guildID}, guildData)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ Error al guardar la configuración.")
			return err
		}

		_, err = ctx.ReplySuccess("Recompensa Configurada", fmt.Sprintf("✅ Al llegar al nivel **%d**, se entregará el rol <@&%s>.", level, roleID))
		return err

	case "remove":
		if len(ctx.Args) < 2 {
			_, err = ctx.ReplyError("Uso Incorrecto", "Uso: `pan!rewards remove <nivel>`")
			return err
		}
		level, err := strconv.ParseInt(ctx.Args[1], 10, 64)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ El nivel debe ser un número válido.")
			return err
		}

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
			_, err = ctx.ReplyError("Error", fmt.Sprintf("❌ No hay ninguna recompensa configurada para el nivel %d.", level))
			return err
		}

		guildData.Levels.Rewards = newRewards
		_, err = database.GlobalGuildDM.Set(bson.M{"id": guildID}, guildData)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ Error al guardar la configuración.")
			return err
		}

		_, err = ctx.ReplySuccess("Recompensa Eliminada", fmt.Sprintf("✅ Recompensa del nivel **%d** eliminada.", level))
		return err

	case "list":
		if len(guildData.Levels.Rewards) == 0 {
			_, err = ctx.ReplySuccess("Recompensas", "📉 No hay recompensas de nivel configuradas en este servidor.")
			return err
		}

		listDesc := ""
		for _, r := range guildData.Levels.Rewards {
			listDesc += fmt.Sprintf("⭐ **Nivel %d**: <@&%s>\n", r.Level, r.RoleID)
		}

		embed := &discordgo.MessageEmbed{
			Title:       "🎁 Recompensas de Nivel",
			Description: listDesc,
			Color:       0x22d3ee,
		}

		_, err = ctx.ReplyEmbed(embed)
		return err
	}

	_, err = ctx.ReplyError("Uso Incorrecto", "Subcomando no reconocido. Usa `add`, `remove` o `list`.")
	return err
}
