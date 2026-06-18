package security

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func createAntibotsCommand() *discord.Command {
	return discord.NewCommand(
		"antibots",
		"Evita que se unan bots según su tipo",
		"security",
		antibotsHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "type",
			Description: "¿Qué tipo de bots expulso?",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Todos los bots",
					Value: "all",
				},
				{
					Name:  "Bots no verificados",
					Value: "only_nv",
				},
				{
					Name:  "Bots verificados",
					Value: "only_v",
				},
				{
					Name:  "Desactivado (Ninguno)",
					Value: "disabled",
				},
			},
		},
	).WithUserPermissions(discordgo.PermissionManageGuild).
		WithBotPermissions(discordgo.PermissionKickMembers)
}

func antibotsHandler(ctx *discord.CommandContext) error {
	option := ctx.GetStringOption("type")

	db := database.Get()
	if db == nil || !db.Connected() {
		return ctx.ReplyEphemeral("❌ La base de datos no está conectada.")
	}

	guildData, err := database.GlobalGuildDM.Get(bson.M{"_id": ctx.Interaction.GuildID})
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error al obtener los datos del servidor: %v", err))
	}

	if option == "disabled" {
		guildData.Protection.Antibots = ""
	} else {
		guildData.Protection.Antibots = option
	}

	_, err = database.GlobalGuildDM.Set(bson.M{"_id": ctx.Interaction.GuildID}, guildData)
	if err != nil {
		return ctx.ReplyEphemeral(fmt.Sprintf("❌ Error al guardar en la base de datos: %v", err))
	}

	status := "activada"
	if option == "disabled" {
		status = "desactivada"
	}

	return ctx.Reply(fmt.Sprintf("🛡️ La protección Anti-Bots ha sido **%s** (Modo: %s).", status, option))
}
