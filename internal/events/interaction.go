// Package events provides event handlers for interaction events
// This file demonstrates how to add custom interaction handlers for buttons, menus, etc.
package events

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

// RegisterInteractionEvents registers all interaction-related event handlers
// Uncomment this function in register.go to enable interaction events
func RegisterInteractionEvents(client *discord.ExtendedClient) {
	client.Session.AddHandler(onInteractionCreate)
}

// onInteractionCreate is called when an interaction is created (buttons, menus, modals, etc.)
// Note: Slash commands are already handled by the CommandHandler
func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Handle message components (buttons, select menus)
	if i.Type == discordgo.InteractionMessageComponent {
		customID := i.MessageComponentData().CustomID
		logger.Debug(fmt.Sprintf("🔘 Componente clickeado: %s", customID), "Interaction")

		// Handle different button/menu IDs
		switch customID {
		case "button_accept":
			handleAcceptButton(s, i)
		case "button_deny":
			handleDenyButton(s, i)
		case "menu_roles":
			handleRoleMenu(s, i)
		default:
			logger.Debug(fmt.Sprintf("Componente no manejado: %s", customID), "Interaction")
		}
		return
	}

	// Handle modal submits
	if i.Type == discordgo.InteractionModalSubmit {
		modalID := i.ModalSubmitData().CustomID
		logger.Debug(fmt.Sprintf("📝 Modal enviado: %s", modalID), "Interaction")

		switch modalID {
		case "modal_feedback":
			handleFeedbackModal(s, i)
		default:
			logger.Debug(fmt.Sprintf("Modal no manejado: %s", modalID), "Interaction")
		}
		return
	}
}

// Example button handlers

func handleAcceptButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ ¡Aceptado!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error respondiendo interacción: %v", err), "Interaction")
	}
}

func handleDenyButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ Denegado",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error respondiendo interacción: %v", err), "Interaction")
	}
}

func handleRoleMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()

	if len(data.Values) > 0 {
		selectedRole := data.Values[0]

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Has seleccionado: <@&%s>", selectedRole),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			logger.Error(fmt.Sprintf("Error respondiendo interacción: %v", err), "Interaction")
		}

		// Add role to user
		err = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, selectedRole)
		if err != nil {
			logger.Error(fmt.Sprintf("Error asignando rol: %v", err), "Interaction")
		}
	}
}

func handleFeedbackModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()

	// Get the feedback text from the modal
	feedback := ""
	for _, component := range data.Components {
		if actionRow, ok := component.(*discordgo.ActionsRow); ok {
			for _, c := range actionRow.Components {
				if textInput, ok := c.(*discordgo.TextInput); ok {
					feedback = textInput.Value
					break
				}
			}
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ ¡Gracias por tu feedback!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error respondiendo interacción: %v", err), "Interaction")
	}

	logger.Info(fmt.Sprintf("Feedback recibido de %s: %s", i.Member.User.Username, feedback), "Interaction")
}
