package economy

import (
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
)

// Register economy text commands
func Register() {
	messagecommands.RegisterCommand("eco", "Comandos de economía global", "pan!eco <comando>", "Economy", func(ctx *messagecommands.MessageContext) error { return ecoRouter(ctx, true) })
	messagecommands.RegisterCommand("ecog", "Comandos de economía local", "pan!ecog <comando>", "Economy", func(ctx *messagecommands.MessageContext) error { return ecoRouter(ctx, false) })

	messagecommands.RegisterCommand("shop", "Tienda global", "pan!shop <comando>", "Economy", func(ctx *messagecommands.MessageContext) error { return shopRouter(ctx, true) })
	messagecommands.RegisterCommand("shopg", "Tienda local", "pan!shopg <comando>", "Economy", func(ctx *messagecommands.MessageContext) error { return shopRouter(ctx, false) })
}

func ecoRouter(ctx *messagecommands.MessageContext, isGlobal bool) error {
	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Debes especificar un comando. Ejemplo: `work`, `balance`, `rob`, `deposit`, `withdraw`, `pay`, `daily`, `weekly`, `crime`, `slut`, `top`")
		return err
	}

	cmd := strings.ToLower(ctx.Args[0])
	ctx.Args = ctx.Args[1:] // Shift args

	switch cmd {
	case "work":
		return workCommand(ctx, isGlobal)
	case "balance":
		return balanceCommand(ctx, isGlobal)
	case "rob":
		return robCommand(ctx, isGlobal)
	case "deposit":
		return depositCommand(ctx, isGlobal)
	case "withdraw":
		return withdrawCommand(ctx, isGlobal)
	case "pay":
		return payCommand(ctx, isGlobal)
	case "daily":
		return dailyCommand(ctx, isGlobal)
	case "weekly":
		return weeklyCommand(ctx, isGlobal)
	case "crime":
		return crimeCommand(ctx, isGlobal)
	case "slut":
		return slutCommand(ctx, isGlobal)
	case "top":
		return topCommand(ctx, isGlobal)
	default:
		_, err := ctx.ReplyError("Comando no encontrado", "Ese comando de economía no existe.")
		return err
	}
}

func shopRouter(ctx *messagecommands.MessageContext, isGlobal bool) error {
	if len(ctx.Args) == 0 {
		return shopCommand(ctx, isGlobal) // view shop
	}

	cmd := strings.ToLower(ctx.Args[0])

	switch cmd {
	case "buy":
		ctx.Args = ctx.Args[1:]
		return buyCommand(ctx, isGlobal)
	case "use":
		ctx.Args = ctx.Args[1:]
		return useCommand(ctx, isGlobal)
	case "inventory", "inv":
		ctx.Args = ctx.Args[1:]
		return inventoryCommand(ctx, isGlobal)
	case "admin":
		ctx.Args = ctx.Args[1:]
		return adminShopCommand(ctx, isGlobal)
	default:
		// view shop with page
		return shopCommand(ctx, isGlobal)
	}
}
