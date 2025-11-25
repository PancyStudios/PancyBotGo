# 🎯 Sistema de Eventos Modular

## 📁 Estructura

El sistema de eventos está organizado en archivos separados por categoría, igual que los comandos:

```
internal/events/
├── register.go      # Registro central de todos los eventos
├── ready.go         # Evento cuando el bot se conecta
├── guild.go         # Eventos de servidores (join/leave)
├── member.go        # Eventos de miembros (join/leave/update)
├── message.go       # Eventos de mensajes (create/update/delete)
└── voice.go         # Eventos de voz (join/leave/move)
```

## 🚀 Cómo Funciona

### Registro Automático

En `main.go`, los eventos se registran automáticamente:

```go
// Register commands using the new commands package
commands.RegisterAll(discordClient)

// Register events using the new events package
events.RegisterAll(discordClient)  // ← Aquí se registran todos los eventos

// Start the bot
discordClient.Start()
```

### Organización por Archivos

Cada categoría de eventos tiene su propio archivo:

#### 1. **ready.go** - Bot Conectado
```go
func onReady(s *discordgo.Session, r *discordgo.Ready) {
    logger.Success("✅ Bot conectado!", "Ready")
    // Tu código aquí
}
```

#### 2. **guild.go** - Eventos de Servidor
```go
func onGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
    logger.Info("➕ Nuevo servidor: " + g.Name, "Guild")
    // Enviar mensaje de bienvenida
}

func onGuildDelete(s *discordgo.Session, g *discordgo.GuildDelete) {
    logger.Info("➖ Servidor removido", "Guild")
}
```

#### 3. **member.go** - Eventos de Miembros
```go
func onGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
    logger.Info("👋 Nuevo miembro: " + m.User.Username, "Member")
    // Enviar mensaje de bienvenida
}

func onGuildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
    logger.Info("👋 Miembro salió: " + m.User.Username, "Member")
}

func onGuildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
    // Detectar cambios en roles, nickname, etc.
}
```

#### 4. **message.go** - Eventos de Mensajes
```go
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.Bot {
        return
    }
    // Responder a menciones, palabras clave, etc.
}

func onMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
    // Detectar mensajes editados
}

func onMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
    // Detectar mensajes eliminados
}
```

#### 5. **voice.go** - Eventos de Voz
```go
func onVoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
    // Detectar usuarios uniéndose/saliendo de canales de voz
}
```

## ➕ Agregar Nuevos Eventos

**Paso 1: Crear un Nuevo Archivo**

Por ejemplo, para eventos de moderación, crea `moderation.go`:

```go
// filepath: internal/events/moderation.go
package events

import (
    "github.com/PancyStudios/PancyBotGo/pkg/discord"
    "github.com/PancyStudios/PancyBotGo/pkg/logger"
    "github.com/bwmarrin/discordgo"
)

// RegisterModerationEvents registers moderation-related event handlers
func RegisterModerationEvents(client *discord.ExtendedClient) {
    client.Session.AddHandler(onGuildBanAdd)
    client.Session.AddHandler(onGuildBanRemove)
}

// onGuildBanAdd is called when a user is banned
func onGuildBanAdd(s *discordgo.Session, b *discordgo.GuildBanAdd) {
    logger.Info("🔨 Usuario baneado: " + b.User.Username, "Moderation")
}

// onGuildBanRemove is called when a ban is removed
func onGuildBanRemove(s *discordgo.Session, b *discordgo.GuildBanRemove) {
    logger.Info("✅ Ban removido: " + b.User.Username, "Moderation")
}
```

### Paso 2: Registrar en register.go

Agrega tu nuevo archivo al registro:

```go
// filepath: internal/events/register.go
func RegisterAll(client *discord.ExtendedClient) {
    logger.System("📋 Registrando eventos del bot...", "Events")

    RegisterReadyEvent(client)
    RegisterGuildEvents(client)
    RegisterMemberEvents(client)
    RegisterMessageEvents(client)
    RegisterVoiceEvents(client)
    
    // ← Agregar aquí
    RegisterModerationEvents(client)

    logger.Success("✅ Todos los eventos registrados", "Events")
}
```

## 📝 Ejemplos de Uso

### Sistema de Bienvenida

Edita `member.go`:

```go
func onGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
    guild, _ := s.Guild(m.GuildID)
    
    // Mensaje público
    if guild.SystemChannelID != "" {
        embed := &discordgo.MessageEmbed{
            Title: "¡Bienvenido/a! 🎉",
            Description: fmt.Sprintf("<@%s> se unió al servidor", m.User.ID),
            Color: 0x00ff00,
        }
        s.ChannelMessageSendEmbed(guild.SystemChannelID, embed)
    }
    
    // DM privado
    channel, _ := s.UserChannelCreate(m.User.ID)
    s.ChannelMessageSend(channel.ID, "¡Bienvenido/a!")
}
```

