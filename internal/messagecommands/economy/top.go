package economy

import (
	"fmt"
	"sort"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"go.mongodb.org/mongo-driver/bson"
)

func topCommand(ctx *messagecommands.MessageContext, isGlobal bool) error {

	var leaderboardStr string
	var embedTitle string
	var embedColor int

	if !isGlobal {
		embedTitle = "🏆 Tabla de Clasificación Local"
		embedColor = 0x2ECC71

		profiles, err := database.LocalEconomyDM.GetAll(bson.M{"guild_id": ctx.Message.GuildID})
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ Hubo un error al obtener la tabla de clasificación local.")
			return err
		}

		sort.Slice(profiles, func(i, j int) bool {
			return (profiles[i].Wallet + profiles[i].Bank) > (profiles[j].Wallet + profiles[j].Bank)
		})

		for i, profile := range profiles {
			if i >= 10 {
				break
			}
			leaderboardStr += fmt.Sprintf("**%d.** <@%s> - 💵 %d (Cartera: %d, Banco: %d)\n", i+1, profile.UserID, profile.Wallet+profile.Bank, profile.Wallet, profile.Bank)
		}

	} else {
		embedTitle = "🏆 Tabla de Clasificación Global"
		embedColor = 0xF1C40F

		profiles, err := database.GlobalEconomyDM.GetAll(bson.M{})
		if err != nil {
			_, err = ctx.ReplyError("Error", "❌ Hubo un error al obtener la tabla de clasificación global.")
			return err
		}

		sort.Slice(profiles, func(i, j int) bool {
			return (profiles[i].StarsWallet + profiles[i].StarsBank) > (profiles[j].StarsWallet + profiles[j].StarsBank)
		})

		for i, profile := range profiles {
			if i >= 10 {
				break
			}
			leaderboardStr += fmt.Sprintf("**%d.** <@%s> - 🌟 %d (Cartera: %d, Banco: %d)\n", i+1, profile.UserID, profile.StarsWallet+profile.StarsBank, profile.StarsWallet, profile.StarsBank)
		}
	}

	if leaderboardStr == "" {
		leaderboardStr = "No hay datos para mostrar en la tabla de clasificación."
	}

	embed := discord.NewEmbed().
		SetTitle(embedTitle).
		SetColor(embedColor).
		SetDescription(leaderboardStr).
		Build()

	_, err := ctx.ReplyEmbed(embed)
	return err
}
