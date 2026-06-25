// Package discord provides the command handler for loading and registering commands.
package discord

import (
	"sync"

	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// CommandHandler manages command loading and registration
type CommandHandler struct {
	client           *ExtendedClient
	slashCommands    []*discordgo.ApplicationCommand
	slashCommandsDev []*discordgo.ApplicationCommand
	mu               sync.RWMutex
}

// NewCommandHandler creates a new CommandHandler
func NewCommandHandler(client *ExtendedClient) *CommandHandler {
	return &CommandHandler{
		client:           client,
		slashCommands:    make([]*discordgo.ApplicationCommand, 0),
		slashCommandsDev: make([]*discordgo.ApplicationCommand, 0),
	}
}

// LoadCommands loads all commands from the commands registry
// In Go, we register commands programmatically instead of reading from files
func (ch *CommandHandler) LoadCommands() error {
	logger.System("Iniciando carga de comandos...", "CommandHandler")

	// Commands are registered programmatically using RegisterCommand
	// Example commands can be added here or in separate packages

	logger.System("Carga finalizada. Los comandos se registrarán programáticamente.", "CommandHandler")
	return nil
}

// RegisterCommand adds a command to the handler
func (ch *CommandHandler) RegisterCommand(cmd *Command) {
	ch.client.Commands.Set(cmd.Name, cmd)

	appCmd := cmd.ToApplicationCommand()

	ch.mu.Lock()
	if cmd.IsDev {
		ch.slashCommandsDev = append(ch.slashCommandsDev, appCmd)
	} else {
		ch.slashCommands = append(ch.slashCommands, appCmd)
	}
	ch.mu.Unlock()

	logger.Debug("Comando registrado: "+cmd.Name, "CommandHandler")
}

// RegisterSubcommand adds a subcommand to an existing command group
func (ch *CommandHandler) RegisterSubcommand(groupName string, cmd *Command) {
	fullName := groupName + "." + cmd.Name
	ch.client.Commands.Set(fullName, cmd)
	logger.Debug("Subcomando registrado: "+fullName, "CommandHandler")
}

// RegisterSubcommandGroup adds a subcommand group
func (ch *CommandHandler) RegisterSubcommandGroup(groupName, subgroupName string, cmd *Command) {
	fullName := groupName + "." + subgroupName + "." + cmd.Name
	ch.client.Commands.Set(fullName, cmd)
	logger.Debug("Subcomando de grupo registrado: "+fullName, "CommandHandler")
}

// BuildCommandGroup creates a command group with subcommands
func (ch *CommandHandler) BuildCommandGroup(name, description string, subcommands ...*Command) *discordgo.ApplicationCommand {
	options := make([]*discordgo.ApplicationCommandOption, 0, len(subcommands))

	for _, cmd := range subcommands {
		fullName := name + "." + cmd.Name
		ch.client.Commands.Set(fullName, cmd)

		optType := discordgo.ApplicationCommandOptionSubCommand
		if len(cmd.Options) > 0 && cmd.Options[0].Type == discordgo.ApplicationCommandOptionSubCommand {
			optType = discordgo.ApplicationCommandOptionSubCommandGroup
			// Also, we need to register the routing for subcommands inside the group
			for _, subOpt := range cmd.Options {
				if subOpt.Type == discordgo.ApplicationCommandOptionSubCommand {
					subFullName := fullName + "." + subOpt.Name
					// the handler logic might need to route differently, but let's just map it to the group root for now
					ch.client.Commands.Set(subFullName, cmd)
				}
			}
		}

		opt := &discordgo.ApplicationCommandOption{
			Type:        optType,
			Name:        cmd.Name,
			Description: cmd.Description,
			Options:     cmd.Options,
		}
		options = append(options, opt)
	}

	return &discordgo.ApplicationCommand{
		Name:        name,
		Description: description,
		Options:     options,
	}
}

// BuildSubcommandGroup creates a subcommand group
func (ch *CommandHandler) BuildSubcommandGroup(groupName, name, description string, subcommands ...*Command) *discordgo.ApplicationCommandOption {
	options := make([]*discordgo.ApplicationCommandOption, 0, len(subcommands))

	for _, cmd := range subcommands {
		fullName := groupName + "." + name + "." + cmd.Name
		ch.client.Commands.Set(fullName, cmd)

		opt := &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        cmd.Name,
			Description: cmd.Description,
			Options:     cmd.Options,
		}
		options = append(options, opt)
	}

	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
		Name:        name,
		Description: description,
		Options:     options,
	}
}

