// Package commands provides music commands for the bot.
package music

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/lavalink"
)

// minVolumeFloat is the minimum volume value for Discord command options
var minVolumeFloat = 0.0


func RegisterAll() {
	messagecommands.RegisterCommand("play", playHandler)
	messagecommands.RegisterCommand("pause", pauseHandler)
	messagecommands.RegisterCommand("skip", skipHandler)
	messagecommands.RegisterCommand("stop", stopHandler)
	messagecommands.RegisterCommand("queue", queueHandler)
	messagecommands.RegisterCommand("volume", volumeHandler)
	messagecommands.RegisterCommand("nowplaying", nowPlayingHandler)
	messagecommands.RegisterCommand("radio", radioHandler)
	
	
}

// Predefined radio stations (Direct HTTP Streams to bypass YouTube)
var radioStations = map[string]string{
	"Lofi Hip Hop Radio":        "https://lfhh.radioca.st/stream",
	"Nightride FM (Synthwave)":  "https://stream.nightride.fm/nightride.m4a",
	"Vaporwave Plaza":           "https://radio.plaza.one/mp3",
	"I Love Radio (Pop/Hits)":   "https://streams.ilovemusic.de/iloveradio1.mp3",
	"I Love Dance (Electronic)": "https://streams.ilovemusic.de/iloveradio2.mp3",
}

// playHandler handles the /play command
func playHandler(ctx *messagecommands.MessageContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		if len(ctx.Args) == 0 { ctx.ReplyError("Error", "Debes especificar una canción."); return }
	query := strings.Join(ctx.Args, " ")
		if query == "" {
			_, err := ctx.Reply("❌ Debes proporcionar una canción para reproducir.")
			if err != nil {
				return
			}
			return
		}

		// Get user's voice channel
		voiceState, err := ctx.Session.State.VoiceState(ctx.Message.GuildID, ctx.Message.Author.ID)
		if err != nil || voiceState.ChannelID == "" {
			_, err := ctx.Reply("❌ Debes estar en un canal de voz.")
			if err != nil {
				return
			}
			return
		}

		// Defer the response since search might take time
		

		// Search for the track
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			_, err := ctx.Reply("❌ El sistema de música no está disponible.")
			if err != nil {
				return
			}
			return
		}

		result, err := lavalinkClient.Search(query)
		if err != nil {
			_, err := ctx.Reply(fmt.Sprintf("❌ Error buscando: %v", err))
			if err != nil {
				return
			}
			return
		}

		tracks := result.GetTracks()
		if result.LoadType == "empty" || len(tracks) == 0 {
			_, err := ctx.Reply("❌ No se encontraron resultados.")
			if err != nil {
				return
			}
			return
		}

		if result.LoadType == "error" && result.Exception != nil {
			_, err := ctx.Reply(fmt.Sprintf("❌ Error: %s", result.Exception.Message))
			if err != nil {
				return
			}
			return
		}

		// Handle playlists
		if result.LoadType == "playlist" {
			// Add all tracks to queue
			for _, track := range tracks {
				track.RequesterID = ctx.Message.Author.ID
				track.RequesterName = ctx.Message.Author.Username
				if err := lavalinkClient.Play(ctx.Message.GuildID, voiceState.ChannelID, ctx.Message.ChannelID, track); err != nil {
					// Log error but continue
					continue
				}
			}
			embed := discord.NewEmbed().
				SetColor(0x5865F2).
				SetTitle("🎵 Playlist añadida a la cola").
				SetDescription(fmt.Sprintf("**Playlist** - %d canciones", len(tracks))).
				Build()
			_, err = ctx.ReplyEmbed(embed)
			if err != nil {
				return
			}
			return
		}

		// Single track
		track := tracks[0]
		track.RequesterID = ctx.Message.Author.ID
		track.RequesterName = ctx.Message.Author.Username

		// Play the track
		if err := lavalinkClient.Play(ctx.Message.GuildID, voiceState.ChannelID, ctx.Message.ChannelID, track); err != nil {
			_, err := ctx.Reply(fmt.Sprintf("❌ Error reproduciendo: %v", err))
			if err != nil {
				return
			}
			return
		}

		embed := discord.NewEmbed().
			SetColor(0x5865F2).
			SetTitle("🎵 Añadido a la cola").
			SetDescription(fmt.Sprintf("[%s](%s)", track.Info.Title, track.Info.URI)).
			SetThumbnail(track.Info.ArtworkURL).
			AddField("Artista", track.Info.Author, true).
			AddField("Duración", formatDuration(track.Info.Length), true).
			Build()

		_, err = ctx.ReplyEmbed(embed)
		if err != nil {
			return
		}
		return
	}()
	return nil
}

