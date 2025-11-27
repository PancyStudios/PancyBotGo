package database

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/logger"
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

// BlacklistCache provides in-memory caching for blacklist entries
type BlacklistCache struct {
	entries  map[string]*models.BlacklistEntry
	mu       sync.RWMutex
	ticker   *time.Ticker
	done     chan bool
	stopOnce sync.Once
}

var blacklistCache = &BlacklistCache{
	entries: make(map[string]*models.BlacklistEntry),
	done:    make(chan bool),
}

// InitBlacklistCache initializes and loads the blacklist cache from the database
// Should be called at bot startup after InitGlobalDataManagers
func InitBlacklistCache() error {
	return RefreshBlacklistCache()
}

// StartBlacklistCacheRefresh starts a goroutine that refreshes the cache every 5 minutes
func StartBlacklistCacheRefresh() {
	blacklistCache.ticker = time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <-blacklistCache.done:
				return
			case <-blacklistCache.ticker.C:
				if err := RefreshBlacklistCache(); err != nil {
					logger.Error("Error refrescando caché de blacklist: "+err.Error(), "BlacklistCache")
				} else {
					logger.Debug("Caché de blacklist refrescada automáticamente", "BlacklistCache")
				}
			}
		}
	}()

	logger.System("Sistema de caché de blacklist iniciado (refresco cada 5 minutos)", "BlacklistCache")
}

// StopBlacklistCacheRefresh stops the cache refresh goroutine
func StopBlacklistCacheRefresh() {
	blacklistCache.stopOnce.Do(func() {
		if blacklistCache.ticker != nil {
			blacklistCache.ticker.Stop()
		}
		close(blacklistCache.done)
	})
}

// RefreshBlacklistCache reloads all blacklist entries from the database into cache
func RefreshBlacklistCache() error {
	dm, err := getBlacklistManager()
	if err != nil {
		return err
	}

	entries, err := dm.GetAll(bson.M{})
	if err != nil {
		return err
	}

	blacklistCache.mu.Lock()
	defer blacklistCache.mu.Unlock()

	// Clear existing cache
	blacklistCache.entries = make(map[string]*models.BlacklistEntry)

	// Load all entries into cache
	for _, entry := range entries {
		blacklistCache.entries[entry.ID] = entry
	}

	logger.Info(fmt.Sprintf("Caché de blacklist cargada: %d entradas", len(blacklistCache.entries)), "BlacklistCache")
	return nil
}

// addToCache adds an entry to the in-memory cache
func addToCache(entry *models.BlacklistEntry) {
	blacklistCache.mu.Lock()
	defer blacklistCache.mu.Unlock()
	blacklistCache.entries[entry.ID] = entry
}

// removeFromCache removes an entry from the in-memory cache
func removeFromCache(id string) {
	blacklistCache.mu.Lock()
	defer blacklistCache.mu.Unlock()
	delete(blacklistCache.entries, id)
}

// getFromCache retrieves an entry from the in-memory cache
func getFromCache(id string) (*models.BlacklistEntry, bool) {
	blacklistCache.mu.RLock()
	defer blacklistCache.mu.RUnlock()
	entry, exists := blacklistCache.entries[id]
	return entry, exists
}

func getBlacklistManager() (*DataManager[models.BlacklistEntry], error) {
	if GlobalBlacklistDM == nil {
		return nil, ErrBlacklistManagerNotInitialized
	}
	return GlobalBlacklistDM, nil
}

// AddToBlacklist adds a user or guild to the blacklist
func AddToBlacklist(id string, blacklistType models.BlacklistType, reason string, createdBy string) (*models.BlacklistEntry, error) {
	// Check cache first for duplicates (fast check)
	if _, exists := getFromCache(id); exists {
		return nil, ErrBlacklistEntryExists
	}

	dm, err := getBlacklistManager()
	if err != nil {
		return nil, err
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

	// Update cache immediately
	addToCache(result)

	return result, nil
}

// RemoveFromBlacklist removes a user or guild from the blacklist
func RemoveFromBlacklist(id string) error {
	// Check cache first
	if _, exists := getFromCache(id); !exists {
		return ErrBlacklistEntryNotFound
	}

	dm, err := getBlacklistManager()
	if err != nil {
		return err
	}

	err = dm.Delete(bson.M{"_id": id})
	if err != nil {
		return err
	}

	// Update cache immediately
	removeFromCache(id)

	return nil
}

// GetBlacklistEntry gets a specific blacklist entry from cache
func GetBlacklistEntry(id string) (*models.BlacklistEntry, error) {
	entry, exists := getFromCache(id)
	if !exists {
		return nil, ErrBlacklistEntryNotFound
	}
	return entry, nil
}

// IsBlacklisted checks if a user or guild is blacklisted (from cache - no DB delay)
func IsBlacklisted(id string) (bool, *models.BlacklistEntry, error) {
	entry, exists := getFromCache(id)
	return exists, entry, nil
}

// IsUserBlacklisted checks if a user is blacklisted (from cache - no DB delay)
func IsUserBlacklisted(userID string) (bool, *models.BlacklistEntry, error) {
	return IsBlacklisted(userID)
}

// IsGuildBlacklisted checks if a guild is blacklisted (from cache - no DB delay)
func IsGuildBlacklisted(guildID string) (bool, *models.BlacklistEntry, error) {
	return IsBlacklisted(guildID)
}

// GetAllBlacklistEntries gets all blacklist entries from cache
func GetAllBlacklistEntries() ([]*models.BlacklistEntry, error) {
	blacklistCache.mu.RLock()
	defer blacklistCache.mu.RUnlock()

	entries := make([]*models.BlacklistEntry, 0, len(blacklistCache.entries))
	for _, entry := range blacklistCache.entries {
		entries = append(entries, entry)
	}
	return entries, nil
}

// GetBlacklistEntriesByType gets all blacklist entries of a specific type from cache
func GetBlacklistEntriesByType(blacklistType models.BlacklistType) ([]*models.BlacklistEntry, error) {
	blacklistCache.mu.RLock()
	defer blacklistCache.mu.RUnlock()

	var entries []*models.BlacklistEntry
	for _, entry := range blacklistCache.entries {
		if entry.Type == blacklistType {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}
