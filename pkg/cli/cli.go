package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

var discordClient *discord.ExtendedClient

// Start begins the interactive REPL in a non-blocking goroutine
func Start(client *discord.ExtendedClient) {
	discordClient = client
	go func() {
		// Wait a moment so startup logs don't clutter the initial prompt
		time.Sleep(1 * time.Second)

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				continue
			}

			if ExecuteCommand(input) {
				return // stop command was executed
			}
		}

		if err := scanner.Err(); err != nil {
			logger.Error(fmt.Sprintf("Error leyendo consola: %v", err), "CLI")
		}
	}()
}

// ExecuteCommand executes a CLI command and returns true if the bot should stop
func ExecuteCommand(input string) bool {
	args := strings.Fields(input)
	if len(args) == 0 {
		return false
	}

	command := strings.ToLower(args[0])

	switch command {
	case "help":
		showHelp()
	case "stats":
		showStats()
	case "clear":
		clearScreen()
	case "premium":
		handlePremium(args[1:])
	case "blacklist":
		handleBlacklist(args[1:])
	case "guilds":
		handleGuilds(args[1:])
	case "user":
		handleUserInfo(args[1:])
	case "guild":
		handleGuildInfo(args[1:])
	case "msg":
		handleMsg(args[1:])
	case "ping":
		if discordClient != nil && discordClient.Session != nil {
			logger.System(fmt.Sprintf("Discord API Latency: %v", discordClient.Session.HeartbeatLatency()), "CLI")
		} else {
			logger.System("Cliente Discord no está listo.", "CLI")
		}
	case "cache":
		if len(args) > 1 && strings.ToLower(args[1]) == "reload" {
			database.InitBlacklistCache()
			// database.InitGuildConfigCache() (if exists, but we can't be sure, so we just clear them)
			database.GlobalGuildDM.ClearCache()
			database.GlobalUserPremiumDM.ClearCache()
			database.GlobalGuildPremiumDM.ClearCache()
			database.GlobalBlacklistDM.ClearCache()
			logger.System("Cachés de MongoDB recargadas exitosamente.", "CLI")
		} else {
			logger.System("Uso: cache reload", "CLI")
		}
	case "stop", "exit", "quit":
		logger.System("Iniciando apagado seguro...", "CLI")
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		return true
	default:
		logger.System(fmt.Sprintf("Comando desconocido: '%s'. Escribe 'help' para ver la lista de comandos.", command), "CLI")
	}
	return false
}

func showHelp() {
	msg := `=== PancyBot CLI Comandos ===
  help       - Muestra esta lista
  stats      - Muestra estadísticas del sistema
  clear      - Limpia la consola
  stop       - Apaga el bot de forma segura
  
  premium user [add/remove/list] <id>   - Gestiona usuarios premium
  premium guild [add/remove/list] <id>  - Gestiona servidores premium
  blacklist [add/remove/list] <id>      - Gestiona lista negra de usuarios
  guilds [number/exit] <id>             - Gestiona servidores activos
  
  user <id>                  - Muestra información completa de un usuario
  guild <id>                 - Muestra información completa de un servidor
  ping                       - Muestra la latencia con Discord
  cache reload               - Limpia la memoria caché de la DB
  msg <canal_id> <texto>     - Envía un mensaje como el bot
=============================`
	logger.System(msg, "CLI")
}

func showStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	msg := fmt.Sprintf(`=== Estadísticas del Sistema ===
  Goroutines  : %d
  Memoria OS  : %.2f MB
  Memoria Uso : %.2f MB
  GC Pausas   : %d
================================`, runtime.NumGoroutine(), float64(m.Sys)/1024/1024, float64(m.Alloc)/1024/1024, m.NumGC)
	logger.System(msg, "CLI")
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func handlePremium(args []string) {
	if len(args) < 2 {
		logger.System("Uso: premium <user|guild> <add|remove|list> [id]", "CLI")
		return
	}

	target := strings.ToLower(args[0])
	action := strings.ToLower(args[1])
	id := ""
	if len(args) > 2 {
		id = args[2]
	}

	if target == "user" {
		switch action {
		case "add":
			if id == "" {
				logger.System("Falta el ID del usuario", "CLI")
				return
			}
			doc := &models.UserPremium{UserID: id, Permanent: true}
			_, err := database.GlobalUserPremiumDM.Set(bson.M{"user": id}, doc)
			if err != nil {
				logger.Error(fmt.Sprintf("Error añadiendo premium: %v", err), "CLI")
				return
			}
			logger.System(fmt.Sprintf("Usuario %s añadido a Premium", id), "CLI")
		case "remove":
			if id == "" {
				logger.System("Falta el ID del usuario", "CLI")
				return
			}
			err := database.GlobalUserPremiumDM.Delete(bson.M{"user": id})
			if err != nil {
				logger.Error(fmt.Sprintf("Error removiendo premium: %v", err), "CLI")
				return
			}
			logger.System(fmt.Sprintf("Usuario %s removido de Premium", id), "CLI")
		case "list":
			users, err := database.GlobalUserPremiumDM.GetAll(bson.M{})
			if err != nil {
				logger.Error(fmt.Sprintf("Error obteniendo lista premium: %v", err), "CLI")
				return
			}
			if len(users) == 0 {
				logger.System("No hay usuarios premium.", "CLI")
				return
			}
			msg := "=== Usuarios Premium ===\n"
			for i, u := range users {
				if i >= 20 {
					msg += fmt.Sprintf("... y %d más\n", len(users)-20)
					break
				}
				msg += fmt.Sprintf("- %s (Permanente: %v)\n", u.UserID, u.Permanent)
			}
			logger.System(msg, "CLI")
		}
	} else if target == "guild" {
		switch action {
		case "add":
			if id == "" {
				logger.System("Falta el ID del servidor", "CLI")
				return
			}
			doc := &models.GuildPremium{GuildID: id, Permanent: true}
			_, err := database.GlobalGuildPremiumDM.Set(bson.M{"guild": id}, doc)
			if err != nil {
				logger.Error(fmt.Sprintf("Error añadiendo premium: %v", err), "CLI")
				return
			}
			logger.System(fmt.Sprintf("Servidor %s añadido a Premium", id), "CLI")
		case "remove":
			if id == "" {
				logger.System("Falta el ID del servidor", "CLI")
				return
			}
			err := database.GlobalGuildPremiumDM.Delete(bson.M{"guild": id})
			if err != nil {
				logger.Error(fmt.Sprintf("Error removiendo premium: %v", err), "CLI")
				return
			}
			logger.System(fmt.Sprintf("Servidor %s removido de Premium", id), "CLI")
		case "list":
			guilds, err := database.GlobalGuildPremiumDM.GetAll(bson.M{})
			if err != nil {
				logger.Error(fmt.Sprintf("Error obteniendo lista premium: %v", err), "CLI")
				return
			}
			if len(guilds) == 0 {
				logger.System("No hay servidores premium.", "CLI")
				return
			}
			msg := "=== Servidores Premium ===\n"
			for i, g := range guilds {
				if i >= 20 {
					msg += fmt.Sprintf("... y %d más\n", len(guilds)-20)
					break
				}
				msg += fmt.Sprintf("- %s (Permanente: %v)\n", g.GuildID, g.Permanent)
			}
			logger.System(msg, "CLI")
		}
	} else {
		logger.System("Objetivo inválido. Usa 'user' o 'guild'", "CLI")
	}
}

func handleBlacklist(args []string) {
	if len(args) < 1 {
		logger.System("Uso: blacklist <add|remove|list> [id]", "CLI")
		return
	}

	action := strings.ToLower(args[0])
	id := ""
	if len(args) > 1 {
		id = args[1]
	}

	switch action {
	case "add":
		if id == "" {
			logger.System("Falta el ID del usuario", "CLI")
			return
		}
		doc := &models.Blacklist{ID: id, Type: models.BlacklistTypeUser, Reason: "Añadido desde CLI", AddedAt: time.Now()}
		_, err := database.GlobalBlacklistDM.Set(bson.M{"_id": id}, doc)
		if err != nil {
			logger.Error(fmt.Sprintf("Error añadiendo blacklist: %v", err), "CLI")
			return
		}
		// Refresh cache
		database.InitBlacklistCache()
		logger.System(fmt.Sprintf("Usuario %s añadido a Blacklist", id), "CLI")
	case "remove":
		if id == "" {
			logger.System("Falta el ID del usuario", "CLI")
			return
		}
		err := database.GlobalBlacklistDM.Delete(bson.M{"_id": id})
		if err != nil {
			logger.Error(fmt.Sprintf("Error removiendo blacklist: %v", err), "CLI")
			return
		}
		// Refresh cache
		database.InitBlacklistCache()
		logger.System(fmt.Sprintf("Usuario %s removido de Blacklist", id), "CLI")
	case "list":
		bans, err := database.GlobalBlacklistDM.GetAll(bson.M{})
		if err != nil {
			logger.Error(fmt.Sprintf("Error obteniendo blacklist: %v", err), "CLI")
			return
		}
		if len(bans) == 0 {
			logger.System("La Blacklist está vacía.", "CLI")
			return
		}
		msg := "=== Blacklist ===\n"
		for i, b := range bans {
			if i >= 20 {
				msg += fmt.Sprintf("... y %d más\n", len(bans)-20)
				break
			}
			msg += fmt.Sprintf("- [%s] %s (Razón: %s)\n", b.Type, b.ID, b.Reason)
		}
		logger.System(msg, "CLI")
	}
}

