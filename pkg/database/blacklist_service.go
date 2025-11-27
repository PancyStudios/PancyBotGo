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

	// Verificar si ya está en la blacklist (usando cache)
	cache := GetBlacklistCache()
	if _, exists := cache.Get(id); exists {
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

	// Update cache after successful DB write
	cache.Add(result)

	return result, nil
}

// RemoveFromBlacklist elimina un usuario o guild de la blacklist
func RemoveFromBlacklist(id string) error {
	if GlobalBlacklistDM == nil {
		return errors.New("blacklist manager not initialized")
	}

	// Verificar que existe (usando cache)
	cache := GetBlacklistCache()
	if _, exists := cache.Get(id); !exists {
		return ErrBlacklistNotFound
	}

	err := GlobalBlacklistDM.Delete(bson.M{"_id": id})
	if err != nil {
		return err
	}

	// Update cache after successful DB delete
	cache.Remove(id)

	return nil
}

// GetBlacklistEntry obtiene una entrada de la blacklist (from cache)
func GetBlacklistEntry(id string) (*models.Blacklist, error) {
	cache := GetBlacklistCache()
	entry, exists := cache.Get(id)
	if !exists {
		return nil, ErrBlacklistNotFound
	}
	return entry, nil
}

// IsBlacklisted verifica si un ID está en la blacklist (from cache)
func IsBlacklisted(id string) bool {
	cache := GetBlacklistCache()
	return cache.IsBlacklisted(id)
}

// IsUserBlacklisted verifica si un usuario está en la blacklist (from cache)
func IsUserBlacklisted(userID string) (bool, *models.Blacklist) {
	cache := GetBlacklistCache()
	return cache.IsUserBlacklisted(userID)
}

// IsGuildBlacklisted verifica si un guild está en la blacklist (from cache)
func IsGuildBlacklisted(guildID string) (bool, *models.Blacklist) {
	cache := GetBlacklistCache()
	return cache.IsGuildBlacklisted(guildID)
}

// GetAllBlacklist obtiene todas las entradas de la blacklist (from cache)
func GetAllBlacklist() ([]*models.Blacklist, error) {
	cache := GetBlacklistCache()
	return cache.GetAll(), nil
}

// GetBlacklistByType obtiene entradas por tipo (from cache)
func GetBlacklistByType(blacklistType models.BlacklistType) ([]*models.Blacklist, error) {
	cache := GetBlacklistCache()
	return cache.GetByType(blacklistType), nil
}