// pauseHandler handles the /pause command
func pauseHandler(ctx *messagecommands.MessageContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			_, err := ctx.Reply("❌ El sistema de música no está disponible.")
			if err != nil {
				return
			}
			return
		}

		player := lavalinkClient.GetPlayer(ctx.Message.GuildID)
		player.Mu.RLock()
		isPaused := player.IsPaused
		player.Mu.RUnlock()

		if err := lavalinkClient.Pause(ctx.Message.GuildID, !isPaused); err != nil {
			_, err := ctx.Reply(fmt.Sprintf("❌ Error: %v", err))
			if err != nil {
				return
			}
			return
		}

		if isPaused {
			_, err := ctx.Reply("▶️ Reproducción resumida.")
			if err != nil {
				return
			}
			return
		}
		_, err := ctx.Reply("⏸️ Reproducción pausada.")
		if err != nil {
			return
		}
		return
	}()
	return nil
}

// skipHandler handles the /skip command
func skipHandler(ctx *messagecommands.MessageContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			_, err := ctx.Reply("❌ El sistema de música no está disponible.")
			if err != nil {
				return
			}
			return
		}

		if err := lavalinkClient.Skip(ctx.Message.GuildID); err != nil {
			_, err := ctx.Reply(fmt.Sprintf("❌ Error: %v", err))
			if err != nil {
				return
			}
			return
		}

		_, err := ctx.Reply("⏭️ Canción saltada.")
		if err != nil {
			return
		}
		return
	}()
	return nil
}

// stopHandler handles the /stop command
func stopHandler(ctx *messagecommands.MessageContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			ctx.Reply("❌ El sistema de música no está disponible.")
			return
		}

		if err := lavalinkClient.Stop(ctx.Message.GuildID); err != nil {
			ctx.Reply(fmt.Sprintf("❌ Error: %v", err))
			return
		}

		lavalinkClient.DestroyPlayer(ctx.Message.GuildID)

		ctx.Reply("⏹️ Reproducción detenida y cola limpiada.")
		return
	}()
	return nil
}

// queueHandler handles the /queue command
func queueHandler(ctx *messagecommands.MessageContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			ctx.Reply("❌ El sistema de música no está disponible.")
			return
		}

		player := lavalinkClient.GetPlayer(ctx.Message.GuildID)
		player.Mu.RLock()
		defer player.Mu.RUnlock()

		if player.CurrentTrack == nil && len(player.Queue) == 0 {
			ctx.Reply("📭 La cola está vacía.")
			return
		}

		var sb strings.Builder
		sb.WriteString("📋 **Cola de reproducción**\n\n")

		if player.CurrentTrack != nil {
			sb.WriteString(fmt.Sprintf("🎵 **Reproduciendo:** [%s](%s) - %s\n\n",
				player.CurrentTrack.Info.Title,
				player.CurrentTrack.Info.URI,
				formatDuration(player.CurrentTrack.Info.Length)))
		}

		if len(player.Queue) > 0 {
			sb.WriteString("**Siguiente:**\n")
			for i, track := range player.Queue {
				if i >= 10 {
					sb.WriteString(fmt.Sprintf("\n... y %d más", len(player.Queue)-10))
					break
				}
				sb.WriteString(fmt.Sprintf("%d. %s - %s\n",
					i+1, track.Info.Title, formatDuration(track.Info.Length)))
			}
		}

		ctx.Reply(sb.String())
		return
	}()
	return nil
}

