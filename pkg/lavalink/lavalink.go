// Package lavalink provides a Lavalink client for music playback.
// It supports connecting to Lavalink nodes, searching for tracks, and controlling playback.
package lavalink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"sync"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/mqtt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

// Volume constants
const (
	MinVolume = 0
	MaxVolume = 1000
)

// NodeConfig holds configuration for a Lavalink node
type NodeConfig struct {
	Name     string
	Host     string
	Port     int
	Password string
	Secure   bool
}

// TrackInfo contains information about a track
type TrackInfo struct {
	Identifier string `json:"identifier"`
	IsSeekable bool   `json:"isSeekable"`
	Author     string `json:"author"`
	Length     int64  `json:"length"`
	IsStream   bool   `json:"isStream"`
	Position   int64  `json:"position"`
	Title      string `json:"title"`
	URI        string `json:"uri"`
	ArtworkURL string `json:"artworkUrl"`
	SourceName string `json:"sourceName"`
}

// Track represents a playable track
type Track struct {
	Encoded string    `json:"encoded"`
	Info    TrackInfo `json:"info"`
}

// Player represents a guild music player
type Player struct {
	GuildID       string
	TextChannelID string
	VoiceChannel  string
	CurrentTrack  *Track
	Queue         []*Track
	Volume        int
	IsPlaying     bool
	IsPaused      bool
	Position      int64
	Mu            sync.RWMutex // Exported for external access
}

// LavalinkClient manages the connection to Lavalink
type LavalinkClient struct {
	session         *discordgo.Session
	nodes           []*Node
	players         map[string]*Player
	mu              sync.RWMutex
	defaultPlatform string
	mqttClient      *mqtt.MqttCommunicator
	progressTickers map[string]*time.Ticker
}

// Node represents a Lavalink node connection
type Node struct {
	config       NodeConfig
	conn         *websocket.Conn
	client       *LavalinkClient
	connected    bool
	reconnecting bool
	sessionId    string
	mu           sync.RWMutex
}

// SearchResult represents a search response from Lavalink v4
type SearchResult struct {
	LoadType  string          `json:"loadType"`
	Data      json.RawMessage `json:"data"`
	Exception *struct {
		Message  string `json:"message"`
		Severity string `json:"severity"`
	} `json:"exception,omitempty"`
}

// PlaylistData represents playlist data from Lavalink
type PlaylistData struct {
	Info struct {
		Name string `json:"name"`
	} `json:"info"`
	PluginInfo interface{} `json:"pluginInfo,omitempty"`
	Tracks     []*Track    `json:"tracks"`
}

// GetTracks returns the tracks from the search result, handling different loadTypes
func (sr *SearchResult) GetTracks() []*Track {
	switch sr.LoadType {
	case "search":
		// For search results, data is an array of tracks
		var tracks []*Track
		if err := json.Unmarshal(sr.Data, &tracks); err != nil {
			logger.Error(fmt.Sprintf("Error unmarshaling search tracks: %v", err), "Lavalink")
			return []*Track{}
		}
		return tracks
	case "track":
		// For single track, data is a track object
		var track Track
		if err := json.Unmarshal(sr.Data, &track); err != nil {
			logger.Error(fmt.Sprintf("Error unmarshaling single track: %v", err), "Lavalink")
			return []*Track{}
		}
		return []*Track{&track}
	case "playlist":
		// For playlist, data is an object with tracks array
		var playlist PlaylistData
		if err := json.Unmarshal(sr.Data, &playlist); err != nil {
			logger.Error(fmt.Sprintf("Error unmarshaling playlist: %v", err), "Lavalink")
			return []*Track{}
		}
		return playlist.Tracks
	default:
		return []*Track{}
	}
}

// MusicState represents the current music state for MQTT publishing
type MusicState struct {
	GuildID      string        `json:"guildId"`
	IsPlaying    bool          `json:"isPlaying"`
	IsPaused     bool          `json:"isPaused"`
	CurrentTrack *TrackState   `json:"currentTrack"`
	Progress     float64       `json:"progress"`
	Volume       int           `json:"volume"`
	Queue        []*TrackState `json:"queue"`
	Timestamp    int64         `json:"timestamp"`
}

// TrackState represents a track in the music state
type TrackState struct {
	Title     string  `json:"title"`
	Artist    string  `json:"artist"`
	Duration  float64 `json:"duration"`
	Thumbnail string  `json:"thumbnail"`
	URL       string  `json:"url"`
}

