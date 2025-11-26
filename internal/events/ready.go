// Package events provides event handlers for the bot
package events

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// RegisterReadyEvent registers the ready event handler
func RegisterReadyEvent(client *discord.ExtendedClient) {
	client.Session.AddHandler(onReady)
}

type StatusOption struct {
	Type discordgo.ActivityType
	Text string
}

//    const activities = [
//        `PancyBot | ${version}`,
//        `pan! | ${version}`,
//        `PancyBot Studios | ${version}`
//    ]

var statusList = []StatusOption{
	{discordgo.ActivityTypeGame, "PancyBot | %s"},
	{discordgo.ActivityTypeListening, "PancyStudios | %s"},
}

// onReady solo se llama cuando está completamente conectado y listo
func onReady(s *discordgo.Session, r *discordgo.Ready) {
	logger.Success(fmt.Sprintf("✅ Bot conectado: %s#%s", r.User.Username, r.User.Discriminator), "Ready")
	logger.Info(fmt.Sprintf("📊 Conectado a %d servidores", len(r.Guilds)), "Ready")

	// Establecer estado del bot
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		rotateStatus(s)
		for {
			select {
			case <-ticker.C:
				rotateStatus(s)
			}
		}
	}()

	logger.Debug("Estado del bot establecido correctamente", "Ready")
}

func rotateStatus(s *discordgo.Session) {
	idx := rand.Intn(len(statusList))
	selected := statusList[idx]

	statusText := selected.Text
	if strings.Contains(statusText, "%d") {
		guildCount := len(s.State.Guilds)
		statusText = fmt.Sprintf(statusText, guildCount)
	} else if strings.Contains(statusText, "%s") {
		statusText = fmt.Sprintf(statusText, config.Version)
	}

	err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: statusText,
				Type: selected.Type,
			},
		},
		Status: "dnd",
		AFK:    false,
	})

	if err != nil {
		logger.Error(fmt.Sprintf("❌ Error rotando estado: %v", err), "Ready")
	}
}
