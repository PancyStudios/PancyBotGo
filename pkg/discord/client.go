// Package discord provides the Discord bot client and related structures.
// It wraps discordgo with additional functionality for command and event handling.
package discord

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord/premium"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// DiscordGoLogger wraps the custom logger to implement discordgo.Logger interface
// Note: discordgo.Logger is a function, not an interface
func init() {
	discordgo.Logger = func(msgL int, caller int, format string, a ...interface{}) {
		logger.Info(fmt.Sprintf(format, a...), "DiscordGo")
	}
}

// ExtendedClient wraps discordgo.Session with additional functionality
type ExtendedClient struct {
	Session        *discordgo.Session
	Commands       *CommandCollection
	CommandHandler *CommandHandler
	EventHandler   *EventHandler
	StartTime      time.Time
	mu             sync.RWMutex
	isReady        bool
}

// CommandCollection holds registered commands
type CommandCollection struct {
	commands map[string]*Command
	mu       sync.RWMutex
}

// NewCommandCollection creates a new CommandCollection
func NewCommandCollection() *CommandCollection {
	return &CommandCollection{
		commands: make(map[string]*Command),
	}
}

// Set adds or updates a command
func (cc *CommandCollection) Set(name string, cmd *Command) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.commands[name] = cmd
}

// Get retrieves a command by name
func (cc *CommandCollection) Get(name string) (*Command, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	cmd, ok := cc.commands[name]
	return cmd, ok
}

// Size returns the number of commands
func (cc *CommandCollection) Size() int {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return len(cc.commands)
}

// All returns all commands
func (cc *CommandCollection) All() map[string]*Command {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	result := make(map[string]*Command)
	for k, v := range cc.commands {
		result[k] = v
	}
	return result
}

var (
	client *ExtendedClient
	once   sync.Once
)

// Init initializes the global Discord client
func Init(token string) (*ExtendedClient, error) {
	var err error
	once.Do(func() {
		client, err = NewClient(token)
	})
	return client, err
}

// Get returns the global Discord client
func Get() *ExtendedClient {
	return client
}

// NewClient creates a new ExtendedClient
func NewClient(token string) (*ExtendedClient, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	// Set intents
	session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsMessageContent

	// Configure session
	session.ShardCount = 1 // Auto sharding equivalent
	session.SyncEvents = false
	session.StateEnabled = true
	session.LogLevel = discordgo.LogWarning

	c := &ExtendedClient{
		Session:  session,
		Commands: NewCommandCollection(),
		isReady:  false,
	}

	// Initialize handlers
	c.CommandHandler = NewCommandHandler(c)
	c.EventHandler = NewEventHandler(c)

	return c, nil
}

// Start initializes and starts the bot
func (c *ExtendedClient) Start() error {
	// Load commands
	if err := c.CommandHandler.LoadCommands(); err != nil {
		logger.Error("Failed to load commands: "+err.Error(), "Client")
		return err
	}

	// Load events
	if err := c.EventHandler.LoadEvents(); err != nil {
		logger.Error("Failed to load events: "+err.Error(), "Client")
		return err
	}

	// Add ready handler
	c.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		c.mu.Lock()
		c.isReady = true
		c.mu.Unlock()

		logger.Success("Bot conectado como: "+r.User.Username, "Client")

		// Register commands with Discord
		c.CommandHandler.RegisterCommands()

		// Señal para que Pterodactyl marque el servidor como "Running"
		fmt.Println("Sistemas Basicos Iniciados Correctamente")
	})

	// Add interaction handler
	c.Session.AddHandler(c.handleInteraction)

	// Set start time
	c.StartTime = time.Now()

	// Open connection
	err := c.Session.Open()
	if err != nil {
		return err
	}
	return nil
}

