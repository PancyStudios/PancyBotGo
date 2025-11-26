// Package main provides a utility to sync Discord slash commands.
// This removes stale commands from Discord and ensures only currently-defined commands are registered.
//
// Usage:
//   go run cmd/sync-commands/main.go [options]
//
// Options:
//   -list           List all registered commands (global and guild)
//   -clean          Remove all commands without registering new ones
//   -guild <id>     Target a specific guild instead of global commands
//   -sync           Sync commands (remove stale, register current) - default behavior
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/PancyStudios/PancyBotGo/internal/commands"
	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

func main() {
	// Parse command line flags
	listCmd := flag.Bool("list", false, "List all registered commands")
	cleanCmd := flag.Bool("clean", false, "Remove all commands without registering new ones")
	guildID := flag.String("guild", "", "Target a specific guild (leave empty for global)")
	syncCmd := flag.Bool("sync", false, "Sync commands (remove stale, register current)")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.Init(cfg.ErrorWebhook, cfg.LogsWebhook)
	defer log.Close()

	logger.System("Iniciando utilidad de sincronización de comandos...", "SyncCommands")

	// Initialize Discord client
	client, err := discord.NewClient(cfg.BotToken)
	if err != nil {
		logger.Critical(fmt.Sprintf("Error creating Discord client: %v", err), "SyncCommands")
		os.Exit(1)
	}

	// Open connection to Discord
	if err := client.Session.Open(); err != nil {
		logger.Critical(fmt.Sprintf("Error connecting to Discord: %v", err), "SyncCommands")
		os.Exit(1)
	}
	defer client.Session.Close()

	logger.Success("Conectado a Discord", "SyncCommands")

	// Register commands to know what we should have
	commands.RegisterAll(client)

	// Execute the requested action
	switch {
	case *listCmd:
		listCommands(client, *guildID)
	case *cleanCmd:
		cleanCommands(client, *guildID)
	case *syncCmd:
		syncCommands(client, *guildID)
	default:
		// Default: sync commands
		syncCommands(client, *guildID)
	}

	logger.Success("Operación completada exitosamente", "SyncCommands")
}

// listCommands lists all commands registered with Discord
func listCommands(client *discord.ExtendedClient, guildID string) {
	logger.Info("📋 Listando comandos registrados...", "SyncCommands")

	var cmds []*discordgo.ApplicationCommand
	var err error

	if guildID != "" {
		logger.Info(fmt.Sprintf("Obteniendo comandos del servidor: %s", guildID), "SyncCommands")
		cmds, err = client.CommandHandler.ListGuildCommands(guildID)
	} else {
		logger.Info("Obteniendo comandos globales", "SyncCommands")
		cmds, err = client.CommandHandler.ListGlobalCommands()
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Error obteniendo comandos: %v", err), "SyncCommands")
		return
	}

	if len(cmds) == 0 {
		logger.Info("No hay comandos registrados", "SyncCommands")
		return
	}

	logger.Info(fmt.Sprintf("Comandos encontrados: %d", len(cmds)), "SyncCommands")
	for i, cmd := range cmds {
		logger.Info(fmt.Sprintf("  %d. /%s - %s (ID: %s)", i+1, cmd.Name, cmd.Description, cmd.ID), "SyncCommands")
	}
}

// cleanCommands removes all commands from Discord
func cleanCommands(client *discord.ExtendedClient, guildID string) {
	logger.Info("🧹 Eliminando todos los comandos...", "SyncCommands")

	var err error
	if guildID != "" {
		logger.Info(fmt.Sprintf("Eliminando comandos del servidor: %s", guildID), "SyncCommands")
		err = client.CommandHandler.UnregisterGuildCommands(guildID)
	} else {
		logger.Info("Eliminando comandos globales", "SyncCommands")
		err = client.CommandHandler.UnregisterCommands()
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Error eliminando comandos: %v", err), "SyncCommands")
		return
	}

	logger.Success("✅ Todos los comandos han sido eliminados", "SyncCommands")
}

// syncCommands removes stale commands and registers current ones
func syncCommands(client *discord.ExtendedClient, guildID string) {
	logger.Info("🔄 Sincronizando comandos...", "SyncCommands")

	if guildID != "" {
		logger.Info(fmt.Sprintf("Sincronizando comandos del servidor: %s", guildID), "SyncCommands")
		// Remove guild commands
		if err := client.CommandHandler.UnregisterGuildCommands(guildID); err != nil {
			logger.Error(fmt.Sprintf("Error eliminando comandos de guild: %v", err), "SyncCommands")
			return
		}
		// Note: For guild-specific registration, you would need to modify RegisterCommands
		// to support guild-specific registration
		logger.Warn("Nota: La registración de comandos específicos de guild no está completamente implementada", "SyncCommands")
	} else {
		// Sync global commands
		if err := client.CommandHandler.SyncCommands(); err != nil {
			logger.Error(fmt.Sprintf("Error sincronizando comandos: %v", err), "SyncCommands")
			return
		}
	}

	logger.Success("✅ Comandos sincronizados correctamente", "SyncCommands")
}
