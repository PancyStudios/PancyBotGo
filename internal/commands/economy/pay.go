package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createPayCommand(isGlobal bool) *discord.Command {
	return discord.NewCommand(
		"pay",
		"💵 | Paga monedas a otro usuario",
		"economy",
		func(ctx *discord.CommandContext) error {
			return payHandler(ctx, isGlobal)
		},
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tipo",
			Description: "💰 | Elige si quieres transferir economía Local o Global",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{Name: "Local (Servidor)", Value: "local"},
				{Name: "Global (Estrellas)", Value: "global"},
			},
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuario",
			Description: "💰 | El usuario que recibirá el dinero",
			Required:    true,
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "cantidad",
			Description: "💰 | La cantidad a enviar",
			Required:    true,
		},
	)
}

func payHandler(ctx *discord.CommandContext, isGlobal bool) error {
	

	var targetUser *discordgo.User
	if ctx.HasOption("usuario") {
		targetUser = ctx.GetUserOption("usuario")
	}

	amount := ctx.GetIntOption("cantidad")
	userID := ctx.Interaction.Member.User.ID
	guildID := ctx.Interaction.GuildID

	if targetUser == nil || targetUser.ID == userID {
		ctx.Reply("❌ " + "No puedes transferirte dinero a ti mismo.")
		return nil
	}
	if targetUser.Bot {
		ctx.Reply("❌ " + "Los bots no tienen economía.")
		return nil
	}
	if amount <= 0 {
		ctx.Reply("❌ " + "La cantidad debe ser mayor a 0.")
		return nil
	}

	var err error
	if !isGlobal {
		err = database.TransferLocalBalance(guildID, userID, targetUser.ID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				ctx.Reply("❌ " + "No tienes suficientes monedas locales en tu cartera.")
			} else {
				ctx.Reply("❌ " + "Error al procesar la transferencia local.")
			}
			return err
		}
		ctx.Reply(fmt.Sprintf("Has transferido **💵 %d** monedas locales a %s.", amount, targetUser.Mention()))
	} else {
		err = database.TransferStars(userID, targetUser.ID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				ctx.Reply("❌ " + "No tienes suficientes estrellas en tu cartera.")
			} else {
				ctx.Reply("❌ " + "Error al procesar la transferencia estelar.")
			}
			return err
		}
		ctx.Reply(fmt.Sprintf("Has transferido **🌟 %d** estrellas a %s.", amount, targetUser.Mention()))
	}
	return nil
}