var (
	lavalinkClient *LavalinkClient
	once           sync.Once
)

// Init initializes the global Lavalink client
func Init(session *discordgo.Session, nodeConfigs []NodeConfig) *LavalinkClient {
	once.Do(func() {
		lavalinkClient = NewLavalinkClient(session, nodeConfigs)
	})
	return lavalinkClient
}

// Get returns the global Lavalink client
func Get() *LavalinkClient {
	return lavalinkClient
}

// NewLavalinkClient creates a new Lavalink client
func NewLavalinkClient(session *discordgo.Session, nodeConfigs []NodeConfig) *LavalinkClient {
	logger.Debug("Initializing Lavalink Client", "Lavalink")

	client := &LavalinkClient{
		session:         session,
		nodes:           make([]*Node, 0),
		players:         make(map[string]*Player),
		defaultPlatform: "dzsearch",
		mqttClient:      mqtt.Get(),
		progressTickers: make(map[string]*time.Ticker),
	}

	// Initialize nodes
	for _, config := range nodeConfigs {
		node := &Node{
			config: config,
			client: client,
		}
		client.nodes = append(client.nodes, node)
	}

	// Add voice state update handler
	session.AddHandler(client.voiceStateUpdate)
	session.AddHandler(client.voiceServerUpdate)

	return client
}

// Connect connects to all Lavalink nodes
func (c *LavalinkClient) Connect() error {
	for _, node := range c.nodes {
		go node.connect()
	}
	return nil
}

// connect establishes connection to a Lavalink node
func (n *Node) connect() {
	n.mu.Lock()
	if n.connected || n.reconnecting {
		n.mu.Unlock()
		return
	}
	n.reconnecting = true
	n.mu.Unlock()

	scheme := "ws"
	if n.config.Secure {
		scheme = "wss"
	}

	url := fmt.Sprintf("%s://%s:%d/v4/websocket", scheme, n.config.Host, n.config.Port)

	headers := http.Header{}
	headers.Set("Authorization", n.config.Password)
	headers.Set("User-Id", n.client.session.State.User.ID)
	headers.Set("Client-Name", "PancyBot-Go/1.0")

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		logger.Error(fmt.Sprintf("Error al conectar con Lavalink %s: %v", n.config.Name, err), "Lavalink")
		n.mu.Lock()
		n.reconnecting = false
		n.mu.Unlock()

		// Retry connection
		go func() {
			time.Sleep(5 * time.Second)
			n.connect()
		}()
		return
	}

	n.mu.Lock()
	n.conn = conn
	n.connected = true
	n.reconnecting = false
	n.mu.Unlock()

	logger.Success(fmt.Sprintf("Conectado con Lavalink server: %s", n.config.Name), "Lavalink")

	// Start reading messages
	go n.readMessages()
}

