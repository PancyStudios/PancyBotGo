# Obtenemos la versión desde el último tag de git (ej: v1.0.5) o el hash corto
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || git rev-parse --short HEAD)
FECHA=$(date +%F)

printf "Building PancyBot version: %s\n" "$VERSION"

# Compilamos inyectando las variables
go build -ldflags "-X 'github.com/PancyStudios/PancyBotGo/pkg/config.Version=$VERSION' -X 'github.com/PancyStudios/PancyBotGo/pkg/config.BuildTime=$FECHA'" -o PancyBot.x86_64 cmd/bot/main.go
printf "Build completed: PancyBot.x86_64\n"