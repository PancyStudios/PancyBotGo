package utils

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
)

func statusCommand(ctx *messagecommands.MessageContext) error {
	db := database.Get()
	dbStatus, _ := db.GetStatus()
	
	guildCount := len(ctx.Session.State.Guilds)

	msg := fmt.Sprintf(
		"📊 **Estado del Bot**\n"+
			"• Bot: 🟢 Online\n"+
			"• Base de datos: %s\n"+
			"• Servidores: %d",
		dbStatus,
		guildCount,
	)
	
	_, err := ctx.ReplySuccess("Estado del Bot", msg)
	return err
}