// RegisterCommands registers all slash commands with Discord
func (ch *CommandHandler) RegisterCommands() {
	cfg := config.Get()

	logger.Info("🔄 Registrando comandos globales...", "CommandHandler")

	// Register global commands
	ch.mu.RLock()
	slashCommands := make([]*discordgo.ApplicationCommand, len(ch.slashCommands))
	copy(slashCommands, ch.slashCommands)
	slashCommandsDev := make([]*discordgo.ApplicationCommand, len(ch.slashCommandsDev))
	copy(slashCommandsDev, ch.slashCommandsDev)
	ch.mu.RUnlock()

	_, err := ch.client.Session.ApplicationCommandBulkOverwrite(
		ch.client.Session.State.User.ID,
		"",
		slashCommands,
	)
	if err != nil {
		logger.Error("Error en BulkOverwrite global: "+err.Error(), "CommandHandler")
	} else {
		logger.Success("✅ Comandos globales registrados y sincronizados.", "CommandHandler")
	}

	// Register dev commands in dev guild
	if cfg.DevGuildID != "" && len(slashCommandsDev) > 0 {
		logger.Info("🔄 Registrando comandos de desarrollo en el servidor "+cfg.DevGuildID+"...", "CommandHandler")

		_, err := ch.client.Session.ApplicationCommandBulkOverwrite(
			ch.client.Session.State.User.ID,
			cfg.DevGuildID,
			slashCommandsDev,
		)
		if err != nil {
			logger.Error("Error en BulkOverwrite dev: "+err.Error(), "CommandHandler")
		}

		logger.Success("✅ Comandos de desarrollo registrados.", "CommandHandler")
	}
}

// UnregisterCommands removes all registered commands from Discord
func (ch *CommandHandler) UnregisterCommands() error {
	commands, err := ch.client.Session.ApplicationCommands(ch.client.Session.State.User.ID, "")
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		err := ch.client.Session.ApplicationCommandDelete(ch.client.Session.State.User.ID, "", cmd.ID)
		if err != nil {
			logger.Error("Error eliminando comando "+cmd.Name+": "+err.Error(), "CommandHandler")
		}
	}

	logger.Success("Comandos globales eliminados.", "CommandHandler")
	return nil
}

// UnregisterGuildCommands removes all guild-specific commands from Discord
func (ch *CommandHandler) UnregisterGuildCommands(guildID string) error {
	commands, err := ch.client.Session.ApplicationCommands(ch.client.Session.State.User.ID, guildID)
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		err := ch.client.Session.ApplicationCommandDelete(ch.client.Session.State.User.ID, guildID, cmd.ID)
		if err != nil {
			logger.Error("Error eliminando comando de guild "+cmd.Name+": "+err.Error(), "CommandHandler")
		}
	}

	logger.Success("Comandos de guild eliminados para "+guildID, "CommandHandler")
	return nil
}

// ListGlobalCommands returns all global commands registered with Discord
func (ch *CommandHandler) ListGlobalCommands() ([]*discordgo.ApplicationCommand, error) {
	return ch.client.Session.ApplicationCommands(ch.client.Session.State.User.ID, "")
}

// ListGuildCommands returns all guild-specific commands registered with Discord
func (ch *CommandHandler) ListGuildCommands(guildID string) ([]*discordgo.ApplicationCommand, error) {
	return ch.client.Session.ApplicationCommands(ch.client.Session.State.User.ID, guildID)
}

// SyncCommands removes stale commands and registers only the current ones
// This ensures Discord only shows commands that are currently defined
func (ch *CommandHandler) SyncCommands() error {
	logger.Info("🔄 Sincronizando comandos con Discord...", "CommandHandler")

	// Remove all existing global commands first
	if err := ch.UnregisterCommands(); err != nil {
		logger.Error("Error eliminando comandos existentes: "+err.Error(), "CommandHandler")
		return err
	}

	// Register current commands (RegisterCommands doesn't return an error)
	ch.RegisterCommands()

	logger.Success("✅ Sincronización de comandos completada", "CommandHandler")
	return nil
}

// AddGlobalCommand adds a command to the global command list
func (ch *CommandHandler) AddGlobalCommand(cmd *discordgo.ApplicationCommand) {
	ch.mu.Lock()
	ch.slashCommands = append(ch.slashCommands, cmd)
	ch.mu.Unlock()
}

// AddDevCommand adds a command to the dev command list
func (ch *CommandHandler) AddDevCommand(cmd *discordgo.ApplicationCommand) {
	ch.mu.Lock()
	ch.slashCommandsDev = append(ch.slashCommandsDev, cmd)
	ch.mu.Unlock()
}

// GetRegisteredCommands returns the list of global application commands registered in memory
func (ch *CommandHandler) GetRegisteredCommands() []*discordgo.ApplicationCommand {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.slashCommands
}