// readMessages reads messages from the Lavalink websocket
func (n *Node) readMessages() {
	for {
		_, message, err := n.conn.ReadMessage()
		if err != nil {
			logger.Warn(fmt.Sprintf("Error leyendo mensaje de Lavalink: %v", err), "Lavalink")
			n.handleDisconnect()
			return
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(message, &payload); err != nil {
			continue
		}

		n.handleMessage(payload)
	}
}

// handleMessage processes incoming Lavalink messages
func (n *Node) handleMessage(payload map[string]interface{}) {
	op, ok := payload["op"].(string)
	if !ok {
		return
	}

	switch op {
	case "ready":
		// Extract sessionId from ready event
		if sessionId, ok := payload["sessionId"].(string); ok {
			n.mu.Lock()
			n.sessionId = sessionId
			n.mu.Unlock()
			logger.Info(fmt.Sprintf("Lavalink ready with session: %s", sessionId), "Lavalink")
		} else {
			logger.Info("Lavalink ready", "Lavalink")
		}
	case "playerUpdate":
		n.handlePlayerUpdate(payload)
	case "event":
		n.handleEvent(payload)
	case "stats":
		// Handle node statistics if needed
	}
}

// handlePlayerUpdate handles player position updates
func (n *Node) handlePlayerUpdate(payload map[string]interface{}) {
	guildID, ok := payload["guildId"].(string)
	if !ok {
		return
	}

	state, ok := payload["state"].(map[string]interface{})
	if !ok {
		return
	}

	position, _ := state["position"].(float64)

	n.client.mu.RLock()
	player, exists := n.client.players[guildID]
	n.client.mu.RUnlock()

	if exists {
		player.Mu.Lock()
		player.Position = int64(position)
		player.Mu.Unlock()
	}
}

// handleEvent handles Lavalink events
func (n *Node) handleEvent(payload map[string]interface{}) {
	eventType, ok := payload["type"].(string)
	if !ok {
		return
	}

	guildID, _ := payload["guildId"].(string)

	switch eventType {
	case "TrackStartEvent":
		n.client.handleTrackStart(guildID, payload)
	case "TrackEndEvent":
		n.client.handleTrackEnd(guildID, payload)
	case "TrackExceptionEvent":
		logger.Error(fmt.Sprintf("Track exception in guild %s", guildID), "Lavalink")
	case "TrackStuckEvent":
		logger.Warn(fmt.Sprintf("Track stuck in guild %s", guildID), "Lavalink")
	case "WebSocketClosedEvent":
		logger.Warn(fmt.Sprintf("WebSocket closed for guild %s", guildID), "Lavalink")
	}
}

// handleDisconnect handles node disconnection
func (n *Node) handleDisconnect() {
	n.mu.Lock()
	n.connected = false
	if n.conn != nil {
		n.conn.Close()
	}
	n.mu.Unlock()

	logger.Warn(fmt.Sprintf("Desconectado de Lavalink: %s. Reintentando...", n.config.Name), "Lavalink")

	time.Sleep(5 * time.Second)
	go n.connect()
}

// GetPlayer gets or creates a player for a guild
func (c *LavalinkClient) GetPlayer(guildID string) *Player {
	c.mu.Lock()
	defer c.mu.Unlock()

	if player, exists := c.players[guildID]; exists {
		return player
	}

	player := &Player{
		GuildID: guildID,
		Volume:  100,
		Queue:   make([]*Track, 0),
	}
	c.players[guildID] = player
	return player
}

// DestroyPlayer destroys a player for a guild
func (c *LavalinkClient) DestroyPlayer(guildID string) {
	c.mu.Lock()
	delete(c.players, guildID)
	c.mu.Unlock()

	c.stopProgressUpdates(guildID)

	// Send destroy command to Lavalink via REST API
	for _, node := range c.nodes {
		node.mu.RLock()
		isConnected := node.connected
		node.mu.RUnlock()

		if isConnected {
			if err := node.destroyPlayer(guildID); err != nil {
				logger.Error(fmt.Sprintf("Error destroying player: %v", err), "Lavalink")
			}
			break
		}
	}
}

// Search searches for tracks
func (c *LavalinkClient) Search(query string) (*SearchResult, error) {
	for _, node := range c.nodes {
		node.mu.RLock()
		isConnected := node.connected
		config := node.config
		nodeName := config.Name
		node.mu.RUnlock()

		if !isConnected {
			logger.Debug(fmt.Sprintf("Node %s no está conectado, saltando...", nodeName), "Lavalink")
			continue
		}

		logger.Debug(fmt.Sprintf("Usando node %s para búsqueda: %s", nodeName, query), "Lavalink")

		scheme := "http"
		if config.Secure {
			scheme = "https"
		}

		// URL encode the query
		encodedQuery := neturl.QueryEscape(fmt.Sprintf("%s:%s", c.defaultPlatform, query))
		url := fmt.Sprintf("%s://%s:%d/v4/loadtracks?identifier=%s",
			scheme, config.Host, config.Port, encodedQuery)

		logger.Debug(fmt.Sprintf("URL de búsqueda: %s", url), "Lavalink")

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.Error(fmt.Sprintf("Error creando request HTTP: %v", err), "Lavalink")
			continue
		}
		req.Header.Set("Authorization", config.Password)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error(fmt.Sprintf("Error ejecutando request HTTP a %s: %v", url, err), "Lavalink")
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logger.Error(fmt.Sprintf("Lavalink respondió con status %d", resp.StatusCode), "Lavalink")
			continue
		}

		// Read the response body
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error(fmt.Sprintf("Error leyendo body de respuesta: %v", err), "Lavalink")
			continue
		}

		var result SearchResult
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			logger.Error(fmt.Sprintf("Error decodificando respuesta JSON: %v", err), "Lavalink")
			continue
		}

		tracks := result.GetTracks()
		logger.Debug(fmt.Sprintf("LoadType: %s, Tracks encontrados: %d", result.LoadType, len(tracks)), "Lavalink")
		return &result, nil
	}

	logger.Error(fmt.Sprintf("No hay nodos Lavalink disponibles. Total de nodos: %d", len(c.nodes)), "Lavalink")
	return nil, fmt.Errorf("no available Lavalink nodes")
}

