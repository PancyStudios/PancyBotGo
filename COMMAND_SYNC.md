# Sincronización de Comandos de Discord

Este documento explica cómo eliminar comandos residuales de Discord (slash commands `/`) que ya no existen en tu código actual.

## Problema

Cuando migras código de TypeScript a Go, o cuando eliminas o renombras comandos, Discord puede mantener los comandos antiguos visibles en la UI. Esto sucede porque Discord cachea los comandos registrados y no los elimina automáticamente cuando cambias tu código.

## Solución

Hemos creado una utilidad de sincronización que:
1. Lista todos los comandos registrados en Discord
2. Elimina comandos obsoletos
3. Registra solo los comandos actuales definidos en tu código

## Uso

### Opción 1: Usar la Utilidad de Sincronización (Recomendado)

La utilidad `sync-commands` te permite gestionar los comandos de Discord sin ejecutar el bot completo.

#### Listar comandos actuales

Para ver todos los comandos registrados en Discord:

```bash
go run cmd/sync-commands/main.go -list
```

Para ver comandos de un servidor específico:

```bash
go run cmd/sync-commands/main.go -list -guild TU_GUILD_ID
```

#### Sincronizar comandos

Para eliminar comandos obsoletos y registrar solo los actuales:

```bash
go run cmd/sync-commands/main.go -sync
```

O simplemente (sync es el comportamiento por defecto):

```bash
go run cmd/sync-commands/main.go
```

Para sincronizar comandos de un servidor específico:

```bash
go run cmd/sync-commands/main.go -sync -guild TU_GUILD_ID
```

#### Limpiar todos los comandos

Para eliminar TODOS los comandos sin registrar los nuevos:

```bash
go run cmd/sync-commands/main.go -clean
```

⚠️ **Advertencia**: Esto eliminará todos los comandos. Tendrás que ejecutar el bot o usar `-sync` después para registrarlos de nuevo.

### Opción 2: Usar la API del Bot Directamente

Si prefieres gestionar comandos desde tu propio código, puedes usar los métodos del `CommandHandler`:

```go
// En tu código Go
client := discord.Get()

// Listar comandos globales
commands, err := client.CommandHandler.ListGlobalCommands()

// Listar comandos de un guild específico
guildCommands, err := client.CommandHandler.ListGuildCommands("GUILD_ID")

// Eliminar todos los comandos globales
err = client.CommandHandler.UnregisterCommands()

// Eliminar comandos de un guild específico
err = client.CommandHandler.UnregisterGuildCommands("GUILD_ID")

// Sincronizar comandos (eliminar obsoletos y registrar actuales)
err = client.CommandHandler.SyncCommands()
```

### Opción 3: Habilitar Sincronización Automática al Inicio

Puedes modificar `cmd/bot/main.go` para sincronizar comandos automáticamente cada vez que inicies el bot. 

**Nota**: Esto es útil durante el desarrollo, pero NO se recomienda en producción ya que los comandos globales pueden tardar hasta 1 hora en propagarse.

## Diferencias entre Comandos Globales y de Guild

### Comandos Globales
- Disponibles en todos los servidores donde está el bot
- Tardan hasta 1 hora en propagarse
- Usa sin el flag `-guild`

### Comandos de Guild (Servidor)
- Solo disponibles en un servidor específico
- Se propagan instantáneamente (útil para desarrollo)
- Usa con el flag `-guild TU_GUILD_ID`

## Mejores Prácticas

1. **Durante Desarrollo**: Usa comandos de guild para testeo rápido
   ```bash
   go run cmd/sync-commands/main.go -sync -guild TU_GUILD_ID
   ```

2. **Antes de Producción**: Sincroniza comandos globales
   ```bash
   go run cmd/sync-commands/main.go -sync
   ```

3. **Después de Eliminar Comandos**: Ejecuta sync para limpiar comandos obsoletos
   ```bash
   go run cmd/sync-commands/main.go -sync
   ```

4. **Migración desde TypeScript**: Ejecuta clean primero, luego sync
   ```bash
   go run cmd/sync-commands/main.go -clean
   go run cmd/sync-commands/main.go -sync
   ```

## Solución de Problemas

### Los comandos antiguos siguen apareciendo

- **Comandos Globales**: Pueden tardar hasta 1 hora en actualizarse. Ten paciencia.
- **Comandos de Guild**: Deberían actualizarse inmediatamente. Intenta:
  1. Cerrar y reabrir Discord
  2. Verificar que estás usando el guild ID correcto

### Error de autenticación

Asegúrate de que tu archivo `.env` tiene el token correcto:
```env
BOT_TOKEN=tu_token_aqui
```

### No se encuentran comandos

Ejecuta `-list` primero para ver qué comandos están registrados actualmente en Discord:
```bash
go run cmd/sync-commands/main.go -list
```

## Comandos Disponibles Actualmente

Los comandos se definen en:
- `internal/commands/register.go` - Registro principal
- `internal/commands/mod/register.go` - Comandos de moderación
- `internal/commands/util.go` - Comandos de utilidad
- `internal/commands/music.go` - Comandos de música

Comandos actuales:
- `/mod ban` - Banear usuario
- `/mod kick` - Expulsar usuario
- `/mod warn` - Advertir usuario
- `/mod mute` - Silenciar usuario

## Métodos Añadidos al SDK

### En `pkg/discord/command.go`:

- `ReplyEphemeralEmbed(embed)` - Enviar respuesta embed efímera (solo visible para el usuario)

### En `pkg/discord/command_handler.go`:

- `ListGlobalCommands()` - Lista comandos globales
- `ListGuildCommands(guildID)` - Lista comandos de un servidor
- `UnregisterGuildCommands(guildID)` - Elimina comandos de un servidor
- `SyncCommands()` - Sincroniza comandos (elimina obsoletos, registra actuales)

## Soporte

Si tienes problemas, revisa los logs del bot. El logger mostrará errores específicos sobre la sincronización de comandos.