// volumeHandler handles the /volume command
func volumeHandler(ctx *messagecommands.MessageContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		if len(ctx.Args) == 0 { ctx.ReplyError("Error", "Debes especificar un volumen."); return }
		level, err := strconv.Atoi(ctx.Args[0])
		if err != nil {
			ctx.ReplyError("Error", "Volumen inválido.")
			return
		}

		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			ctx.Reply("❌ El sistema de música no está disponible.")
			return
		}

		if err := lavalinkClient.SetVolume(ctx.Message.GuildID, level); err != nil {
			ctx.Reply(fmt.Sprintf("❌ Error: %v", err))
			return
		}

		ctx.Reply(fmt.Sprintf("🔊 Volumen ajustado a %d%%", level))
		return
	}()
	return nil
}

// nowPlayingHandler handles the /nowplaying command
func nowPlayingHandler(ctx *messagecommands.MessageContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			ctx.Reply("❌ El sistema de música no está disponible.")
			return
		}

		player := lavalinkClient.GetPlayer(ctx.Message.GuildID)
		player.Mu.RLock()
		defer player.Mu.RUnlock()

		if player.CurrentTrack == nil {
			ctx.Reply("🔇 No hay nada reproduciéndose.")
			return
		}

		track := player.CurrentTrack
		progress := float64(player.Position) / float64(track.Info.Length) * 100

		embed := discord.NewEmbed().
			SetColor(0x5865F2).
			SetTitle("🎵 Reproduciendo ahora").
			SetDescription(fmt.Sprintf("[%s](%s)", track.Info.Title, track.Info.URI)).
			SetThumbnail(track.Info.ArtworkURL).
			AddField("Artista", track.Info.Author, true).
			AddField("Progreso", fmt.Sprintf("%s / %s (%.1f%%)", formatDuration(player.Position), formatDuration(track.Info.Length), progress), true).
			AddField("Volumen", fmt.Sprintf("%d%%", player.Volume), true).
			Build()
		ctx.ReplyEmbed(embed)
		return
	}()
	return nil
}

// playAutoComplete handles autocomplete for the play command

// formatDuration formats milliseconds to mm:ss format
func formatDuration(ms int64) string {
	seconds := ms / 1000
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// radioHandler handles the /radio command
func radioHandler(ctx *messagecommands.MessageContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		if len(ctx.Args) == 0 { ctx.ReplyError("Error", "Debes especificar una estación."); return }
		stationURL := ctx.Args[0]

		// Get user's voice channel
		voiceState, err := ctx.Session.State.VoiceState(ctx.Message.GuildID, ctx.Message.Author.ID)
		if err != nil || voiceState.ChannelID == "" {
			_, err := ctx.Reply("❌ Debes estar en un canal de voz.")
			if err != nil {
				return
			}
			return
		}

		

		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			_, err := ctx.Reply("❌ El sistema de música no está disponible.")
			if err != nil {
				return
			}
			return
		}

		// Find the station name
		var stationName string
		for name, url := range radioStations {
			if url == stationURL {
				stationName = name
				break
			}
		}

		result, err := lavalinkClient.Search(stationURL)
		if err != nil {
			_, err := ctx.Reply(fmt.Sprintf("❌ Error buscando la radio: %v", err))
			if err != nil {
				return
			}
			return
		}

		tracks := result.GetTracks()
		if result.LoadType == "empty" || len(tracks) == 0 {
			_, err := ctx.Reply("❌ No se pudo conectar a la estación de radio.")
			if err != nil {
				return
			}
			return
		}

		track := tracks[0]

		// Stop current playback to override with radio
		if err := lavalinkClient.Stop(ctx.Message.GuildID); err != nil {
			// ignore error if nothing was playing
		}

		// Play the radio track
		if err := lavalinkClient.Play(ctx.Message.GuildID, voiceState.ChannelID, ctx.Message.ChannelID, track); err != nil {
			_, err := ctx.Reply(fmt.Sprintf("❌ Error reproduciendo: %v", err))
			if err != nil {
				return
			}
			return
		}

		embed := discord.NewEmbed().
			SetColor(0xF1C40F).
			SetTitle("📻 Radio Sintonizada").
			SetDescription(fmt.Sprintf("Reproduciendo **%s** 24/7\n\n[%s](%s)", stationName, track.Info.Title, track.Info.URI)).
			SetThumbnail(track.Info.ArtworkURL).
			Build()

		_, err = ctx.ReplyEmbed(embed)
		if err != nil {
			return
		}
	}()
	return nil
}
