package utils

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

func pingCommand(ctx *messagecommands.MessageContext) error {
	latency := ctx.Session.HeartbeatLatency().Milliseconds()
	_, err := ctx.Reply(fmt.Sprintf("🏓 Pong! Latencia: %dms", latency))
	return err
}
