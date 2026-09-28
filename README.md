# PancyBot Go

Esta es la versión en Go de los sistemas esenciales de PancyBot. Incluye todas las funcionalidades principales del bot de Discord, reescritas desde TypeScript a Go.

## Estructura del Proyecto

```
PancyBotGo/
├── cmd/
│   └── bot/
│       └── main.go              # Punto de entrada principal
├── internal/
│   └── commands/                # 📁 AQUÍ VAN TUS COMANDOS
│       ├── register.go          # Registro de todos los comandos
│       ├── utils/               # Comandos de utilidad (ping, status, etc.)
│       ├── music.go             # Comandos de música (play, pause, etc.)
│       └── mod/                 # 📁 Grupo de subcomandos /mod
│           ├── register.go      # Registra el grupo /mod
│           ├── ban.go           # /mod ban
│           ├── kick.go          # /mod kick
│           ├── warn.go          # /mod warn
│           └── mute.go          # /mod mute
├── pkg/
│   ├── config/                  # Gestión de configuración
│   ├── logger/                  # Sistema de logs con colores y webhooks
│   ├── database/                # Conexión MongoDB con DataManager y caché
│   ├── mqtt/                    # Comunicación MQTT para mensajería
│   ├── discord/                 # Cliente Discord con manejadores de comandos y eventos
│   ├── lavalink/                # 🎵 Cliente Lavalink para música
│   ├── web/                     # Servidor web HTTP con Gin
│   └── errors/                  # Manejo de errores y recuperación
└── go.mod                       # Módulo Go y dependencias
```

## 🚀 Cómo Añadir un Nuevo Comando

### Opción 1: Añadir a un archivo existente

Si tu comando pertenece a una categoría existente (util, music), añádelo al archivo correspondiente:

```go
// En internal/commands/utils/register.go

// Mi nuevo comando
miComandoCmd := discord.NewCommand(
    "micomando",                    // Nombre del comando
    "Descripción de mi comando",    // Descripción
    "util",                         // Categoría
    func(ctx *discord.CommandContext) error {
        // Tu lógica aquí
        return ctx.Reply("¡Hola desde mi comando!")
    },
)
client.CommandHandler.RegisterCommand(miComandoCmd)
client.CommandHandler.AddGlobalCommand(miComandoCmd.ToApplicationCommand())
```

### Opción 2: Crear una nueva categoría

1. **Crear archivo** en `internal/commands/`:

```go
// internal/commands/fun.go
package commands

import (
    "github.com/PancyStudios/PancyBotGo/pkg/discord"
)

func RegisterFunCommands(client *discord.ExtendedClient) {
    // Comando meme
    memeCmd := discord.NewCommand(
        "meme",
        "Muestra un meme random",
        "fun",
        func(ctx *discord.CommandContext) error {
            return ctx.Reply("🎭 Aquí va tu meme!")
        },
    )
    client.CommandHandler.RegisterCommand(memeCmd)
    client.CommandHandler.AddGlobalCommand(memeCmd.ToApplicationCommand())
}
```

2. **Registrar en `register.go`**:

```go
// internal/commands/register.go
func RegisterAll(client *discord.ExtendedClient) {
    RegisterUtilCommands(client)
    RegisterMusicCommands(client)
    RegisterFunCommands(client)  // ← Añadir esta línea
}
```

### Comandos con Opciones

```go
import "github.com/bwmarrin/discordgo"

cmd := discord.NewCommand(
    "saludar",
    "Saluda a alguien",
    "fun",
    func(ctx *discord.CommandContext) error {
        usuario := ctx.GetUserOption("usuario")
        mensaje := ctx.GetStringOption("mensaje")
        return ctx.Reply(fmt.Sprintf("¡Hola %s! %s", usuario.Username, mensaje))
    },
).WithOptions(
    &discordgo.ApplicationCommandOption{
        Type:        discordgo.ApplicationCommandOptionUser,
        Name:        "usuario",
        Description: "Usuario a saludar",
        Required:    true,
    },
    &discordgo.ApplicationCommandOption{
        Type:        discordgo.ApplicationCommandOptionString,
        Name:        "mensaje",
        Description: "Mensaje adicional",
        Required:    false,
    },
)
```

### Opción 3: Grupos de Subcomandos (`/mod ban`, `/mod kick`, etc.)

Para crear comandos como `/mod ban`, `/mod kick`, etc., usa grupos de subcomandos.

**Estructura de archivos (un comando por archivo):**
```
internal/commands/mod/
├── register.go    # Registra el grupo /mod
├── ban.go         # /mod ban
├── kick.go        # /mod kick
├── warn.go        # /mod warn
└── mute.go        # /mod mute
```

**1. Crear `internal/commands/mod/register.go`:**
```go
package mod

import "github.com/PancyStudios/PancyBotGo/pkg/discord"

func RegisterModCommands(client *discord.ExtendedClient) {
    // Crear subcomandos (cada uno puede estar en su propio archivo)
    banCmd := createBanCommand()
    kickCmd := createKickCommand()
    warnCmd := createWarnCommand()

    // Construir el grupo /mod con todos los subcomandos
    modGroup := client.CommandHandler.BuildCommandGroup(
        "mod",                      // Nombre del grupo
        "Comandos de moderación",   // Descripción
        banCmd,                     // Subcomandos...
        kickCmd,
        warnCmd,
    )

    // Registrar el grupo
    client.CommandHandler.AddGlobalCommand(modGroup)
}
```