### Respuestas Automáticas

Edita `message.go`:

```go
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.Bot {
        return
    }
    
    content := strings.ToLower(m.Content)
    
    if strings.Contains(content, "hola bot") {
        s.ChannelMessageSend(m.ChannelID, "¡Hola! 👋")
    }
    
    if strings.Contains(content, "ayuda") {
        s.ChannelMessageSend(m.ChannelID, "Usa /help para ver los comandos")
    }
}
```

### Sistema de Logs de Auditoría

Crea un nuevo archivo `audit.go`:

```go
package events

import (
    "fmt"
    "time"
    "github.com/PancyStudios/PancyBotGo/pkg/discord"
    "github.com/bwmarrin/discordgo"
)

// Canal de logs (configurable)
const LogChannelID = "TU_CANAL_ID_AQUI"

func RegisterAuditEvents(client *discord.ExtendedClient) {
    client.Session.AddHandler(logMemberJoin)
    client.Session.AddHandler(logMemberLeave)
    client.Session.AddHandler(logMessageDelete)
}

func logMemberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
    embed := &discordgo.MessageEmbed{
        Title: "👤 Nuevo Miembro",
        Description: fmt.Sprintf("<@%s> se unió", m.User.ID),
        Color: 0x00ff00,
        Timestamp: time.Now().Format(time.RFC3339),
    }
    s.ChannelMessageSendEmbed(LogChannelID, embed)
}

func logMemberLeave(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
    embed := &discordgo.MessageEmbed{
        Title: "👋 Miembro Salió",
        Description: m.User.Username,
        Color: 0xff0000,
        Timestamp: time.Now().Format(time.RFC3339),
    }
    s.ChannelMessageSendEmbed(LogChannelID, embed)
}

func logMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
    embed := &discordgo.MessageEmbed{
        Title: "🗑️ Mensaje Eliminado",
        Description: "Canal: <#" + m.ChannelID + ">",
        Color: 0xffa500,
        Timestamp: time.Now().Format(time.RFC3339),
    }
    s.ChannelMessageSendEmbed(LogChannelID, embed)
}
```

## 🎨 Eventos Disponibles

| Categoría | Evento | Archivo | Descripción |
|-----------|--------|---------|-------------|
| **Bot** | Ready | `ready.go` | Bot conectado |
| **Guild** | GuildCreate | `guild.go` | Bot se une a servidor |
| | GuildDelete | `guild.go` | Bot es removido |
| **Member** | GuildMemberAdd | `member.go` | Nuevo miembro |
| | GuildMemberRemove | `member.go` | Miembro sale |
| | GuildMemberUpdate | `member.go` | Miembro actualizado |
| **Message** | MessageCreate | `message.go` | Nuevo mensaje |
| | MessageUpdate | `message.go` | Mensaje editado |
| | MessageDelete | `message.go` | Mensaje eliminado |
| **Voice** | VoiceStateUpdate | `voice.go` | Estado de voz cambia |

### Otros Eventos Disponibles

Puedes agregar más eventos usando `AddHandler`:

```go
// En cualquier archivo de eventos
client.Session.AddHandler(func(s *discordgo.Session, r *discordgo.GuildRoleCreate) {
    logger.Info("Rol creado: " + r.Role.Name, "Events")
})

client.Session.AddHandler(func(s *discordgo.Session, b *discordgo.GuildBanAdd) {
    logger.Info("Usuario baneado: " + b.User.Username, "Events")
})

client.Session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
    logger.Debug("Reacción agregada", "Events")
})
```

## ✅ Ventajas de Este Sistema

1. **Organizado** - Cada categoría en su propio archivo
2. **Modular** - Fácil agregar/remover eventos
3. **Mantenible** - Código limpio y estructurado
4. **Escalable** - Crece con tu bot
5. **Similar a Comandos** - Misma estructura que el sistema de comandos

## 🔧 Personalización

### Deshabilitar Eventos

Comenta la línea en `register.go`:

```go
func RegisterAll(client *discord.ExtendedClient) {
    RegisterReadyEvent(client)
    RegisterGuildEvents(client)
    // RegisterMessageEvents(client)  // ← Comentar para deshabilitar
    RegisterVoiceEvents(client)
}
```

### Cambiar Comportamiento

Edita directamente el archivo correspondiente:

- Sistema de bienvenida → `member.go`
- Respuestas automáticas → `message.go`
- Logs de voz → `voice.go`

## 📚 Recursos

- **Eventos de Discord**: https://discord.com/developers/docs/topics/gateway-events
- **discordgo**: https://pkg.go.dev/github.com/bwmarrin/discordgo

---

**¡Listo!** Ahora tu bot tiene un sistema de eventos completamente modular y organizado. 🚀

