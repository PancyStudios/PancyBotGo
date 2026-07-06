package dev

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/PancyStudios/PancyBotGo/internal/messagecommands"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
)

func codegenCommand(ctx *messagecommands.MessageContext) error {
	if !isDev(ctx.Message.Author.ID) {
		_, err := ctx.ReplyError("Acceso Denegado", "Este comando es solo para la desarrolladora.")
		return err
	}

	if len(ctx.Args) == 0 {
		_, err := ctx.ReplyError("Uso Incorrecto", "Uso: `pan!codegen <user/guild> [duracion_dias] [cantidad]`")
		return err
	}

	codeTypeStr := strings.ToLower(ctx.Args[0])
	var premiumType models.PremiumCodeType
	if codeTypeStr == "user" {
		premiumType = models.PremiumCodeTypeUser
	} else if codeTypeStr == "guild" {
		premiumType = models.PremiumCodeTypeGuild
	} else {
		_, err := ctx.ReplyError("Error", "El tipo debe ser `user` o `guild`.")
		return err
	}

	duration := int64(0)
	if len(ctx.Args) > 1 {
		duration, _ = strconv.ParseInt(ctx.Args[1], 10, 64)
	}

	cantidad := 1
	if len(ctx.Args) > 2 {
		cantidad, _ = strconv.Atoi(ctx.Args[2])
	}

	permanent := false
	if duration == 0 {
		permanent = true
	}

	var generatedCodes []string

	for i := 0; i < cantidad; i++ {
		code := generateRandomCode()

		_, err := database.CreatePremiumCode(
			code,
			premiumType,
			int(duration),
			permanent,
			ctx.Message.Author.ID,
		)

		if err == nil {
			generatedCodes = append(generatedCodes, code)
		}
	}

	_, err := ctx.ReplySuccess("Códigos Generados", fmt.Sprintf("✅ Se generaron %d códigos.\n\n```\n%s\n```", len(generatedCodes), strings.Join(generatedCodes, "\n")))
	return err
}

func generateRandomCode() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
