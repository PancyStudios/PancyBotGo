package config

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

// Register registers all configuration commands
func Register(client *discord.ExtendedClient) {
	suggestConfigCmd := createSuggestConfigCommand()
	confessConfigCmd := createConfessConfigCommand()
	verifyChannelCmd := createVerifyChannelCommand()
	verifyRoleCmd := createVerifyRoleCommand()
	sendVerifyCmd := createSendVerifyCommand()

	// Create the base /config command
	configCmd := discord.NewCommand(
		"config",
		"✨ | Comandos de configuración del servidor",
		"config",
		HandleSubcommand,
	).WithUserPermissions(discordgo.PermissionManageGuild) // Require ManageServer permission

	// Add subcommands
	configCmd.Options = append(configCmd.Options, welcomeSubcommand())
	configCmd.Options = append(configCmd.Options, farewellSubcommand())
	configCmd.Options = append(configCmd.Options, autoroleSubcommand())
	configCmd.Options = append(configCmd.Options, logsSubcommand())

	// Register the command with the client
	client.CommandHandler.RegisterCommand(configCmd)
	client.CommandHandler.AddGlobalCommand(configCmd.ToApplicationCommand())

	cmds := []*discord.Command{
		suggestConfigCmd,
		confessConfigCmd,
		verifyChannelCmd,
		verifyRoleCmd,
		sendVerifyCmd,
	}

	for _, cmd := range cmds {
		cmd.WithUserPermissions(discordgo.PermissionManageGuild)
		client.CommandHandler.RegisterCommand(cmd)
		client.CommandHandler.AddGlobalCommand(cmd.ToApplicationCommand())
	}
}

// helper to safely run subcommand handlers
func HandleSubcommand(ctx *discord.CommandContext) error {
	options := ctx.Interaction.ApplicationCommandData().Options
	if len(options) == 0 {
		return ctx.ReplyEphemeral("❌ Comando inválido.")
	}

	subcommand := options[0].Name
	switch subcommand {
	case "welcome":
		return handleWelcome(ctx, options[0].Options)
	case "farewell":
		return handleFarewell(ctx, options[0].Options)
	case "autorole":
		return handleAutorole(ctx, options[0].Options)
	case "logs":
		return handleLogs(ctx, options[0].Options)
	default:
		return ctx.ReplyEphemeral("❌ Subcomando no encontrado.")
	}
}
