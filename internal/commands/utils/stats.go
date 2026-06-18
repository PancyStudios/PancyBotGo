// filepath: /home/turbis/GolandProjects/PancyBotGo/internal/commands/utils/stats.go
package utils

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/bwmarrin/discordgo"
)

// createBotinfoCommand creates the /utils botinfo subcommand
func createBotinfoCommand() *discord.Command {
	return discord.NewCommand(
		"botinfo",
		"Muestra estadísticas detalladas del bot",
		"utils",
		botinfoHandler,
	)
}

// botinfoHandler handles the /utils botinfo command
func botinfoHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()

		// Get memory stats
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		// Get CPU stats (simplified)
		numGoroutines := runtime.NumGoroutine()
		numCPU := runtime.NumCPU()

		// Get bot version (hardcoded for now)
		botVersion := config.Version

		// Get Go version
		goVersion := strings.TrimPrefix(runtime.Version(), "go")

		// Get discordgo version
		discordgoVersion := discordgo.VERSION

		// Get guild and member count
		guildCount := ctx.Client.GuildCount()
		memberCount := 0
		for _, guild := range ctx.Session.State.Guilds {
			memberCount += guild.MemberCount
		}

		// Calculate uptime
		uptime := time.Since(ctx.Client.StartTime)

		// Calculate Ping
		ping := ctx.Session.HeartbeatLatency()

		// Create embed
		embed := &discordgo.MessageEmbed{
			Title: "📊 Información del Bot",
			Color: 0x5865F2,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "🤖 Versión del Bot",
					Value:  botVersion,
					Inline: true,
				},
				{
					Name:   "🐹 Versión de Go",
					Value:  goVersion,
					Inline: true,
				},
				{
					Name:   "📚 Versión de DiscordGo",
					Value:  discordgoVersion,
					Inline: true,
				},
				{
					Name:   "🖥 Uso de RAM",
					Value:  fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
					Inline: true,
				},
				{
					Name:   "⚙ ️Uso de CPU",
					Value:  fmt.Sprintf("%d Goroutines / %d CPUs", numGoroutines, numCPU),
					Inline: true,
				},
				{
					Name:   "⏱ Uptime",
					Value:  formatDuration(uptime),
					Inline: true,
				},
				{
					Name:   "🏠 Guilds",
					Value:  fmt.Sprintf("%d", guildCount),
					Inline: true,
				},
				{
					Name:   "👥 Miembros",
					Value:  fmt.Sprintf("%d", memberCount),
					Inline: true,
				},
				{
					Name:   "🌐 Ping",
					Value:  fmt.Sprintf("%d ms", ping.Milliseconds()),
					Inline: true,
				},
				{
					Name:   "📅 Compilación",
					Value:  config.BuildTime,
					Inline: true,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "💫 - Developed by PancyStudios",
				IconURL: ctx.Client.Session.State.User.AvatarURL(""),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		ctx.ReplyEmbed(embed)
	}()
	return nil
}

// formatDuration formats a time.Duration into a human-readable string
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
