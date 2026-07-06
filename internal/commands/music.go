// Package commands provides music commands for the bot.
package commands

import (
	"fmt"
	"strings"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/lavalink"
	"github.com/bwmarrin/discordgo"
)

// minVolumeFloat is the minimum volume value for Discord command options
var minVolumeFloat = 0.0

// RegisterMusicCommands registers all music commands
func RegisterMusicCommands(client *discord.ExtendedClient) {
	// Play command
	playCmd := discord.NewCommand(
		"play",
		"✨ | Reproduce una canción o la añade a la cola",
		"music",
		playHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "query",
			Description:  "Nombre de la canción o URL",
			Required:     true,
			Autocomplete: true,
		},
	).WithAutoComplete(playAutoComplete).RequiresVoice()
	client.CommandHandler.RegisterCommand(playCmd)

	// Pause command
	pauseCmd := discord.NewCommand(
		"pause",
		"✨ | Pausa o resume la reproducción",
		"music",
		pauseHandler,
	).RequiresVoice()
	client.CommandHandler.RegisterCommand(pauseCmd)

	// Skip command
	skipCmd := discord.NewCommand(
		"skip",
		"✨ | Salta a la siguiente canción",
		"music",
		skipHandler,
	).RequiresVoice()
	client.CommandHandler.RegisterCommand(skipCmd)

	// Stop command
	stopCmd := discord.NewCommand(
		"stop",
		"✨ | Detiene la reproducción y limpia la cola",
		"music",
		stopHandler,
	).RequiresVoice()
	client.CommandHandler.RegisterCommand(stopCmd)

	// Queue command
	queueCmd := discord.NewCommand(
		"queue",
		"✨ | Muestra la cola de reproducción",
		"music",
		queueHandler,
	)
	client.CommandHandler.RegisterCommand(queueCmd)

	// Volume command
	volumeCmd := discord.NewCommand(
		"volume",
		"✨ | Ajusta el volumen de reproducción",
		"music",
		volumeHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "level",
			Description: "Nivel de volumen (0-100)",
			Required:    true,
			MinValue:    &minVolumeFloat,
			MaxValue:    100,
		},
	).RequiresVoice()
	client.CommandHandler.RegisterCommand(volumeCmd)

	// NowPlaying command
	npCmd := discord.NewCommand(
		"nowplaying",
		"✨ | Muestra la canción que se está reproduciendo",
		"music",
		nowPlayingHandler,
	)
	client.CommandHandler.RegisterCommand(npCmd)

	// Radio command
	var radioChoices []*discordgo.ApplicationCommandOptionChoice
	for name, url := range radioStations {
		radioChoices = append(radioChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: url,
		})
	}

	radioCmd := discord.NewCommand(
		"radio",
		"📻 | Sintoniza una estación de radio 24/7",
		"music",
		radioHandler,
	).WithOptions(
		&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "estacion",
			Description: "La estación de radio a sintonizar",
			Required:    true,
			Choices:     radioChoices,
		},
	).RequiresVoice()
	client.CommandHandler.RegisterCommand(radioCmd)
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
func playHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		query := ctx.GetStringOption("query")
		if query == "" {
			err := ctx.ReplyEphemeral("❌ Debes proporcionar una canción para reproducir.")
			if err != nil {
				return
			}
			return
		}

		// Get user's voice channel
		voiceState, err := ctx.Session.State.VoiceState(ctx.Interaction.GuildID, ctx.User().ID)
		if err != nil || voiceState.ChannelID == "" {
			err := ctx.ReplyEphemeral("❌ Debes estar en un canal de voz.")
			if err != nil {
				return
			}
			return
		}

		// Defer the response since search might take time
		ctx.Defer()

		// Search for the track
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			err := ctx.EditReply("❌ El sistema de música no está disponible.")
			if err != nil {
				return
			}
			return
		}

		result, err := lavalinkClient.Search(query)
		if err != nil {
			err := ctx.EditReply(fmt.Sprintf("❌ Error buscando: %v", err))
			if err != nil {
				return
			}
			return
		}

		tracks := result.GetTracks()
		if result.LoadType == "empty" || len(tracks) == 0 {
			err := ctx.EditReply("❌ No se encontraron resultados.")
			if err != nil {
				return
			}
			return
		}

		if result.LoadType == "error" && result.Exception != nil {
			err := ctx.EditReply(fmt.Sprintf("❌ Error: %s", result.Exception.Message))
			if err != nil {
				return
			}
			return
		}

		// Handle playlists
		if result.LoadType == "playlist" {
			// Add all tracks to queue
			for _, track := range tracks {
				track.RequesterID = ctx.User().ID
				track.RequesterName = ctx.User().Username
				if err := lavalinkClient.Play(ctx.Interaction.GuildID, voiceState.ChannelID, ctx.Interaction.ChannelID, track); err != nil {
					// Log error but continue
					continue
				}
			}
			embed := discord.NewEmbed().
				SetColor(0x5865F2).
				SetTitle("🎵 Playlist añadida a la cola").
				SetDescription(fmt.Sprintf("**Playlist** - %d canciones", len(tracks))).
				Build()
			err = ctx.EditReplyEmbed(embed)
			if err != nil {
				return
			}
			return
		}

		// Single track
		track := tracks[0]
		track.RequesterID = ctx.User().ID
		track.RequesterName = ctx.User().Username

		// Play the track
		if err := lavalinkClient.Play(ctx.Interaction.GuildID, voiceState.ChannelID, ctx.Interaction.ChannelID, track); err != nil {
			err := ctx.EditReply(fmt.Sprintf("❌ Error reproduciendo: %v", err))
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

		err = ctx.EditReplyEmbed(embed)
		if err != nil {
			return
		}
		return
	}()
	return nil
}

