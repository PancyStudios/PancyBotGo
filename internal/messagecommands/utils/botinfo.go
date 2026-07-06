package utils

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/bwmarrin/discordgo"
)

var botStartTime = time.Now()

func botinfoCommand(ctx *messagecommands.MessageContext) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	numGoroutines := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()
	botVersion := config.Version
	goVersion := strings.TrimPrefix(runtime.Version(), "go")
	discordgoVersion := discordgo.VERSION

	memberCount := 0
	guildCount := len(ctx.Session.State.Guilds)
	for _, guild := range ctx.Session.State.Guilds {
		memberCount += guild.MemberCount
	}

	uptime := time.Since(botStartTime)
	ping := ctx.Session.HeartbeatLatency()

	embed := &discordgo.MessageEmbed{
		Title: "📊 Información del Bot",
		Color: 0x5865F2,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "🤖 Versión del Bot", Value: botVersion, Inline: true},
			{Name: "🐹 Versión de Go", Value: goVersion, Inline: true},
			{Name: "📚 Versión de DiscordGo", Value: discordgoVersion, Inline: true},
			{Name: "🖥 Uso de RAM", Value: fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024), Inline: true},
			{Name: "⚙ ️Uso de CPU", Value: fmt.Sprintf("%d Goroutines / %d CPUs", numGoroutines, numCPU), Inline: true},
			{Name: "⏱ Uptime", Value: formatDuration(uptime), Inline: true},
			{Name: "🏠 Guilds", Value: fmt.Sprintf("%d", guildCount), Inline: true},
			{Name: "👥 Miembros", Value: fmt.Sprintf("%d", memberCount), Inline: true},
			{Name: "🌐 Ping", Value: fmt.Sprintf("%d ms", ping.Milliseconds()), Inline: true},
			{Name: "📅 Compilación", Value: config.BuildTime, Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "💫 - Developed by PancyStudios",
			IconURL: ctx.Session.State.User.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err := ctx.ReplyEmbed(embed)
	return err
}

func formatDuration(dur time.Duration) string {
	days := int(dur.Hours() / 24)
	hours := int(dur.Hours()) % 24
	minutes := int(dur.Minutes()) % 60
	seconds := int(dur.Seconds()) % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d días", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d horas", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d minutos", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%d segundos", seconds))
	}

	return strings.Join(parts, ", ")
}
