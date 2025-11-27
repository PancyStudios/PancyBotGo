package database

import (
	"errors"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrBlacklistNotFound  = errors.New("entrada no encontrada en la blacklist")
	ErrAlreadyBlacklisted = errors.New("ya está en la blacklist")
)

// AddToBlacklist añade un usuario o guild a la blacklist
func AddToBlacklist(id string, blacklistType models.BlacklistType, reason string, addedBy string) (*models.Blacklist, error) {
	if GlobalBlacklistDM == nil {
		return nil, errors.New("blacklist manager not initialized")
	}

	// Verificar si ya está en la blacklist
	existing, err := GetBlacklistEntry(id)
	if err == nil && existing != nil {
		return nil, ErrAlreadyBlacklisted
	}

	entry := models.Blacklist{
		ID:      id,
		Type:    blacklistType,
		Reason:  reason,
		AddedBy: addedBy,
		AddedAt: time.Now(),
	}

	result, err := GlobalBlacklistDM.Set(bson.M{"_id": id}, entry)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RemoveFromBlacklist elimina un usuario o guild de la blacklist
func RemoveFromBlacklist(id string) error {
	if GlobalBlacklistDM == nil {
		return errors.New("blacklist manager not initialized")
	}

	// Verificar que existe
	_, err := GetBlacklistEntry(id)
	if err != nil {
		return ErrBlacklistNotFound
	}

	return GlobalBlacklistDM.Delete(bson.M{"_id": id})
}

// GetBlacklistEntry obtiene una entrada de la blacklist
func GetBlacklistEntry(id string) (*models.Blacklist, error) {
	if GlobalBlacklistDM == nil {
		return nil, errors.New("blacklist manager not initialized")
	}

	entry, err := GlobalBlacklistDM.Get(bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, ErrBlacklistNotFound
	}

	return entry, nil
}

// IsBlacklisted verifica si un ID está en la blacklist
func IsBlacklisted(id string) bool {
	entry, err := GetBlacklistEntry(id)
	return err == nil && entry != nil
}

// IsUserBlacklisted verifica si un usuario está en la blacklist
func IsUserBlacklisted(userID string) (bool, *models.Blacklist) {
	entry, err := GetBlacklistEntry(userID)
	if err != nil || entry == nil || entry.Type != models.BlacklistTypeUser {
		return false, nil
	}
	return true, entry
}

// IsGuildBlacklisted verifica si un guild está en la blacklist
func IsGuildBlacklisted(guildID string) (bool, *models.Blacklist) {
	entry, err := GetBlacklistEntry(guildID)
	if err != nil || entry == nil || entry.Type != models.BlacklistTypeGuild {
		return false, nil
	}
	return true, entry
}

// GetAllBlacklist obtiene todas las entradas de la blacklist
func GetAllBlacklist() ([]*models.Blacklist, error) {
	if GlobalBlacklistDM == nil {
		return nil, errors.New("blacklist manager not initialized")
	}

	return GlobalBlacklistDM.GetAll(bson.M{})
}

// GetBlacklistByType obtiene entradas por tipo
func GetBlacklistByType(blacklistType models.BlacklistType) ([]*models.Blacklist, error) {
	if GlobalBlacklistDM == nil {
		return nil, errors.New("blacklist manager not initialized")
	}

	return GlobalBlacklistDM.GetAll(bson.M{"type": blacklistType})
}