// handleInteraction handles incoming Discord interactions
func (c *ExtendedClient) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer errors.RecoverMiddleware()()

	if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		data := i.ApplicationCommandData()
		commandName := data.Name

		// Build full command name for subcommands
		if len(data.Options) > 0 {
			opt := data.Options[0]
			if opt.Type == discordgo.ApplicationCommandOptionSubCommandGroup {
				if len(opt.Options) > 0 {
					commandName = data.Name + "." + opt.Name + "." + opt.Options[0].Name
				}
			} else if opt.Type == discordgo.ApplicationCommandOptionSubCommand {
				commandName = data.Name + "." + opt.Name
			}
		}

		cmd, ok := c.Commands.Get(commandName)
		if !ok {
			return
		}

		if cmd.AutoComplete != nil {
			ctx := &CommandContext{
				Session:     s,
				Interaction: i,
				Client:      c,
			}
			if !c.canUseAutocomplete(ctx, cmd, commandName) {
				_ = ctx.SendAutoCompleteResults(nil)
				return
			}
			cmd.AutoComplete(ctx)
		}
		return
	}

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	commandName := data.Name

	// Build full command name for subcommands
	if len(data.Options) > 0 {
		opt := data.Options[0]
		if opt.Type == discordgo.ApplicationCommandOptionSubCommandGroup {
			if len(opt.Options) > 0 {
				commandName = data.Name + "." + opt.Name + "." + opt.Options[0].Name
			}
		} else if opt.Type == discordgo.ApplicationCommandOptionSubCommand {
			commandName = data.Name + "." + opt.Name
		}
	}

	cmd, ok := c.Commands.Get(commandName)
	if !ok {
		logger.Warn("Command not found: "+commandName+". Será limpiado por la sincronización (BulkOverwrite).", "Client")

		// Responder efímeramente para que no se quede "pensando"
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ Este comando ha sido movido, agrupado o ya no existe. Actualiza tu cliente de Discord o espera a que se sincronicen.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		return
	}

	// Check maintenance mode
	botConfig := config.GetBotConfig()
	if botConfig.IsMaintenance() && !strings.HasPrefix(commandName, "dev") {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ **El bot está en mantenimiento.** Los desarrolladores están trabajando en mejoras. Intenta de nuevo más tarde.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Check disabled commands
	if botConfig.IsCommandDisabled(commandName) {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ Este comando ha sido deshabilitado temporalmente por los administradores.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Execute command
	ctx := &CommandContext{
		Session:     s,
		Interaction: i,
		Client:      c,
	}

	// Check blacklist
	if err := c.checkBlacklist(ctx); err != nil {
		return
	}

	if err := c.checkCommandRequirements(ctx, cmd, commandName); err != nil {
		logger.Warn("Command requirement rejected for "+commandName+": "+err.Error(), "Client")
		return
	}

	if !cmd.PremiumType.IsNone() {
		guildID := i.GuildID
		userID := ctx.User().ID
		allowed, msg := premium.Check(cmd.PremiumType, userID, guildID)
		if !allowed {
			ctx.ReplyEphemeral(msg)
			return
		}
	}

	if err := cmd.Run(ctx); err != nil {
		logger.Error("Error executing command "+commandName+": "+err.Error(), "Client")
	}
}

func (c *ExtendedClient) canUseAutocomplete(ctx *CommandContext, cmd *Command, commandName string) bool {
	if (strings.HasPrefix(commandName, "dev.") || cmd.IsDev) && !ctx.IsDev() {
		return false
	}
	user := ctx.User()
	if user == nil {
		return false
	}
	if cmd.UserPermissions != 0 {
		permissions, err := ctx.Session.UserChannelPermissions(user.ID, ctx.Interaction.ChannelID)
		if err != nil || permissions&cmd.UserPermissions != cmd.UserPermissions {
			return false
		}
	}
	if cmd.RequiresDB {
		db := database.Get()
		if db == nil || !db.Connected() {
			return false
		}
	}
	if cmd.InVoiceChannel {
		if ctx.Session.State == nil || ctx.Interaction.GuildID == "" {
			return false
		}
		voiceState, err := ctx.Session.State.VoiceState(ctx.Interaction.GuildID, user.ID)
		if err != nil || voiceState == nil || voiceState.ChannelID == "" {
			return false
		}
	}
	return true
}

func (c *ExtendedClient) checkCommandRequirements(ctx *CommandContext, cmd *Command, commandName string) error {
	if strings.HasPrefix(commandName, "dev.") || cmd.IsDev {
		if !ctx.IsDev() {
			_ = ctx.ReplyEphemeral("❌ Este comando está reservado para el equipo de desarrollo.")
			return fmt.Errorf("developer access required")
		}
	}

	if cmd.UserPermissions != 0 {
		user := ctx.User()
		if user == nil {
			_ = ctx.ReplyEphemeral("❌ No se pudo validar tu identidad.")
			return fmt.Errorf("interaction user is missing")
		}
		permissions, err := ctx.Session.UserChannelPermissions(user.ID, ctx.Interaction.ChannelID)
		if err != nil || permissions&cmd.UserPermissions != cmd.UserPermissions {
			_ = ctx.ReplyEphemeral("❌ No tienes los permisos necesarios para usar este comando.")
			return fmt.Errorf("missing user permissions")
		}
	}

	if cmd.BotPermissions != 0 {
		if ctx.Session.State == nil || ctx.Session.State.User == nil {
			_ = ctx.ReplyEphemeral("❌ No se pudieron validar los permisos del bot.")
			return fmt.Errorf("bot state is unavailable")
		}
		permissions, err := ctx.Session.UserChannelPermissions(ctx.Session.State.User.ID, ctx.Interaction.ChannelID)
		if err != nil || permissions&cmd.BotPermissions != cmd.BotPermissions {
			_ = ctx.ReplyEphemeral("❌ No tengo los permisos necesarios para ejecutar este comando.")
			return fmt.Errorf("missing bot permissions")
		}
	}

	if cmd.RequiresDB {
		db := database.Get()
		if db == nil || !db.Connected() {
			_ = ctx.ReplyEphemeral("❌ La base de datos no está disponible temporalmente.")
			return fmt.Errorf("database is unavailable")
		}
	}

	if cmd.InVoiceChannel {
		user := ctx.User()
		if user == nil || ctx.Interaction.GuildID == "" || ctx.Session.State == nil {
			_ = ctx.ReplyEphemeral("❌ Debes usar este comando dentro de un servidor y estar en un canal de voz.")
			return fmt.Errorf("voice context is unavailable")
		}
		voiceState, err := ctx.Session.State.VoiceState(ctx.Interaction.GuildID, user.ID)
		if err != nil || voiceState == nil || voiceState.ChannelID == "" {
			_ = ctx.ReplyEphemeral("❌ Debes estar en un canal de voz para usar este comando.")
			return fmt.Errorf("user is not in a voice channel")
		}
	}

	return nil
}

// Stop stops the bot and closes the session
func (c *ExtendedClient) Stop() error {
	c.mu.Lock()
	c.isReady = false
	c.mu.Unlock()

	if c.Session != nil {
		return c.Session.Close()
	}
	return nil
}

// IsReady returns true if the bot is ready
func (c *ExtendedClient) IsReady() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isReady
}

// GuildCount returns the number of guilds the bot is in
func (c *ExtendedClient) GuildCount() int {
	if c.Session == nil || c.Session.State == nil {
		return 0
	}
	c.Session.State.RLock()
	defer c.Session.State.RUnlock()
	return len(c.Session.State.Guilds)
}

// GetConfig returns the bot configuration
func (c *ExtendedClient) GetConfig() *config.Config {
	return config.Get()
}

// checkBlacklist verifica si el usuario o guild está en la blacklist
func (c *ExtendedClient) checkBlacklist(ctx *CommandContext) error {
	return c.BlacklistMiddleware(ctx)
}
