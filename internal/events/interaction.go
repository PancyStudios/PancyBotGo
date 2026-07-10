// Package events provides event handlers for interaction events
// This file demonstrates how to add custom interaction handlers for buttons, menus, etc.
package events

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PancyStudios/PancyBotGo/internal/commands/embeds"
	slashHelpCommands "github.com/PancyStudios/PancyBotGo/internal/commands/help"
	"github.com/PancyStudios/PancyBotGo/internal/commands/economy"
	helpMsgCommands "github.com/PancyStudios/PancyBotGo/internal/messagecommands/help"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

// RegisterInteractionEvents registers all interaction-related event handlers
// Uncomment this function in register.go to enable interaction events
func RegisterInteractionEvents(client *discord.ExtendedClient) {
	client.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		onInteractionCreate(s, i, client)
	})
}

// onInteractionCreate is called when an interaction is created (buttons, menus, modals, etc.)
// Note: Slash commands are already handled by the CommandHandler
func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate, client *discord.ExtendedClient) {
	// Handle message components (buttons, select menus)
	if i.Type == discordgo.InteractionMessageComponent {
		customID := i.MessageComponentData().CustomID
		logger.Debug(fmt.Sprintf("🔘 Componente clickeado: %s", customID), "Interaction")

		if embeds.HandleInteraction(s, i) {
			return
		}

		if helpMsgCommands.HandleInteraction(s, i) {
			return
		}

		if slashHelpCommands.HandleSlashInteraction(s, i, client) {
			return
		}

		// Handle different button/menu IDs
		switch customID {
		case "button_accept":
			handleAcceptButton(s, i)
		case "button_deny":
			handleDenyButton(s, i)
		case "menu_roles":
			handleRoleMenu(s, i)
		case "btn_verify_user":
			handleVerifyUser(s, i)
		default:
			if len(customID) >= 9 && customID[:9] == "shop_nav_" {
				handleShopNavigation(s, i)
			} else {
				logger.Debug(fmt.Sprintf("Componente no manejado: %s", customID), "Interaction")
			}
		}
		return
	}

	// Handle modal submits
	if i.Type == discordgo.InteractionModalSubmit {
		modalID := i.ModalSubmitData().CustomID
		logger.Debug(fmt.Sprintf("📝 Modal enviado: %s", modalID), "Interaction")

		if embeds.HandleInteraction(s, i) {
			return
		}

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

func handleShopNavigation(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	
	if customID == "shop_nav_menu" {
		embed, components := economy.ShopMenu()
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{embed},
				Components: components,
			},
		})
		if err != nil {
			logger.Error(fmt.Sprintf("Error respondiendo interacción shop_nav_menu: %v", err), "Interaction")
		}
		return
	}

	// Format is shop_nav_{type}_{page}
	parts := strings.Split(customID, "_")
	if len(parts) != 4 {
		logger.Error(fmt.Sprintf("Invalid customID format: %s", customID), "Interaction")
		return
	}
	
	shopType := parts[2]
	page, err := strconv.Atoi(parts[3])
	if err != nil {
		logger.Error(fmt.Sprintf("Error parseando página %s: %v", customID, err), "Interaction")
		return
	}

	embed, components, err := economy.RenderShopPage(i.GuildID, shopType, page)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Ocurrió un error al cargar la tienda.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error actualizando mensaje de la tienda: %v", err), "Interaction")
	}
}

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

func handleVerifyUser(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildDoc, err := database.GlobalGuildDM.Get(bson.M{"id": i.GuildID})
	if err != nil || guildDoc == nil || !guildDoc.Protection.Verification.Enable || guildDoc.Protection.Verification.Role == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ El sistema de verificación no está activo o el rol no está configurado en este servidor.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	verifyRoleID := guildDoc.Protection.Verification.Role

	// Check Account Age if configured
	if guildDoc.Protection.Verification.MinAccountAgeDays > 0 {
		// Snowflake to Timestamp formula
		// Timestamp = (Snowflake >> 22) + DiscordEpoch (1420070400000)
		idInt, _ := strconv.ParseInt(i.Member.User.ID, 10, 64)
		timestampMs := (idInt >> 22) + 1420070400000
		createdAt := time.UnixMilli(timestampMs)
		accountAge := time.Since(createdAt)

		minAgeDuration := time.Duration(guildDoc.Protection.Verification.MinAccountAgeDays) * 24 * time.Hour
		if accountAge < minAgeDuration {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("❌ Tu cuenta es demasiado reciente. Debe tener al menos **%d días** de antigüedad para verificarte.", guildDoc.Protection.Verification.MinAccountAgeDays),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
	}

	// Check if user already has the role
	hasRole := false
	for _, r := range i.Member.Roles {
		if r == verifyRoleID {
			hasRole = true
			break
		}
	}

	if hasRole {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "✅ Ya estás verificado.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	err = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, verifyRoleID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error añadiendo rol de verificación: %v", err), "Interaction")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ No pude añadirte el rol. Es posible que me falten permisos o el rol esté por encima del mío.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🎉 ¡Te has verificado exitosamente!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
