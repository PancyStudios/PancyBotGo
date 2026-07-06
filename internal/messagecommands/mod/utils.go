package mod

import "github.com/bwmarrin/discordgo"

// getHighestRolePosition calcula la posición más alta de los roles de un miembro
func getHighestRolePosition(guild *discordgo.Guild, member *discordgo.Member) int {
	highest := 0
	roleMap := make(map[string]int)
	for _, r := range guild.Roles {
		roleMap[r.ID] = r.Position
	}

	for _, roleID := range member.Roles {
		if pos, ok := roleMap[roleID]; ok {
			if pos > highest {
				highest = pos
			}
		}
	}
	return highest
}