// pauseHandler handles the /pause command
func pauseHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			err := ctx.ReplyEphemeral("❌ El sistema de música no está disponible.")
			if err != nil {
				return
			}
			return
		}

		player := lavalinkClient.GetPlayer(ctx.Interaction.GuildID)
		player.Mu.RLock()
		isPaused := player.IsPaused
		player.Mu.RUnlock()

		if err := lavalinkClient.Pause(ctx.Interaction.GuildID, !isPaused); err != nil {
			err := ctx.ReplyEphemeral(fmt.Sprintf("❌ Error: %v", err))
			if err != nil {
				return
			}
			return
		}

		if isPaused {
			err := ctx.Reply("▶️ Reproducción resumida.")
			if err != nil {
				return
			}
			return
		}
		err := ctx.Reply("⏸️ Reproducción pausada.")
		if err != nil {
			return
		}
		return
	}()
	return nil
}

// skipHandler handles the /skip command
func skipHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			err := ctx.ReplyEphemeral("❌ El sistema de música no está disponible.")
			if err != nil {
				return
			}
			return
		}

		if err := lavalinkClient.Skip(ctx.Interaction.GuildID); err != nil {
			err := ctx.ReplyEphemeral(fmt.Sprintf("❌ Error: %v", err))
			if err != nil {
				return
			}
			return
		}

		err := ctx.Reply("⏭️ Canción saltada.")
		if err != nil {
			return
		}
		return
	}()
	return nil
}

// stopHandler handles the /stop command
func stopHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			ctx.ReplyEphemeral("❌ El sistema de música no está disponible.")
			return
		}

		if err := lavalinkClient.Stop(ctx.Interaction.GuildID); err != nil {
			ctx.ReplyEphemeral(fmt.Sprintf("❌ Error: %v", err))
			return
		}

		lavalinkClient.DestroyPlayer(ctx.Interaction.GuildID)

		ctx.Reply("⏹️ Reproducción detenida y cola limpiada.")
		return
	}()
	return nil
}

// queueHandler handles the /queue command
func queueHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			ctx.ReplyEphemeral("❌ El sistema de música no está disponible.")
			return
		}

		player := lavalinkClient.GetPlayer(ctx.Interaction.GuildID)
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
func volumeHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		level := int(ctx.GetIntOption("level"))

		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			ctx.ReplyEphemeral("❌ El sistema de música no está disponible.")
			return
		}

		if err := lavalinkClient.SetVolume(ctx.Interaction.GuildID, level); err != nil {
			ctx.ReplyEphemeral(fmt.Sprintf("❌ Error: %v", err))
			return
		}

		ctx.Reply(fmt.Sprintf("🔊 Volumen ajustado a %d%%", level))
		return
	}()
	return nil
}

// nowPlayingHandler handles the /nowplaying command
func nowPlayingHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			ctx.ReplyEphemeral("❌ El sistema de música no está disponible.")
			return
		}

		player := lavalinkClient.GetPlayer(ctx.Interaction.GuildID)
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
func playAutoComplete(ctx *discord.CommandContext) {
	go func() {
		defer errors.RecoverMiddleware()()
		query := ctx.GetStringOption("query")

		if query == "" {
			return
		}

		// If it looks like a URL, don't provide autocomplete
		if strings.HasPrefix(query, "http://") || strings.HasPrefix(query, "https://") {
			return
		}

		// Search for tracks
		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			return
		}

		result, err := lavalinkClient.Search(query)
		if err != nil {
			return
		}

		tracks := result.GetTracks()
		if len(tracks) == 0 {
			return
		}

		choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, 10)
		for i, track := range tracks {
			if i >= 10 {
				break
			}
			name := fmt.Sprintf("🎧 | %s - %s", track.Info.Author, track.Info.Title)
			if len(name) > 100 {
				name = name[:97] + "..."
			}
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  name,
				Value: track.Info.URI,
			})
		}

		ctx.SendAutoCompleteChoices(choices)
	}()
	return
}

// formatDuration formats milliseconds to mm:ss format
func formatDuration(ms int64) string {
	seconds := ms / 1000
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// radioHandler handles the /radio command
func radioHandler(ctx *discord.CommandContext) error {
	go func() {
		defer errors.RecoverMiddleware()()
		stationURL := ctx.GetStringOption("estacion")

		// Get user's voice channel
		voiceState, err := ctx.Session.State.VoiceState(ctx.Interaction.GuildID, ctx.User().ID)
		if err != nil || voiceState.ChannelID == "" {
			err := ctx.ReplyEphemeral("❌ Debes estar en un canal de voz.")
			if err != nil {
				return
			}
			return
		}

		ctx.Defer()

		lavalinkClient := lavalink.Get()
		if lavalinkClient == nil {
			err := ctx.EditReply("❌ El sistema de música no está disponible.")
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
			err := ctx.EditReply(fmt.Sprintf("❌ Error buscando la radio: %v", err))
			if err != nil {
				return
			}
			return
		}

		tracks := result.GetTracks()
		if result.LoadType == "empty" || len(tracks) == 0 {
			err := ctx.EditReply("❌ No se pudo conectar a la estación de radio.")
			if err != nil {
				return
			}
			return
		}

		track := tracks[0]

		// Stop current playback to override with radio
		if err := lavalinkClient.Stop(ctx.Interaction.GuildID); err != nil {
			// ignore error if nothing was playing
		}

		// Play the radio track
		if err := lavalinkClient.Play(ctx.Interaction.GuildID, voiceState.ChannelID, ctx.Interaction.ChannelID, track); err != nil {
			err := ctx.EditReply(fmt.Sprintf("❌ Error reproduciendo: %v", err))
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

		err = ctx.EditReplyEmbed(embed)
		if err != nil {
			return
		}
	}()
	return nil
}
