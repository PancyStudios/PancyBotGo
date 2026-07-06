package dev

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"github.com/bwmarrin/discordgo"
)

func blacklistCommand(ctx *messagecommands.MessageContext) error {
	if !isDev(ctx.Message.Author.ID) {
		_, err := ctx.ReplyError("Acceso Denegado", "Este comando es solo para la desarrolladora.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!blacklist <add/remove/list> [user/guild] [id] [razon]`")
		return err
	}

	action := strings.ToLower(ctx.Args[0])

	switch action {
	case "add":
		if len(ctx.Args) < 3 {
			_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!blacklist add <user/guild> <id> [razon]`")
			return err
		}
		blType := strings.ToLower(ctx.Args[1])
		id := ctx.Args[2]
		reason := "Sin razón especificada."
		if len(ctx.Args) > 3 {
			reason = strings.Join(ctx.Args[3:], " ")
		}

		if blType != "user" && blType != "guild" {
			_, err := ctx.ReplyError("Error", "El tipo debe ser `user` o `guild`.")
			return err
		}

		var dbBlType models.BlacklistType
		if blType == "user" {
			dbBlType = models.BlacklistTypeUser
		} else {
			dbBlType = models.BlacklistTypeGuild
		}

		_, err := database.GetBlacklistEntry(id)
		if err == nil {
			_, err = ctx.ReplyError("Error", "Este ID ya está en la blacklist.")
			return err
		}

		_, err = database.AddToBlacklist(id, dbBlType, reason, ctx.Message.Author.ID)
		if err != nil {
			_, err = ctx.ReplyError("Error", "No se pudo añadir a la blacklist.")
			return err
		}

		_, err = ctx.ReplySuccess("Blacklist Añadida", fmt.Sprintf("✅ `%s` (%s) fue añadido a la blacklist por: %s", id, blType, reason))
		return err

	case "remove":
		if len(ctx.Args) < 2 {
			_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!blacklist remove <id>`")
			return err
		}
		id := ctx.Args[1]

		err := database.RemoveFromBlacklist(id)
		if err != nil {
			_, err = ctx.ReplyError("Error", "No se pudo eliminar de la blacklist (puede que no exista).")
			return err
		}

		_, err = ctx.ReplySuccess("Blacklist Eliminada", fmt.Sprintf("✅ `%s` fue eliminado de la blacklist.", id))
		return err

	case "list":
		entries := database.GetBlacklistCache().GetAll()

		if len(entries) == 0 {
			_, err := ctx.ReplySuccess("Blacklist", "La blacklist está vacía.")
			return err
		}

		desc := ""
		for _, e := range entries {
			desc += fmt.Sprintf("• **%s** (%s) - Razón: %s\n", e.ID, string(e.Type), e.Reason)
		}

		embed := &discordgo.MessageEmbed{
			Title:       "📋 Blacklist Global",
			Description: desc,
			Color:       0xff0000,
		}

		_, err := ctx.ReplyEmbed(embed)
		return err
	}

	_, err := ctx.ReplyError("Uso Incorrecto", "Acción no reconocida.")
	return err
}
