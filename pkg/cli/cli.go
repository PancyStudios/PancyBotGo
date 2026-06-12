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

	"github.com/PancyStudios/PancyBotGo/pkg/logger"
)

// Start begins the interactive REPL in a non-blocking goroutine
func Start() {
	go func() {
		// Wait a moment so startup logs don't clutter the initial prompt
		time.Sleep(1 * time.Second)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("\n> ")

		for scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				fmt.Print("> ")
				continue
			}

			args := strings.Fields(input)
			command := strings.ToLower(args[0])

			switch command {
			case "help":
				showHelp()
			case "stats":
				showStats()
			case "clear":
				clearScreen()
			case "stop", "exit", "quit":
				fmt.Println("Iniciando apagado seguro...")
				// Send SIGINT to our own process to trigger graceful shutdown in main.go
				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
				return // exit the REPL goroutine
			default:
				fmt.Printf("Comando desconocido: '%s'. Escribe 'help' para ver la lista de comandos.\n", command)
			}
			fmt.Print("> ")
		}

		if err := scanner.Err(); err != nil {
			logger.Error(fmt.Sprintf("Error leyendo consola: %v", err), "CLI")
		}
	}()
}

func showHelp() {
	fmt.Println("=== PancyBot CLI Comandos ===")
	fmt.Println("  help  - Muestra esta lista de comandos")
	fmt.Println("  stats - Muestra estadísticas del bot y uso de recursos")
	fmt.Println("  clear - Limpia la consola")
	fmt.Println("  stop  - Apaga el bot de forma segura (guarda cachés)")
	fmt.Println("=============================")
}

func showStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Println("=== Estadísticas del Sistema ===")
	fmt.Printf("  Goroutines  : %d\n", runtime.NumGoroutine())
	fmt.Printf("  Memoria OS  : %.2f MB\n", float64(m.Sys)/1024/1024)
	fmt.Printf("  Memoria Uso : %.2f MB\n", float64(m.Alloc)/1024/1024)
	fmt.Printf("  GC Pausas   : %d\n", m.NumGC)
	fmt.Println("================================")
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
