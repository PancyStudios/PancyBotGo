package discord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

// TestReplyEphemeralEmbedExists verifies that the ReplyEphemeralEmbed method exists
// and has the correct signature (compile-time check)
func TestReplyEphemeralEmbedExists(t *testing.T) {
	// This test verifies that ReplyEphemeralEmbed method exists and has the correct signature
	// by checking that we can reference the method
	
	// Create a type that matches the expected method signature
	type replyEphemeralEmbedFunc func(*CommandContext, *discordgo.MessageEmbed) error
	
	// Verify the method exists by assigning it to a variable
	var _ replyEphemeralEmbedFunc = (*CommandContext).ReplyEphemeralEmbed
	
	// If the above line compiles, the method exists with the correct signature
	t.Log("✅ ReplyEphemeralEmbed method exists with correct signature: func(*CommandContext, *discordgo.MessageEmbed) error")
}

// TestCommandCreation verifies that commands can be created with the builder pattern
func TestCommandCreation(t *testing.T) {
	handler := func(ctx *CommandContext) error {
		return nil
	}

	cmd := NewCommand("test", "Test command", "test", handler)
	
	if cmd == nil {
		t.Fatal("NewCommand returned nil")
	}

	if cmd.Name != "test" {
		t.Errorf("Name = %v, want %v", cmd.Name, "test")
	}

	if cmd.Description != "Test command" {
		t.Errorf("Description = %v, want %v", cmd.Description, "Test command")
	}

	if cmd.Category != "test" {
		t.Errorf("Category = %v, want %v", cmd.Category, "test")
	}

	if cmd.Run == nil {
		t.Error("Run function is nil")
	}
}

// TestCommandWithOptions verifies the WithOptions builder method
func TestCommandWithOptions(t *testing.T) {
	handler := func(ctx *CommandContext) error {
		return nil
	}

	option := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "test-option",
		Description: "Test option",
		Required:    true,
	}

	cmd := NewCommand("test", "Test command", "test", handler).
		WithOptions(option)

	if cmd.Options == nil {
		t.Fatal("Options is nil")
	}

	if len(cmd.Options) != 1 {
		t.Fatalf("Options length = %v, want %v", len(cmd.Options), 1)
	}

	if cmd.Options[0].Name != "test-option" {
		t.Errorf("Option name = %v, want %v", cmd.Options[0].Name, "test-option")
	}
}

// TestCommandWithPermissions verifies the permission builder methods
func TestCommandWithPermissions(t *testing.T) {
	handler := func(ctx *CommandContext) error {
		return nil
	}

	cmd := NewCommand("test", "Test command", "test", handler).
		WithUserPermissions(discordgo.PermissionAdministrator).
		WithBotPermissions(discordgo.PermissionSendMessages)

	if cmd.UserPermissions != discordgo.PermissionAdministrator {
		t.Errorf("UserPermissions = %v, want %v", cmd.UserPermissions, discordgo.PermissionAdministrator)
	}

	if cmd.BotPermissions != discordgo.PermissionSendMessages {
		t.Errorf("BotPermissions = %v, want %v", cmd.BotPermissions, discordgo.PermissionSendMessages)
	}
}

// TestCommandAsDev verifies the AsDev builder method
func TestCommandAsDev(t *testing.T) {
	handler := func(ctx *CommandContext) error {
		return nil
	}

	cmd := NewCommand("test", "Test command", "test", handler).AsDev()

	if !cmd.IsDev {
		t.Error("IsDev should be true after calling AsDev()")
	}
}

// TestToApplicationCommand verifies conversion to Discord application command
func TestToApplicationCommand(t *testing.T) {
	handler := func(ctx *CommandContext) error {
		return nil
	}

	option := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "test-option",
		Description: "Test option",
		Required:    true,
	}

	cmd := NewCommand("test", "Test command", "test", handler).
		WithOptions(option)

	appCmd := cmd.ToApplicationCommand()

	if appCmd == nil {
		t.Fatal("ToApplicationCommand returned nil")
	}

	if appCmd.Name != "test" {
		t.Errorf("ApplicationCommand Name = %v, want %v", appCmd.Name, "test")
	}

	if appCmd.Description != "Test command" {
		t.Errorf("ApplicationCommand Description = %v, want %v", appCmd.Description, "Test command")
	}

	if len(appCmd.Options) != 1 {
		t.Fatalf("ApplicationCommand Options length = %v, want %v", len(appCmd.Options), 1)
	}
}
