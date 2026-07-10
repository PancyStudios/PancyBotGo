// Package main is the entry point for the PancyBot Go application.
// It initializes all systems and starts the Discord bot.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PancyStudios/PancyBotGo/internal/commands"
	"github.com/PancyStudios/PancyBotGo/internal/events"
	msgconfig "github.com/PancyStudios/PancyBotGo/internal/messagecommands/config"
	"github.com/PancyStudios/PancyBotGo/internal/messagecommands/economy"
	funMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/fun"
	modMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/mod"
	utilsMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/utils"
	helpMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/help"
	iaMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/ia"
	levelsMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/levels"
	devMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/dev"
	embedsMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/embeds"
	musicMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/music"
	premiumMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/premium"
	reactionMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/reaction"
	securityMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/security"
	"github.com/PancyStudios/PancyBotGo/pkg/api"
	"github.com/PancyStudios/PancyBotGo/pkg/cli"
	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/lavalink"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/mqtt"
	"github.com/PancyStudios/PancyBotGo/pkg/scheduler"
	"github.com/PancyStudios/PancyBotGo/pkg/web"
)

func main() {
	// Parse flags
	useDashboard := flag.Bool("dashboard", false, "Abre el panel de control local en el navegador")
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

	logger.System("Iniciando PancyBot Go...", "Main")
	logger.Info(fmt.Sprintf("Directorio de trabajo: %s", getCurrentDir()), "Main")

	// Initialize error handler
	var discordClient *discord.ExtendedClient
	var lavalinkClient *lavalink.LavalinkClient
	errors.Init(cfg.ErrorWebhook, func() {
		if discordClient != nil {
			err := discordClient.Stop()
			if err != nil {
				return
			}
		}
		if lavalinkClient != nil {
			lavalinkClient.Disconnect()
		}
	})

	// Initialize database
	db, err := database.Init(cfg.MongoDBURL, cfg.DBName)
	if err != nil {
		logger.Error(fmt.Sprintf("Error connecting to database: %v", err), "Main")
		logger.Debug(fmt.Sprintf("Error connecting to database: %v", cfg.MongoDBURL), "Main")
		// Continue without database- it will attempt to reconnect
	}
	defer func() {
		if db != nil {
			err := db.Disconnect()
			if err != nil {
				return
			}
		}
	}()

	// Initialize global DataManagers
	if db != nil {
		database.InitGlobalDataManagers(db)

		// Initialize blacklist cache with all entries from DB
		if err := database.InitBlacklistCache(); err != nil {
			logger.Warn(fmt.Sprintf("Error initializing blacklist cache: %v", err), "Main")
		}

		// Start automatic blacklist cache refresh every 5 minutes
		database.StartBlacklistCacheRefresh()
	}

	// Initialize MQTT
	mqttClientID := "pancybot"
	if !cfg.IsProd() {
		mqttClientID = "pancybot_canary"
	}

	mqttClient := mqtt.Init(
		cfg.MQTTHost,
		cfg.MQTTPort,
		cfg.MQTTUser,
		cfg.MQTTPassword,
		mqttClientID,
	)
	defer mqttClient.Destroy()

	// ──────────────────────────────────────────
	// Publish logs to MQTT
	// ──────────────────────────────────────────
	go func() {
		logChan := logger.Subscribe()
		for logMsg := range logChan {
			payload := map[string]interface{}{
				"message": logMsg,
			}
			mqttClient.Publish(fmt.Sprintf("pancy/logs/%s", cfg.Environment), payload)
		}
	}()

	// Initialize web server
	webServer := web.Init(cfg.LogsWebServerHook)
	web.SetupAPIRoutes(webServer)
	web.SetupDashboardRoutes(webServer)
	webServer.StartAsync(cfg.Port)

	// Open dashboard automatically if requested
	if *useDashboard {
		dashboardURL := fmt.Sprintf("http://localhost:%s/admin", cfg.Port)
		cli.OpenBrowser(dashboardURL)
	}

	// Initialize Discord client
	discordClient, err = discord.Init(cfg.BotToken)
	if err != nil {
		logger.Critical(fmt.Sprintf("Error creating Discord client: %v", err), "Main")
		os.Exit(1)
	}

	// Register commands using the new commands package
	commands.RegisterAll(discordClient)

	// Register text commands (Prefix)
	msgconfig.Register()
	economy.Register()
	funMsgCommands.Register()
	modMsgCommands.Register()
	utilsMsgCommands.Register()
	helpMsgCommands.RegisterAll()
	iaMsgCommands.RegisterAll()
	levelsMsgCommands.RegisterAll()
	devMsgCommands.RegisterAll()
	embedsMsgCommands.RegisterAll()
	musicMsgCommands.RegisterAll()
	premiumMsgCommands.RegisterAll()
	reactionMsgCommands.RegisterAll()
	securityMsgCommands.RegisterAll()

	logger.Info("Initialising Discord Bot...", "App")

	// Register events using the new events package
	events.RegisterAll(discordClient)

	// Start the bot
	if err := discordClient.Start(); err != nil {
		logger.Critical(fmt.Sprintf("Error starting Discord client: %v", err), "Main")
		os.Exit(1)
	}
	defer func(discordClient *discord.ExtendedClient) {
		err := discordClient.Stop()
		if err != nil {

		}
	}(discordClient)

	// Start Tempban scheduler
	scheduler.StartTempBanScheduler(discordClient)

	// Initialize Lavalink after Discord is connected
	lavalinkClient = lavalink.Init(discordClient.Session, []lavalink.NodeConfig{
		{
			Name:     "PancyBeta",
			Host:     cfg.LinkServer,
			Port:     2333,
			Password: cfg.LinkPassword,
			Secure:   false,
		},
	})

	err = lavalinkClient.Connect()
	if err != nil {
		return
	}
	defer lavalinkClient.Disconnect()

	// Register MQTT handlers for remote control via API
	lavalink.RegisterMusicHandlers(mqttClient, lavalinkClient)
	api.RegisterAPIHandlers(mqttClient, discordClient)
	api.RegisterDevHandlers(mqttClient, discordClient)

	logger.Success("PancyBot Go iniciado correctamente!", "Main")

	// Start CLI REPL
	cli.Start(discordClient)

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Stop blacklist cache auto-refresh
	database.StopBlacklistCacheRefresh()

	logger.System("Apagando PancyBot Go...", "Main")
}

// getCurrentDir returns the current working directory
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}
