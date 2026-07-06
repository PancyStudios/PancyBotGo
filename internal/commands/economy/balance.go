package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createBalanceCommand() *discord.Command {
	return discord.NewCommand(
		"balance",
		"💸 | Revisa tu balance estelar y de monedas locales",
		"economy",
		balanceHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "💰 | Usuario para ver su balance",
			Required:    false,
		},
	)
}

func balanceHandler(ctx *discord.CommandContext) error {
	targetUser := ctx.Interaction.Member.User

	if ctx.HasOption("usuario") {
		targetUser = ctx.GetUserOption("usuario")
	}

	// Get local profile
	localProfile, err := database.GetLocalProfile(ctx.Interaction.GuildID, targetUser.ID)
	if err != nil {
		ctx.Reply("❌ " + "Hubo un error al obtener la economía local.")
		return err
	}

	// Get global profile
	globalProfile, err := database.GetGlobalProfile(targetUser.ID)
	if err != nil {
		ctx.Reply("❌ " + "Hubo un error al obtener la economía global.")
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

	ctx.ReplyEmbed(embed)
	return nil
}
