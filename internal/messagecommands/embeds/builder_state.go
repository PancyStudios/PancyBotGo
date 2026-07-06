package embeds

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	builderStateMap = make(map[string]*discordgo.MessageEmbed)
	builderMutex    sync.RWMutex
)

func getBuilderState(userID string) *discordgo.MessageEmbed {
	builderMutex.RLock()
	defer builderMutex.RUnlock()

	if state, exists := builderStateMap[userID]; exists {
		return state
	}

	// Default state
	newState := &discordgo.MessageEmbed{
		Title:       "Mi Nuevo Embed",
		Description: "📝 | Descripción vacía.",
		Color:       0x3498DB,
	}

	// We don't save it here yet, let the caller save it to avoid map pollution
	return newState
}

func saveBuilderState(userID string, embed *discordgo.MessageEmbed) {
	builderMutex.Lock()
	defer builderMutex.Unlock()
	builderStateMap[userID] = embed
}

func clearBuilderState(userID string) {
	builderMutex.Lock()
	defer builderMutex.Unlock()
	delete(builderStateMap, userID)
}
