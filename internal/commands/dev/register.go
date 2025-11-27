package dev

import (
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

// Register registers all dev commands as /dev subcommands (only in dev guild)
func Register(client *discord.ExtendedClient) {
	// Create individual subcommands
	codegenCmd := CreateCodeGenCommand()
	codelistCmd := CreateCodeListCommand()
	codedelCmd := CreateCodeDelCommand()
	evalCmd := CreateEvalCommand()

	// Create blacklist subcommands
	blacklistAddCmd := CreateBlacklistAddCommand()
	blacklistRemoveCmd := CreateBlacklistRemoveCommand()
	blacklistListCmd := CreateBlacklistListCommand()

	// Build the blacklist subcommand group
	blacklistGroup := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
		Name:        "blacklist",
		Description: "Gestión de blacklist",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        blacklistAddCmd.Name,
				Description: blacklistAddCmd.Description,
				Options:     blacklistAddCmd.Options,
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        blacklistRemoveCmd.Name,
				Description: blacklistRemoveCmd.Description,
				Options:     blacklistRemoveCmd.Options,
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        blacklistListCmd.Name,
				Description: blacklistListCmd.Description,
				Options:     blacklistListCmd.Options,
			},
		},
	}

	// Build the /dev command group with all subcommands
	devGroup := &discordgo.ApplicationCommand{
		Name:        "dev",
		Description: "Comandos de desarrollo",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        codegenCmd.Name,
				Description: codegenCmd.Description,
				Options:     codegenCmd.Options,
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        codelistCmd.Name,
				Description: codelistCmd.Description,
				Options:     codelistCmd.Options,
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        codedelCmd.Name,
				Description: codedelCmd.Description,
				Options:     codedelCmd.Options,
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        evalCmd.Name,
				Description: evalCmd.Description,
				Options:     evalCmd.Options,
			},
			blacklistGroup,
		},
	}

	// Register the individual commands in the command map
	client.Commands.Set("dev.codegen", codegenCmd)
	client.Commands.Set("dev.codelist", codelistCmd)
	client.Commands.Set("dev.codedel", codedelCmd)
	client.Commands.Set("dev.eval", evalCmd)
	client.Commands.Set("dev.blacklist.add", blacklistAddCmd)
	client.Commands.Set("dev.blacklist.remove", blacklistRemoveCmd)
	client.Commands.Set("dev.blacklist.list", blacklistListCmd)

	// Register the command group as dev-only command
	client.CommandHandler.AddDevCommand(devGroup)
}