// Play starts playing a track
func (c *LavalinkClient) Play(guildID, voiceChannelID, textChannelID string, track *Track) error {
	player := c.GetPlayer(guildID)
	player.Mu.Lock()
	player.VoiceChannel = voiceChannelID
	player.TextChannelID = textChannelID
	player.Mu.Unlock()

	// Join voice channel
	if err := c.session.ChannelVoiceJoinManual(guildID, voiceChannelID, false, true); err != nil {
		return fmt.Errorf("error joining voice channel: %w", err)
	}

	// Add to queue or play
	player.Mu.Lock()
	if player.IsPlaying {
		player.Queue = append(player.Queue, track)
		player.Mu.Unlock()
		return nil
	}

	player.CurrentTrack = track
	player.IsPlaying = true
	player.Mu.Unlock()

	// Send play command via REST API
	for _, node := range c.nodes {
		node.mu.RLock()
		isConnected := node.connected
		node.mu.RUnlock()

		if isConnected {
			payload := map[string]interface{}{
				"track": map[string]interface{}{
					"encoded": track.Encoded,
				},
			}
			if err := node.updatePlayer(guildID, payload); err != nil {
				logger.Error(fmt.Sprintf("Error sending play command: %v", err), "Lavalink")
				continue
			}
			logger.Debug(fmt.Sprintf("Track enviado a reproducir: %s", track.Info.Title), "Lavalink")
			break
		}
	}

	return nil
}

// Pause pauses or resumes playback
func (c *LavalinkClient) Pause(guildID string, pause bool) error {
	player := c.GetPlayer(guildID)
	player.Mu.Lock()
	player.IsPaused = pause
	player.Mu.Unlock()

	for _, node := range c.nodes {
		node.mu.RLock()
		isConnected := node.connected
		node.mu.RUnlock()

		if isConnected {
			payload := map[string]interface{}{
				"paused": pause,
			}
			if err := node.updatePlayer(guildID, payload); err != nil {
				logger.Error(fmt.Sprintf("Error sending pause command: %v", err), "Lavalink")
				continue
			}
			return nil
		}
	}
	return fmt.Errorf("no available nodes")
}

// Stop stops playback
func (c *LavalinkClient) Stop(guildID string) error {
	player := c.GetPlayer(guildID)
	player.Mu.Lock()
	player.IsPlaying = false
	player.CurrentTrack = nil
	player.Queue = make([]*Track, 0)
	player.Mu.Unlock()

	c.stopProgressUpdates(guildID)

	for _, node := range c.nodes {
		node.mu.RLock()
		isConnected := node.connected
		node.mu.RUnlock()

		if isConnected {
			// In Lavalink v4, to stop, we set track to null
			payload := map[string]interface{}{
				"track": map[string]interface{}{
					"encoded": nil,
				},
			}
			if err := node.updatePlayer(guildID, payload); err != nil {
				logger.Error(fmt.Sprintf("Error sending stop command: %v", err), "Lavalink")
				continue
			}
			return nil
		}
	}
	return fmt.Errorf("no available nodes")
}

// Skip skips to the next track
func (c *LavalinkClient) Skip(guildID string) error {
	player := c.GetPlayer(guildID)
	player.Mu.Lock()

	if len(player.Queue) == 0 {
		player.Mu.Unlock()
		return c.Stop(guildID)
	}

	nextTrack := player.Queue[0]
	player.Queue = player.Queue[1:]
	player.CurrentTrack = nextTrack
	player.Mu.Unlock()

	for _, node := range c.nodes {
		node.mu.RLock()
		isConnected := node.connected
		node.mu.RUnlock()

		if isConnected {
			payload := map[string]interface{}{
				"track": map[string]interface{}{
					"encoded": nextTrack.Encoded,
				},
			}
			if err := node.updatePlayer(guildID, payload); err != nil {
				logger.Error(fmt.Sprintf("Error sending skip command: %v", err), "Lavalink")
				continue
			}
			return nil
		}
	}
	return fmt.Errorf("no available nodes")
}