**2. Crear `internal/commands/mod/ban.go` (un comando por archivo):**
```go
package mod

import (
    "fmt"
    "github.com/PancyStudios/PancyBotGo/pkg/discord"
    "github.com/bwmarrin/discordgo"
)

func createBanCommand() *discord.Command {
    return discord.NewCommand(
        "ban",
        "Banea a un usuario",
        "mod",
        banHandler,
    ).WithOptions(
        &discordgo.ApplicationCommandOption{
            Type:        discordgo.ApplicationCommandOptionUser,
            Name:        "usuario",
            Description: "Usuario a banear",
            Required:    true,
        },
    ).WithUserPermissions(discordgo.PermissionBanMembers)
}

func banHandler(ctx *discord.CommandContext) error {
    user := ctx.GetUserOption("usuario")
    // Lógica del ban...
    return ctx.Reply(fmt.Sprintf("🔨 %s ha sido baneado", user.Username))
}
```

**3. Registrar en `internal/commands/register.go`:**
```go
import "github.com/PancyStudios/PancyBotGo/internal/commands/mod"

func RegisterAll(client *discord.ExtendedClient) {
    RegisterUtilCommands(client)
    RegisterMusicCommands(client)
    mod.RegisterModCommands(client)  // ← Añadir esta línea
}
```

## Sistemas Convertidos

### 1. Sistema de Configuración (`pkg/config/`)
- Carga de variables de entorno desde `.env`
- Configuración centralizada para todos los servicios
- Valores por defecto para desarrollo

### 2. Sistema de Logging (`pkg/logger/`)
- Logging con colores para la consola
- Múltiples niveles: Critical, Error, Warn, Success, Info, Debug, System
- Integración con webhooks de Discord
- Logs a archivos con rotación

### 3. Base de Datos (`pkg/database/`)
- Conexión a MongoDB con reconexión automática
- DataManager genérico con caché LRU
- Cola de operaciones offline para sincronización

### 4. Comunicación MQTT (`pkg/mqtt/`)
- Cliente MQTT con publicación/suscripción
- Sistema de request/response con correlationId
- Soporte para wildcards en topics

### 5. Cliente Discord (`pkg/discord/`)
- Wrapper sobre discordgo
- Manejador de comandos slash
- Manejador de eventos
- Contexto de comandos enriquecido

### 6. Servidor Web (`pkg/web/`)
- Servidor HTTP basado en Gin
- Rate limiting integrado
- Logging de requests a webhooks
- Rutas API para status y salud

### 7. Manejo de Errores (`pkg/errors/`)
- Contador de errores con auto-shutdown
- Reporte a webhooks
- Recuperación de panics

### 8. 🎵 Sistema de Música Lavalink (`pkg/lavalink/`)
- Conexión a nodos Lavalink con reconexión automática
- Búsqueda de canciones (Deezer, YouTube, SoundCloud)
- Cola de reproducción con gestión completa
- Publicación de eventos via MQTT
- Comandos: play, pause, skip, stop, queue, volume, nowplaying

## Dependencias

- **discordgo**: Cliente Discord para Go
- **mongo-driver**: Driver oficial de MongoDB
- **paho.mqtt.golang**: Cliente MQTT
- **gin-gonic/gin**: Framework web HTTP
- **logrus**: Logging estructurado
- **godotenv**: Carga de archivos .env
- **gorilla/websocket**: WebSocket para Lavalink

## Requisitos

- Go 1.21+
- MongoDB
- Broker MQTT (opcional)
- Servidor Lavalink (para música)
- Token de bot de Discord

## Instalación

```bash
# Clonar el repositorio
git clone https://github.com/PancyStudios/PancyBotGo.git
cd PancyBotGo

# Verificar dependencias
go mod verify

# Compilar
go build -o pancybot ./cmd/bot
```

## Configuración

Crear un archivo `.env` en el directorio raíz con las siguientes variables:

```env
# Discord
BOT_TOKEN=tu_token_de_discord
DEV_GUILD_ID=id_del_servidor_de_desarrollo
DEVELOPER_ID=id_del_usuario_desarrollador

# MongoDB
MONGODB_URL=mongodb://localhost:27017
DB_NAME=PancyBot

# MQTT
MQTT_Host=localhost
MQTT_Port=1883
MQTT_User=
MQTT_Password=
# Las credenciales MQTT son obligatorias en producción; limita además los topics con ACLs en el broker.

# Lavalink (para música)
linkserver=localhost
linkpassword=youshallnotpass

# Web Server
PORT=3000
ADMIN_USER=admin
ADMIN_PASSWORD=una_contraseña_segura  # obligatorio en producción

# Environment
BOT_ENV=dev  # o 'prod' para producción

# Webhooks (opcional)
errorWebhook=url_webhook_errores
logsWebhook=url_webhook_logs
logsWebServerWebhook=url_webhook_web
```

## Ejecución

```bash
# Ejecutar directamente
go run ./cmd/bot

# O ejecutar el binario compilado
./pancybot
```

## Tests

```bash
# Ejecutar todos los tests
go test ./...

# Con verbose
go test -v ./...

# Coverage
go test -cover ./...
```

## Diferencias con la Versión TypeScript

| Característica | TypeScript | Go |
|----------------|------------|-----|
| Carga de comandos | Dinámica desde archivos | Registro programático |
| Genéricos | Parcial | Completo con Go 1.18+ |
| Concurrencia | Async/await | Goroutines y channels |
| Tipado | Estático (compilación) | Estático (compilación) |
| Performance | V8 JIT | Compilado nativamente |

## Ejemplo: Registrar un Evento

```go
// Registrar evento de mensaje
client.EventHandler.OnMessageCreate(func(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.Bot {
        return
    }
    logger.Info("Mensaje recibido: " + m.Content, "Messages")
})
```

## Licencia

MIT - Ver [LICENSE](../LICENCE) para más detalles.

## Autor

Desarrollado por PancyStudios / ImTurbis
