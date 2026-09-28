package discord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestToApplicationCommandIncludesDefaultPermissions(t *testing.T) {
	command := NewCommand("ban", "ban a member", "mod", nil).
		WithUserPermissions(discordgo.PermissionBanMembers)

	applicationCommand := command.ToApplicationCommand()
	if applicationCommand.DefaultMemberPermissions == nil {
		t.Fatal("DefaultMemberPermissions was not set")
	}
	if *applicationCommand.DefaultMemberPermissions != discordgo.PermissionBanMembers {
		t.Fatalf("unexpected permissions: got %d", *applicationCommand.DefaultMemberPermissions)
	}
}
