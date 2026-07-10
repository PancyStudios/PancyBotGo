package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createDepositCommand(isGlobal bool) *discord.Command {
	return discord.NewCommand(
		"deposit",
		"🏦 | Deposita tus monedas en el banco",
		"economy",
		func(ctx *discord.CommandContext) error {
			return depositHandler(ctx, isGlobal)
		},
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "cantidad",
			Description: "💰 | Cantidad a depositar",
			Required:    true,
		},
	)
}

func depositHandler(ctx *discord.CommandContext, isGlobal bool) error {

	amount := ctx.GetIntOption("cantidad")
	userID := ctx.Interaction.Member.User.ID
	guildID := ctx.Interaction.GuildID

	if amount <= 0 {
		ctx.Reply("❌ " + "La cantidad debe ser mayor a 0.")
		return nil
	}

	var err error
	if !isGlobal {
		err = database.DepositLocal(guildID, userID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				ctx.Reply("❌ " + "No tienes suficientes monedas locales en tu cartera.")
			} else if err == database.ErrBankFull {
				ctx.Reply("❌ " + "El banco local no tiene suficiente capacidad para ese depósito.")
			} else {
				ctx.Reply("❌ " + "Error al depositar.")
			}
			return err
		}
		ctx.Reply(fmt.Sprintf("Has depositado **💵 %d** a tu banco local.", amount))
	} else {
		err = database.DepositStars(userID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				ctx.Reply("❌ " + "No tienes suficientes estrellas en tu cartera.")
			} else if err == database.ErrBankFull {
				ctx.Reply("❌ " + "Tu banco estelar está al límite de su capacidad.")
			} else {
				ctx.Reply("❌ " + "Error al depositar.")
			}
			return err
		}
		ctx.Reply(fmt.Sprintf("Has depositado **🌟 %d** a tu banco estelar.", amount))
	}
	return nil
}
