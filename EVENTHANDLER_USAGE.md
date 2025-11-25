# Guía de Uso del Event Handler

## ¿Qué es el Event Handler?

El Event Handler es un sistema que te permite registrar funciones que se ejecutan cuando ocurren eventos específicos en Discord (como cuando alguien se une al servidor, envía un mensaje, cambia su estado de voz, etc.).

## Estructura Básica

El Event Handler se inicializa automáticamente cuando creas el cliente de Discord y proporciona métodos helper para registrar eventos comunes.

## Eventos Disponibles

### 1. **OnReady** - Bot Listo
Se ejecuta cuando el bot se conecta exitosamente a Discord.

```go
client.EventHandler.OnReady(func(s *discordgo.Session, r *discordgo.Ready) {
    logger.Info(fmt.Sprintf("Bot conectado como: %s", r.User.Username), "Events")
    logger.Info(fmt.Sprintf("Conectado a %d servidores", len(r.Guilds)), "Events")
})
```

### 2. **OnGuildCreate** - Bot se une a un servidor
Se ejecuta cuando el bot se une a un servidor nuevo.

```go
client.EventHandler.OnGuildCreate(func(s *discordgo.Session, g *discordgo.GuildCreate) {
    logger.Info(fmt.Sprintf("Bot agregado al servidor: %s (ID: %s)", g.Name, g.ID), "Events")
    
    // Ejemplo: Enviar mensaje de bienvenida
    if g.SystemChannelID != "" {
        s.ChannelMessageSend(g.SystemChannelID, "¡Hola! Gracias por agregarme al servidor 👋")
    }
})
```

### 3. **OnGuildDelete** - Bot es removido de un servidor
Se ejecuta cuando el bot es expulsado o el servidor es eliminado.

```go
client.EventHandler.OnGuildDelete(func(s *discordgo.Session, g *discordgo.GuildDelete) {
    logger.Info(fmt.Sprintf("Bot removido del servidor ID: %s", g.ID), "Events")
})
```

### 4. **OnMessageCreate** - Nuevo mensaje
Se ejecuta cuando se crea un mensaje en cualquier canal.

```go
client.EventHandler.OnMessageCreate(func(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Ignorar mensajes del bot
    if m.Author.Bot {
        return
    }
    
    logger.Debug(fmt.Sprintf("Mensaje de %s: %s", m.Author.Username, m.Content), "Events")
    
    // Ejemplo: Responder a menciones
    if len(m.Mentions) > 0 {
        for _, mention := range m.Mentions {
            if mention.ID == s.State.User.ID {
                s.ChannelMessageSend(m.ChannelID, "¡Me mencionaste! Usa comandos slash para interactuar conmigo.")
                break
            }
        }
    }
})
```

### 5. **OnMessageUpdate** - Mensaje editado
Se ejecuta cuando un mensaje es editado.

```go
client.EventHandler.OnMessageUpdate(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
    if m.Author != nil && !m.Author.Bot {
        logger.Debug(fmt.Sprintf("Mensaje editado por %s en canal %s", m.Author.Username, m.ChannelID), "Events")
    }
})
```

### 6. **OnMessageDelete** - Mensaje eliminado
Se ejecuta cuando un mensaje es eliminado.

```go
client.EventHandler.OnMessageDelete(func(s *discordgo.Session, m *discordgo.MessageDelete) {
    logger.Debug(fmt.Sprintf("Mensaje eliminado en canal %s", m.ChannelID), "Events")
})
```

### 7. **OnGuildMemberAdd** - Nuevo miembro
Se ejecuta cuando alguien se une al servidor.

```go
client.EventHandler.OnGuildMemberAdd(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
    logger.Info(fmt.Sprintf("Nuevo miembro: %s#%s se unió al servidor %s", 
        m.User.Username, m.User.Discriminator, m.GuildID), "Events")
    
    // Ejemplo: Enviar mensaje de bienvenida por DM
    channel, err := s.UserChannelCreate(m.User.ID)
    if err == nil {
        s.ChannelMessageSend(channel.ID, fmt.Sprintf("¡Bienvenido/a %s! 🎉", m.User.Username))
    }
    
    // Ejemplo: Asignar rol automático
    // roleID := "123456789012345678"
    // s.GuildMemberRoleAdd(m.GuildID, m.User.ID, roleID)
})
```

### 8. **OnGuildMemberRemove** - Miembro sale
Se ejecuta cuando alguien sale o es expulsado del servidor.

```go
client.EventHandler.OnGuildMemberRemove(func(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
    logger.Info(fmt.Sprintf("Miembro %s#%s salió del servidor %s", 
        m.User.Username, m.User.Discriminator, m.GuildID), "Events")
})
```