func handleGuilds(args []string) {
	if discordClient == nil || discordClient.Session == nil {
		logger.System("Cliente Discord no está listo.", "CLI")
		return
	}

	if len(args) < 1 {
		logger.System("Uso: guilds <number|exit> [id]", "CLI")
		return
	}

	action := strings.ToLower(args[0])
	switch action {
	case "number":
		guilds := len(discordClient.Session.State.Guilds)
		logger.System(fmt.Sprintf("PancyBot está actualmente en %d servidores.", guilds), "CLI")
	case "exit":
		if len(args) < 2 {
			logger.System("Falta el ID del servidor a abandonar", "CLI")
			return
		}
		id := args[1]
		err := discordClient.Session.GuildLeave(id)
		if err != nil {
			logger.Error(fmt.Sprintf("Error abandonando servidor %s: %v", id, err), "CLI")
			return
		}
		logger.System(fmt.Sprintf("Servidor %s abandonado con éxito.", id), "CLI")
	default:
		logger.System("Uso: guilds <number|exit> [id]", "CLI")
	}
}

func handleUserInfo(args []string) {
	if discordClient == nil || discordClient.Session == nil {
		logger.System("Cliente Discord no está listo.", "CLI")
		return
	}
	if len(args) < 1 {
		logger.System("Uso: user <id>", "CLI")
		return
	}
	id := args[0]
	user, err := discordClient.Session.User(id)

	msg := fmt.Sprintf("=== Información de Usuario ===\nID: %s\n", id)
	if err == nil {
		msg += fmt.Sprintf("Tag: %s#%s\nBot: %v\n", user.Username, user.Discriminator, user.Bot)
	} else {
		msg += "Usuario no encontrado en Discord.\n"
	}

	prem, err := database.GlobalUserPremiumDM.Get(bson.M{"user": id})
	if err == nil && prem != nil {
		msg += fmt.Sprintf("Premium: Sí (Permanente: %v)\n", prem.Permanent)
	} else {
		msg += "Premium: No\n"
	}

	ban, err := database.GlobalBlacklistDM.Get(bson.M{"_id": id})
	if err == nil && ban != nil && ban.Type == models.BlacklistTypeUser {
		msg += fmt.Sprintf("Blacklist: Sí (Razón: %s)\n", ban.Reason)
	} else {
		msg += "Blacklist: No\n"
	}

	logger.System(msg+"==============================", "CLI")
}

func handleGuildInfo(args []string) {
	if discordClient == nil || discordClient.Session == nil {
		logger.System("Cliente Discord no está listo.", "CLI")
		return
	}
	if len(args) < 1 {
		logger.System("Uso: guild <id>", "CLI")
		return
	}
	id := args[0]
	guild, err := discordClient.Session.Guild(id)

	msg := fmt.Sprintf("=== Información de Servidor ===\nID: %s\n", id)
	if err == nil {
		msg += fmt.Sprintf("Nombre: %s\nMiembros: %d\nDueño ID: %s\n", guild.Name, guild.MemberCount, guild.OwnerID)
	} else {
		msg += "Servidor no encontrado en la caché de Discord (¿el bot está ahí?).\n"
	}

	prem, err := database.GlobalGuildPremiumDM.Get(bson.M{"guild": id})
	if err == nil && prem != nil {
		msg += fmt.Sprintf("Premium: Sí (Permanente: %v)\n", prem.Permanent)
	} else {
		msg += "Premium: No\n"
	}

	logger.System(msg+"===============================", "CLI")
}

func handleMsg(args []string) {
	if discordClient == nil || discordClient.Session == nil {
		logger.System("Cliente Discord no está listo.", "CLI")
		return
	}
	if len(args) < 2 {
		logger.System("Uso: msg <canal_id> <texto...>", "CLI")
		return
	}
	channelID := args[0]
	text := strings.Join(args[1:], " ")

	_, err := discordClient.Session.ChannelMessageSend(channelID, text)
	if err != nil {
		logger.Error(fmt.Sprintf("Error enviando mensaje: %v", err), "CLI")
	} else {
		logger.System(fmt.Sprintf("Mensaje enviado al canal %s", channelID), "CLI")
	}
}
