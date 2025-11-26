package events

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

func RegisterShardEvents(client *discord.ExtendedClient) {
	client.Session.AddHandler(onShardDisconnect)
	client.Session.AddHandler(onShardResumed)
}

func onShardDisconnect(s *discordgo.Session, event *discordgo.Disconnect) {
	var shardID = s.ShardID
	logger.Info(fmt.Sprintf("🔌 Shard %d desconectado.", shardID), "Shard")
}

func onShardResumed(s *discordgo.Session, event *discordgo.Resumed) {
	var shardID = s.ShardID
	logger.Success(fmt.Sprintf("✅ Shard %d reanudado.", shardID), "Shard")
}
