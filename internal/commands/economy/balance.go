package economy

import (
	"fmt"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createBalanceCommand() *discord.Command {
	return discord.NewCommand(
		"balance",
		"💰 Revisa tu cartera y banco (Estrellas globales y monedas locales)",
		"economy",
		balanceHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "Usuario para ver su balance",
			Required:    false,
		},
	)
}

func balanceHandler(ctx *discord.CommandContext) error {
	targetUser := ctx.Interaction.Member.User
	
	// CommandContext doesn't have an explicit option for user object, but we can extract it
	if len(ctx.Interaction.ApplicationCommandData().Options) > 0 {
		targetUser = ctx.Interaction.ApplicationCommandData().Options[0].UserValue(ctx.Session)
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

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Balance de %s", targetUser.Username),
		Color:       0xF1C40F,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: targetUser.AvatarURL("")},
		Description: "Aquí tienes el resumen de tu economía.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "🌐 Economía Global (Estrellas)",
				Value:  fmt.Sprintf("**Cartera:** 🌟 %d\n**Banco:** 🏦 %d / %d", globalProfile.StarsWallet, globalProfile.StarsBank, globalProfile.BankCapacity),
				Inline: false,
			},
			{
				Name:   "🏠 Economía Local (Servidor)",
				Value:  fmt.Sprintf("**Cartera:** 💵 %d\n**Banco:** 🏦 %d / %d", localProfile.Wallet, localProfile.Bank, localProfile.BankCapacity),
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	ctx.ReplyEmbed(embed)
	return nil
}