### 9. **OnGuildMemberUpdate** - Miembro actualizado
Se ejecuta cuando un miembro cambia (roles, nickname, etc).

```go
client.EventHandler.OnGuildMemberUpdate(func(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
    logger.Debug(fmt.Sprintf("Miembro %s actualizado en servidor %s", m.User.Username, m.GuildID), "Events")
})
```

### 10. **OnVoiceStateUpdate** - Estado de voz cambia
Se ejecuta cuando alguien se une/sale de un canal de voz o cambia su estado.

```go
client.EventHandler.OnVoiceStateUpdate(func(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
    // Usuario se unió a un canal
    if v.ChannelID != "" && v.BeforeUpdate != nil && v.BeforeUpdate.ChannelID == "" {
        channel, _ := s.Channel(v.ChannelID)
        logger.Debug(fmt.Sprintf("Usuario %s se unió al canal de voz: %s", v.UserID, channel.Name), "Events")
    }
    
    // Usuario salió de un canal
    if v.ChannelID == "" && v.BeforeUpdate != nil && v.BeforeUpdate.ChannelID != "" {
        logger.Debug(fmt.Sprintf("Usuario %s salió del canal de voz", v.UserID), "Events")
    }
    
    // Usuario cambió de canal
    if v.ChannelID != "" && v.BeforeUpdate != nil && v.BeforeUpdate.ChannelID != "" && v.ChannelID != v.BeforeUpdate.ChannelID {
        logger.Debug(fmt.Sprintf("Usuario %s cambió de canal de voz", v.UserID), "Events")
    }
})
```

### 11. **OnInteractionCreate** - Interacción creada
Se ejecuta cuando se crea una interacción (slash command, botón, menú, etc).

```go
client.EventHandler.OnInteractionCreate(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
    // Nota: El CommandHandler ya maneja las interacciones de comandos
    // Este handler es útil para botones, menús, modales, etc.
    
    if i.Type == discordgo.InteractionMessageComponent {
        logger.Debug(fmt.Sprintf("Componente interactuado: %s", i.MessageComponentData().CustomID), "Events")
    }
})
```

## Implementación Completa - Ejemplo

Aquí está cómo implementar un sistema completo de eventos en tu `main.go`:

```go
package main

import (
    "fmt"
    "github.com/PancyStudios/PancyBotGo/pkg/discord"
    "github.com/PancyStudios/PancyBotGo/pkg/logger"
    "github.com/bwmarrin/discordgo"
)

func setupEvents(client *discord.ExtendedClient) {
    // Evento: Bot listo
    client.EventHandler.OnReady(func(s *discordgo.Session, r *discordgo.Ready) {
        logger.Success(fmt.Sprintf("✅ Bot online: %s", r.User.Username), "Events")
        logger.Info(fmt.Sprintf("📊 Servidores: %d", len(r.Guilds)), "Events")
    })

    // Evento: Nuevo servidor
    client.EventHandler.OnGuildCreate(func(s *discordgo.Session, g *discordgo.GuildCreate) {
        logger.Info(fmt.Sprintf("➕ Servidor nuevo: %s", g.Name), "Events")
    })

    // Evento: Nuevo miembro
    client.EventHandler.OnGuildMemberAdd(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
        logger.Info(fmt.Sprintf("👋 Bienvenido: %s", m.User.Username), "Events")
        
        // Enviar mensaje de bienvenida
        guild, _ := s.Guild(m.GuildID)
        if guild.SystemChannelID != "" {
            s.ChannelMessageSend(guild.SystemChannelID, 
                fmt.Sprintf("¡Bienvenido/a <@%s> al servidor! 🎉", m.User.ID))
        }
    })

    // Evento: Estado de voz
    client.EventHandler.OnVoiceStateUpdate(func(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
        // Detectar cuando alguien se une a voz
        if v.ChannelID != "" && (v.BeforeUpdate == nil || v.BeforeUpdate.ChannelID == "") {
            logger.Debug(fmt.Sprintf("🎤 Usuario %s en canal de voz", v.UserID), "Events")
        }
    })

    // Evento: Mensajes (para respuestas automáticas)
    client.EventHandler.OnMessageCreate(func(s *discordgo.Session, m *discordgo.MessageCreate) {
        if m.Author.Bot {
            return
        }

        // Responder a palabras clave
        if strings.Contains(strings.ToLower(m.Content), "hola bot") {
            s.ChannelMessageSend(m.ChannelID, "¡Hola! 👋 Usa `/help` para ver mis comandos.")
        }
    })
}

func main() {
    // ... código de inicialización ...
    
    // Configurar eventos ANTES de iniciar el bot
    setupEvents(discordClient)
    
    // Iniciar bot
    if err := discordClient.Start(); err != nil {
        logger.Critical(fmt.Sprintf("Error: %v", err), "Main")
        os.Exit(1)
    }
    
    // ... resto del código ...
}
```

