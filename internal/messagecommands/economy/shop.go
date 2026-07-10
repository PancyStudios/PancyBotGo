package economy

import (
	"github.com/PancyStudios/PancyBotGo/internal/commands/economy"
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/bwmarrin/discordgo"
)

func shopCommand(ctx *messagecommands.MessageContext) error {
	embed, components := economy.ShopMenu()

	_, err := ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})

	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ No se pudo cargar el catálogo de la tienda.")
		return err
	}

	return nil
}
