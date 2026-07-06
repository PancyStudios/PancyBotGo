package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

func balanceCommand(ctx *messagecommands.MessageContext) error {
	targetUser := ctx.Message.Author
	if len(ctx.Args) > 0 {
		parsedUserID := ctx.ParseUser(0)
		if parsedUserID != "" {
			member, err := ctx.Session.GuildMember(ctx.Message.GuildID, parsedUserID)
			if err == nil {
				targetUser = member.User
			}
		}
	}

	localProfile, err := database.GetLocalProfile(ctx.Message.GuildID, targetUser.ID)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Hubo un error al obtener la economía local.")
		return err
	}

	globalProfile, err := database.GetGlobalProfile(targetUser.ID)
	if err != nil {
		_, err = ctx.ReplyError("Error", "❌ Hubo un error al obtener la economía global.")
		return err
	}

	embed := discord.NewEmbed().
		SetTitle(fmt.Sprintf("Balance de %s", targetUser.Username)).
		SetColor(discord.ColorWarning).
		SetThumbnail(targetUser.AvatarURL("")).
		SetDescription("💰 | Aquí tienes el resumen de tu economía.").
		AddField("🌐 Economía Global (Estrellas)", fmt.Sprintf("**Cartera:** 🌟 %d\n**Banco:** 🏦 %d / %d", globalProfile.StarsWallet, globalProfile.StarsBank, globalProfile.BankCapacity), false).
		AddField("🏠 Economía Local (Servidor)", fmt.Sprintf("**Cartera:** 💵 %d\n**Banco:** 🏦 %d / %d", localProfile.Wallet, localProfile.Bank, localProfile.BankCapacity), false).
		Build()

	_, err = ctx.ReplyEmbed(embed)
	return err
}
