package economy

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func createWithdrawCommand() *discord.Command {
	return discord.NewCommand(
		"withdraw",
		"🏧 | Retira monedas de tu banco",
		"economy",
		withdrawHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tipo",
			Description: "💰 | Elige si quieres retirar economía Local o Global",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{Name: "Local (Servidor)", Value: "local"},
				{Name: "Global (Estrellas)", Value: "global"},
			},
		},
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "cantidad",
			Description: "💰 | Cantidad a retirar",
			Required:    true,
		},
	)
}

func withdrawHandler(ctx *discord.CommandContext) error {
	ecoType := ctx.GetStringOption("tipo")
	amount := ctx.GetIntOption("cantidad")
	userID := ctx.Interaction.Member.User.ID
	guildID := ctx.Interaction.GuildID

	if amount <= 0 {
		ctx.Reply("❌ " + "La cantidad debe ser mayor a 0.")
		return nil
	}

	var err error
	if ecoType == "local" {
		err = database.WithdrawLocal(guildID, userID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				ctx.Reply("❌ " + "No tienes suficientes monedas en el banco local.")
			} else {
				ctx.Reply("❌ " + "Error al retirar.")
			}
			return err
		}
		ctx.Reply(fmt.Sprintf("Has retirado **💵 %d** de tu banco local.", amount))
	} else {
		err = database.WithdrawStars(userID, amount)
		if err != nil {
			if err == database.ErrInsufficientFunds {
				ctx.Reply("❌ " + "No tienes suficientes estrellas en el banco estelar.")
			} else {
				ctx.Reply("❌ " + "Error al retirar.")
			}
			return err
		}
		ctx.Reply(fmt.Sprintf("Has retirado **🌟 %d** de tu banco estelar.", amount))
	}
	return nil
}
