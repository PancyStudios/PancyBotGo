package config

import (
	"fmt"

	appconfig "github.com/PancyStudios/PancyBotGo/internal/commands/config"
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func pojCommand(ctx *messagecommands.MessageContext) error {
	if !ctx.HasPermission(discordgo.PermissionManageGuild) {
		_, err := ctx.ReplyError("Acceso Denegado", "No tienes permisos para usar este comando. Necesitas `Gestionar Servidor`.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso", "Uso: `pan!poj <add|remove|list>`")
		return err
	}

	subcommand := ctx.Args[0]
	switch subcommand {
	case "add":
		if len(ctx.Args) < 3 {
			_, err := ctx.ReplyError("Uso Incorrecto", "Uso correcto: `pan!poj add <#canal> <@&rol>`")
			return err
		}

		channelID := messagecommands.CleanMention(ctx.Args[1])
		roleID := messagecommands.CleanMention(ctx.Args[2])

		err := appconfig.AddPojConfig(ctx.Message.GuildID, channelID, roleID)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ Hubo un error al guardar.")
			return err
		}
		_, err = ctx.ReplySuccess("Éxito", fmt.Sprintf("✅ PoJ añadido: el rol <@&%s> será mencionado en <#%s>.", roleID, channelID))
		return err

	case "remove":
		if len(ctx.Args) < 2 {
			_, err := ctx.ReplyError("Uso Incorrecto", "Uso correcto: `pan!poj remove <#canal>`")
			return err
		}

		channelID := messagecommands.CleanMention(ctx.Args[1])
		err := appconfig.RemovePojConfig(ctx.Message.GuildID, channelID)
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ Hubo un error al remover.")
			return err
		}
		_, err = ctx.ReplySuccess("Éxito", fmt.Sprintf("✅ Configuración de PoJ removida del canal <#%s>.", channelID))
		return err

	case "list":
		doc, err := database.GlobalGuildDM.Get(bson.M{"id": ctx.Message.GuildID})
		if err != nil || len(doc.PingOnJoin) == 0 {
			_, err = ctx.Reply("ℹ️ No hay configuraciones de Ping On Join activas.")
			return err
		}

		list := "🔔 **Lista de Ping On Join (PoJ)**\n"
		for _, poj := range doc.PingOnJoin {
			list += fmt.Sprintf("• Canal: <#%s> | Rol: <@&%s>\n", poj.ChannelID, poj.RoleID)
		}

		_, err = ctx.Reply(list)
		return err
	default:
		_, err := ctx.ReplyError("Comando Desconocido", "❌ Subcomando desconocido. Usa `add`, `remove` o `list`.")
		return err
	}
}
