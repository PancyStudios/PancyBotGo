// Package main is the entry point for the PancyBot Go application.
// It initializes all systems and starts the Discord bot.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PancyStudios/PancyBotGo/internal/commands"
	"github.com/PancyStudios/PancyBotGo/internal/events"
	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/errors"
	"github.com/PancyStudios/PancyBotGo/pkg/lavalink"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/mqtt"
	"github.com/PancyStudios/PancyBotGo/pkg/web"
)

func main() {
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

		// Initialize blacklist cache at startup and start auto-refresh
		if err := database.InitBlacklistCache(); err != nil {
			logger.Warn(fmt.Sprintf("Error inicializando caché de blacklist: %v", err), "Main")
		}
		database.StartBlacklistCacheRefresh()
		defer database.StopBlacklistCacheRefresh()
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

	// Initialize web server
	webServer := web.Init(cfg.LogsWebServerHook)
	web.SetupAPIRoutes(webServer)
	webServer.StartAsync(cfg.Port)

	// Initialize Discord client
	discordClient, err = discord.Init(cfg.BotToken)
	if err != nil {
		logger.Critical(fmt.Sprintf("Error creating Discord client: %v", err), "Main")
		os.Exit(1)
	}

	// Register commands using the new commands package
	commands.RegisterAll(discordClient)

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

	logger.Success("PancyBot Go iniciado correctamente!", "Main")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

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
