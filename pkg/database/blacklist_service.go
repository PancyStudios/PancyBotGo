package database

import (
	"errors"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

// Alias de tipos para facilitar el acceso
type BlacklistEntry = models.BlacklistEntry

var (
	ErrBlacklistManagerNotInitialized = errors.New("blacklist data manager not initialized")
	ErrBlacklistEntryNotFound         = errors.New("entrada de blacklist no encontrada")
	ErrBlacklistEntryExists           = errors.New("la entrada ya existe en la blacklist")
)

func getBlacklistManager() (*DataManager[models.BlacklistEntry], error) {
	if GlobalBlacklistDM == nil {
		return nil, ErrBlacklistManagerNotInitialized
	}
	return GlobalBlacklistDM, nil
}

// AddToBlacklist adds a user or guild to the blacklist
func AddToBlacklist(id string, blacklistType models.BlacklistType, reason string, createdBy string) (*models.BlacklistEntry, error) {
	dm, err := getBlacklistManager()
	if err != nil {
		return nil, err
	}

	// Check if entry already exists
	existing, err := dm.Get(bson.M{"_id": id})
	if err == nil && existing != nil {
		return nil, ErrBlacklistEntryExists
	}

	entry := models.BlacklistEntry{
		ID:        id,
		Type:      blacklistType,
		Reason:    reason,
		CreatedAt: time.Now(),
		CreatedBy: createdBy,
	}

	result, err := dm.Set(bson.M{"_id": id}, entry)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RemoveFromBlacklist removes a user or guild from the blacklist
func RemoveFromBlacklist(id string) error {
	dm, err := getBlacklistManager()
	if err != nil {
		return err
	}

	// Check if entry exists
	existing, err := dm.Get(bson.M{"_id": id})
	if err != nil || existing == nil {
		return ErrBlacklistEntryNotFound
	}

	return dm.Delete(bson.M{"_id": id})
}

// GetBlacklistEntry gets a specific blacklist entry
func GetBlacklistEntry(id string) (*models.BlacklistEntry, error) {
	dm, err := getBlacklistManager()
	if err != nil {
		return nil, err
	}

	entry, err := dm.Get(bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, ErrBlacklistEntryNotFound
	}

	return entry, nil
}

// IsBlacklisted checks if a user or guild is blacklisted
func IsBlacklisted(id string) (bool, *models.BlacklistEntry, error) {
	dm, err := getBlacklistManager()
	if err != nil {
		return false, nil, err
	}

	entry, err := dm.Get(bson.M{"_id": id})
	if err != nil {
		return false, nil, err
	}

	return entry != nil, entry, nil
}

// IsUserBlacklisted checks if a user is blacklisted
func IsUserBlacklisted(userID string) (bool, *models.BlacklistEntry, error) {
	return IsBlacklisted(userID)
}

// IsGuildBlacklisted checks if a guild is blacklisted
func IsGuildBlacklisted(guildID string) (bool, *models.BlacklistEntry, error) {
	return IsBlacklisted(guildID)
}

// GetAllBlacklistEntries gets all blacklist entries
func GetAllBlacklistEntries() ([]*models.BlacklistEntry, error) {
	dm, err := getBlacklistManager()
	if err != nil {
		return nil, err
	}

	return dm.GetAll(bson.M{})
}

// GetBlacklistEntriesByType gets all blacklist entries of a specific type
func GetBlacklistEntriesByType(blacklistType models.BlacklistType) ([]*models.BlacklistEntry, error) {
	dm, err := getBlacklistManager()
	if err != nil {
		return nil, err
	}

	return dm.GetAll(bson.M{"type": blacklistType})
}