## Usando el RegisterEvent Directo

Si necesitas registrar un evento personalizado que no tiene un método helper:

```go
// Registrar evento personalizado
client.EventHandler.RegisterEvent(func(s *discordgo.Session, event *discordgo.TypingStart) {
    logger.Debug(fmt.Sprintf("Usuario %s está escribiendo...", event.UserID), "Events")
})
```

## Eventos Disponibles en discordgo

Además de los que tienen helpers, puedes registrar cualquier evento de discordgo:

- `*discordgo.ChannelCreate` - Canal creado
- `*discordgo.ChannelUpdate` - Canal actualizado
- `*discordgo.ChannelDelete` - Canal eliminado
- `*discordgo.GuildBanAdd` - Usuario baneado
- `*discordgo.GuildBanRemove` - Ban removido
- `*discordgo.GuildRoleCreate` - Rol creado
- `*discordgo.GuildRoleUpdate` - Rol actualizado
- `*discordgo.GuildRoleDelete` - Rol eliminado
- `*discordgo.MessageReactionAdd` - Reacción agregada
- `*discordgo.MessageReactionRemove` - Reacción removida
- `*discordgo.PresenceUpdate` - Estado/presencia actualizado
- `*discordgo.TypingStart` - Usuario escribiendo

## Mejores Prácticas

1. **Registra eventos ANTES de llamar a `Start()`**
   ```go
   setupEvents(discordClient)
   discordClient.Start()
   ```

2. **Ignora mensajes de bots en OnMessageCreate**
   ```go
   if m.Author.Bot {
       return
   }
   ```

3. **Maneja errores apropiadamente**
   ```go
   channel, err := s.Channel(channelID)
   if err != nil {
       logger.Error(fmt.Sprintf("Error: %v", err), "Events")
       return
   }
   ```

4. **Usa goroutines para operaciones largas**
   ```go
   client.EventHandler.OnGuildMemberAdd(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
       go func() {
           // Operación larga aquí
       }()
   })
   ```

5. **Logging apropiado**
   - `logger.Debug()` para eventos frecuentes
   - `logger.Info()` para eventos importantes
   - `logger.Error()` para errores

## Ejemplo Completo: Sistema de Logs

```go
func setupAuditLog(client *discord.ExtendedClient, logChannelID string) {
    // Logs de miembros
    client.EventHandler.OnGuildMemberAdd(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
        embed := &discordgo.MessageEmbed{
            Title:       "👋 Nuevo Miembro",
            Description: fmt.Sprintf("<@%s> se unió al servidor", m.User.ID),
            Color:       0x00ff00,
            Timestamp:   time.Now().Format(time.RFC3339),
        }
        s.ChannelMessageSendEmbed(logChannelID, embed)
    })

    client.EventHandler.OnGuildMemberRemove(func(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
        embed := &discordgo.MessageEmbed{
            Title:       "👋 Miembro Salió",
            Description: fmt.Sprintf("%s#%s salió del servidor", m.User.Username, m.User.Discriminator),
            Color:       0xff0000,
            Timestamp:   time.Now().Format(time.RFC3339),
        }
        s.ChannelMessageSendEmbed(logChannelID, embed)
    })

    // Logs de mensajes
    client.EventHandler.OnMessageDelete(func(s *discordgo.Session, m *discordgo.MessageDelete) {
        embed := &discordgo.MessageEmbed{
            Title: "🗑️ Mensaje Eliminado",
            Fields: []*discordgo.MessageEmbedField{
                {
                    Name:  "Canal",
                    Value: fmt.Sprintf("<#%s>", m.ChannelID),
                },
                {
                    Name:  "ID del Mensaje",
                    Value: m.ID,
                },
            },
            Color:     0xffa500,
            Timestamp: time.Now().Format(time.RFC3339),
        }
        s.ChannelMessageSendEmbed(logChannelID, embed)
    })
}
```

## Resumen

El Event Handler te permite:
- ✅ Responder a eventos de Discord
- ✅ Implementar funcionalidades automáticas
- ✅ Crear sistemas de logs y auditoría
- ✅ Mejorar la interactividad del bot
- ✅ Manejar estados de voz, mensajes, miembros, etc.

Para más información sobre eventos de Discord, consulta:
- [Documentación de discordgo](https://pkg.go.dev/github.com/bwmarrin/discordgo)
- [Discord API Documentation](https://discord.com/developers/docs/topics/gateway-events)

