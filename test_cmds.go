//go:build tools

package main

import (
	"fmt"
	"github.com/PancyStudios/PancyBotGo/internal/commands"
	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/discord"
)

func main() {
	config.Load()
	c, _ := discord.Init("dummy")
	commands.RegisterAll(c)

	// Test if it's there
	cmd, ok := c.Commands.Get("help.cmds")
	if ok {
		fmt.Printf("FOUND help.cmds: %v\n", cmd != nil)
	} else {
		fmt.Printf("NOT FOUND!\n")
	}
}
