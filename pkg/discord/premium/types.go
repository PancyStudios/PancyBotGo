package premium

import "github.com/PancyStudios/PancyBotGo/pkg/database"

// Requirement defines the premium requirements for a command.
type Requirement struct {
	User  bool
	Guild bool
}

// RequirementType enumerates possible requirement keys
const (
	requirementTypeNone         = "none"
	requirementTypeUserOnly     = "user"
	requirementTypeGuildOnly    = "guild"
	requirementTypeUserAndGuild = "both"
)

// IsNone returns true if no premium requirement is set.
func (r Requirement) IsNone() bool {
	return !r.User && !r.Guild
}

// Predefined premium requirement helpers.
var (
	RequirementNone         = Requirement{}
	RequirementUserOnly     = Requirement{User: true}
	RequirementGuildOnly    = Requirement{Guild: true}
	RequirementUserAndGuild = Requirement{User: true, Guild: true}
)

// Check verifies the premium status based on the requirement.
func Check(req Requirement, userID, guildID string) (bool, string) {
	if req.IsNone() {
		return true, ""
	}

	if req.User {
		ok, _, err := database.IsUserPremium(userID)
		if err != nil {
			return false, "Error al verificar premium de usuario."
		}
		if !ok {
			return false, "Necesitas premium de usuario para usar este comando."
		}
	}

	if req.Guild && guildID != "" {
		ok, _, err := database.IsGuildPremium(guildID)
		if err != nil {
			return false, "Error al verificar premium del servidor."
		}
		if !ok {
			return false, "Este servidor necesita premium para usar este comando."
		}
	} else if req.Guild {
		return false, "Este comando solo puede usarse en un servidor."
	}

	return true, ""
}