// SetVolume sets the player volume
func (c *LavalinkClient) SetVolume(guildID string, volume int) error {
	if volume < MinVolume {
		volume = MinVolume
	}
	if volume > MaxVolume {
		volume = MaxVolume
	}

	player := c.GetPlayer(guildID)
	player.Mu.Lock()
	player.Volume = volume
	player.Mu.Unlock()

	for _, node := range c.nodes {
		node.mu.RLock()
		isConnected := node.connected
		node.mu.RUnlock()

		if isConnected {
			payload := map[string]interface{}{
				"volume": volume,
			}
			if err := node.updatePlayer(guildID, payload); err != nil {
				logger.Error(fmt.Sprintf("Error sending volume command: %v", err), "Lavalink")
				continue
			}
			return nil
		}
	}
	return fmt.Errorf("no available nodes")
}

// sendOp sends an operation to the Lavalink node (deprecated, kept for compatibility)
func (n *Node) sendOp(data map[string]interface{}) error {
	// This method is deprecated in Lavalink v4
	// Commands should use REST API instead
	return fmt.Errorf("sendOp is deprecated in Lavalink v4, use REST API")
}

// updatePlayer sends a PATCH request to update player state via REST API
func (n *Node) updatePlayer(guildID string, payload map[string]interface{}) error {
	n.mu.RLock()
	sessionId := n.sessionId
	config := n.config
	connected := n.connected
	n.mu.RUnlock()

	if !connected || sessionId == "" {
		return fmt.Errorf("node not connected or no session")
	}

	scheme := "http"
	if config.Secure {
		scheme = "https"
	}

	url := fmt.Sprintf("%s://%s:%d/v4/sessions/%s/players/%s?noReplace=false",
		scheme, config.Host, config.Port, sessionId, guildID)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", config.Password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("lavalink returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// destroyPlayer sends a DELETE request to destroy a player
func (n *Node) destroyPlayer(guildID string) error {
	n.mu.RLock()
	sessionId := n.sessionId
	config := n.config
	connected := n.connected
	n.mu.RUnlock()

	if !connected || sessionId == "" {
		return fmt.Errorf("node not connected or no session")
	}

	scheme := "http"
	if config.Secure {
		scheme = "https"
	}

	url := fmt.Sprintf("%s://%s:%d/v4/sessions/%s/players/%s",
		scheme, config.Host, config.Port, sessionId, guildID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", config.Password)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// handleTrackStart handles track start events
func (c *LavalinkClient) handleTrackStart(guildID string, payload map[string]interface{}) {
	player := c.GetPlayer(guildID)

	player.Mu.RLock()
	track := player.CurrentTrack
	player.Mu.RUnlock()

	if track == nil {
		return
	}

	logger.Info(fmt.Sprintf("Reproduciendo: %s en guild %s", track.Info.Title, guildID), "Lavalink")

	// Start progress updates
	c.startProgressUpdates(guildID)

	// Publish MQTT event
	c.publishMusicEvent(guildID, "playing", player)
}

// handleTrackEnd handles track end events
func (c *LavalinkClient) handleTrackEnd(guildID string, payload map[string]interface{}) {
	c.stopProgressUpdates(guildID)

	player := c.GetPlayer(guildID)

	// Publish MQTT stopped event
	c.publishMusicEvent(guildID, "stopped", player)

	player.Mu.Lock()
	if len(player.Queue) > 0 {
		nextTrack := player.Queue[0]
		player.Queue = player.Queue[1:]
		player.CurrentTrack = nextTrack
		player.Mu.Unlock()

		// Play next track
		for _, node := range c.nodes {
			node.mu.RLock()
			isConnected := node.connected
			node.mu.RUnlock()

			if isConnected {
				payload := map[string]interface{}{
					"track": map[string]interface{}{
						"encoded": nextTrack.Encoded,
					},
				}
				if err := node.updatePlayer(guildID, payload); err != nil {
					logger.Error(fmt.Sprintf("Error playing next track: %v", err), "Lavalink")
					continue
				}
				break
			}
		}
	} else {
		player.IsPlaying = false
		player.CurrentTrack = nil
		player.Mu.Unlock()

		logger.Info(fmt.Sprintf("Cola finalizada en guild %s", guildID), "Lavalink")
	}
}

// startProgressUpdates starts sending progress updates via MQTT
func (c *LavalinkClient) startProgressUpdates(guildID string) {
	c.stopProgressUpdates(guildID)

	ticker := time.NewTicker(5 * time.Second)
	c.mu.Lock()
	c.progressTickers[guildID] = ticker
	c.mu.Unlock()

	go func() {
		for range ticker.C {
			player := c.GetPlayer(guildID)
			player.Mu.RLock()
			isPlaying := player.IsPlaying
			player.Mu.RUnlock()

			if !isPlaying {
				c.stopProgressUpdates(guildID)
				return
			}

			c.publishMusicEvent(guildID, "progress", player)
		}
	}()
}

// stopProgressUpdates stops sending progress updates
func (c *LavalinkClient) stopProgressUpdates(guildID string) {
	c.mu.Lock()
	if ticker, exists := c.progressTickers[guildID]; exists {
		ticker.Stop()
		delete(c.progressTickers, guildID)
	}
	c.mu.Unlock()
}

// publishMusicEvent publishes a music event via MQTT
func (c *LavalinkClient) publishMusicEvent(guildID, event string, player *Player) {
	if c.mqttClient == nil {
		return
	}

	player.Mu.RLock()
	defer player.Mu.RUnlock()

	state := MusicState{
		GuildID:   guildID,
		IsPlaying: player.IsPlaying,
		IsPaused:  player.IsPaused,
		Progress:  float64(player.Position) / 1000,
		Volume:    player.Volume,
		Timestamp: time.Now().UnixMilli(),
	}

	if player.CurrentTrack != nil {
		state.CurrentTrack = &TrackState{
			Title:     player.CurrentTrack.Info.Title,
			Artist:    player.CurrentTrack.Info.Author,
			Duration:  float64(player.CurrentTrack.Info.Length) / 1000,
			Thumbnail: player.CurrentTrack.Info.ArtworkURL,
			URL:       player.CurrentTrack.Info.URI,
		}
	}

	for _, t := range player.Queue {
		state.Queue = append(state.Queue, &TrackState{
			Title:    t.Info.Title,
			Artist:   t.Info.Author,
			Duration: float64(t.Info.Length) / 1000,
		})
	}

	topic := fmt.Sprintf("pancy/music/%s/%s", guildID, event)
	c.mqttClient.Publish(topic, state)
}

// Voice handlers for Discord
func (c *LavalinkClient) voiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v.UserID != s.State.User.ID {
		return
	}

	// Store voice session ID for later use with voice server update
	// In Lavalink v4, voice updates are sent via REST API with both sessionId and server info
	logger.Debug(fmt.Sprintf("Voice state update for guild %s, session: %s", v.GuildID, v.SessionID), "Lavalink")
}

func (c *LavalinkClient) voiceServerUpdate(s *discordgo.Session, v *discordgo.VoiceServerUpdate) {
	// Get voice state from Discord
	guild, err := s.State.Guild(v.GuildID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting guild: %v", err), "Lavalink")
		return
	}

	var voiceState *discordgo.VoiceState
	for _, vs := range guild.VoiceStates {
		if vs.UserID == s.State.User.ID {
			voiceState = vs
			break
		}
	}

	if voiceState == nil {
		logger.Error("Voice state not found", "Lavalink")
		return
	}

	// Send voice update to Lavalink via REST API
	for _, node := range c.nodes {
		node.mu.RLock()
		isConnected := node.connected
		node.mu.RUnlock()

		if isConnected {
			payload := map[string]interface{}{
				"voice": map[string]interface{}{
					"sessionId": voiceState.SessionID,
					"token":     v.Token,
					"endpoint":  v.Endpoint,
				},
			}
			if err := node.updatePlayer(v.GuildID, payload); err != nil {
				logger.Error(fmt.Sprintf("Error sending voice update: %v", err), "Lavalink")
				continue
			}
			logger.Debug(fmt.Sprintf("Voice update sent for guild %s", v.GuildID), "Lavalink")
			break
		}
	}
}

// Disconnect disconnects from all nodes
func (c *LavalinkClient) Disconnect() {
	for _, node := range c.nodes {
		node.mu.Lock()
		if node.conn != nil {
			node.conn.Close()
		}
		node.connected = false
		node.mu.Unlock()
	}

	// Stop all progress tickers
	c.mu.Lock()
	for guildID, ticker := range c.progressTickers {
		ticker.Stop()
		delete(c.progressTickers, guildID)
	}
	c.mu.Unlock()

	logger.System("Lavalink client desconectado", "Lavalink")
}
