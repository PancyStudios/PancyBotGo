# Sistema de Códigos Premium

## Descripción
Sistema completo de códigos premium para PancyBot Go que permite generar y canjear códigos premium tanto para usuarios como para servidores.

## Comandos

### `/premium redeem`
Canjea un código premium.

**Opciones:**
- `codigo` (requerido): El código premium a canjear
- `tipo` (opcional): Tipo de código (`user` o `guild`). Si no se especifica, se detecta automáticamente

**Uso:**
```
/premium redeem codigo:PANC-abc123def456
/premium redeem codigo:PANC-abc123def456 tipo:user
/premium redeem codigo:PANC-abc123def456 tipo:guild
```

**Permisos:**
- Para códigos de usuario: Cualquier usuario puede canjear
- Para códigos de servidor: Solo administradores o el dueño del servidor

### `/dev codegen`
Genera códigos premium.

**⚠️ IMPORTANTE:** Este comando solo está disponible en el servidor de desarrollo configurado en `DevGuildID`.

**Opciones:**
- `tipo` (requerido): Tipo de código a generar (`user` o `guild`)
- `duracion` (opcional): Duración en días (0 para permanente, defecto: 30)
- `cantidad` (opcional): Cantidad de códigos a generar (1-10, defecto: 1)

**Uso:**
```
/dev codegen tipo:user duracion:30 cantidad:5
/dev codegen tipo:guild duracion:365 cantidad:1
/dev codegen tipo:user duracion:0 cantidad:1    # Permanente
```

**Nota:** El comando se registra automáticamente solo en el servidor de desarrollo, por lo que no aparecerá en ningún otro servidor.

## Estructura de Base de Datos

### Colección: `premium`
Almacena premium de usuarios.

```go
{
    "user": "123456789",        // ID del usuario
    "permanent": false,          // Si es permanente
    "expira": 1234567890000     // Timestamp en milisegundos
}
```

### Colección: `premium_guilds`
Almacena premium de servidores.

```go
{
    "guild": "123456789",       // ID del servidor
    "permanent": false,          // Si es permanente
    "expira": 1234567890000     // Timestamp en milisegundos
}
```

### Colección: `premium_codes`
Almacena códigos premium generados.

```go
{
    "_id": "PANC-abc123def456",  // El código
    "type": "user",               // "user" o "guild"
    "duration_days": 30,          // Duración en días
    "permanent": false,           // Si es permanente
    "is_claimed": false,          // Si ya fue canjeado
    "created_at": "...",          // Fecha de creación
    "claimed_by": "123456789",    // ID de quien lo canjeó (opcional)
    "claimed_at": "...",          // Fecha de canje (opcional)
    "created_by": "123456789"     // ID de quien lo creó
}
```

## Servicios de Base de Datos

### Premium de Usuario
- `GetUserPremium(userID string)`: Obtiene el premium de un usuario
- `IsUserPremium(userID string)`: Verifica si un usuario tiene premium activo
- `GrantUserPremium(userID, duration, permanent)`: Otorga premium a un usuario
- `RemoveUserPremium(userID)`: Remueve el premium de un usuario

### Premium de Servidor
- `GetGuildPremium(guildID string)`: Obtiene el premium de un servidor
- `IsGuildPremium(guildID string)`: Verifica si un servidor tiene premium activo
- `GrantGuildPremium(guildID, duration, permanent)`: Otorga premium a un servidor
- `RemoveGuildPremium(guildID)`: Remueve el premium de un servidor

### Códigos Premium
- `CreatePremiumCode(code, type, durationDays, permanent, createdBy)`: Crea un nuevo código
- `GetPremiumCode(code)`: Obtiene información de un código
- `RedeemPremiumCode(code, claimedBy)`: Redime un código de usuario
- `RedeemPremiumCodeForGuild(code, guildID, claimedBy)`: Redime un código de servidor
- `GetAllPremiumCodes()`: Obtiene todos los códigos (admin)
- `DeletePremiumCode(code)`: Elimina un código

## Formato de Códigos

Los códigos se generan con el formato: `PANC-xxxxxxxxxxxx`

Donde `xxxxxxxxxxxx` son 12 caracteres hexadecimales aleatorios.

Ejemplo: `PANC-a1b2c3d4e5f6`

## Seguridad

- El comando `/dev codegen` solo está registrado y disponible en el servidor de desarrollo (configurado en `DevGuildID`)
- Los códigos de servidor solo pueden ser canjeados por administradores del servidor
- Los códigos solo pueden canjearse una vez
- Se registran todos los canjes en los logs
- Los códigos generados son criptográficamente seguros usando `crypto/rand`

## Ejemplos de Uso

### Generar código premium de usuario por 30 días
```
/dev codegen tipo:user duracion:30 cantidad:1
```

### Generar código premium de servidor permanente
```
/dev codegen tipo:guild duracion:0 cantidad:1
```

### Canjear código
```
/premium redeem codigo:PANC-a1b2c3d4e5f6
```

## Integración

Para verificar si un comando requiere premium, se puede usar en el middleware del command handler:

```go
// Verificar premium de usuario
isPremium, _, err := database.IsUserPremium(userID)
if err != nil {
    // Manejar error
}
if !isPremium {
    // Usuario no premium
}

// Verificar premium de servidor
isPremium, _, err := database.IsGuildPremium(guildID)
if err != nil {
    // Manejar error
}
if !isPremium {
    // Servidor no premium
}
```

## Notas

- Los premiums expirados se eliminan automáticamente al consultarlos
- El campo `expira` se almacena en milisegundos (Unix timestamp * 1000)
- Los códigos permanentes tienen `duration_days = 0` y `permanent = true`

