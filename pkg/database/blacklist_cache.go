// Package database provides the BlacklistCache for in-memory blacklist operations.
package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

// BlacklistCache provides an in-memory cache for blacklist entries
type BlacklistCache struct {
	entries     map[string]*models.Blacklist
	mu          sync.RWMutex
	stopRefresh chan struct{}
	refreshing  bool
}

var (
	blacklistCache *BlacklistCache
	cacheOnce      sync.Once
)

// GetBlacklistCache returns the global blacklist cache instance
func GetBlacklistCache() *BlacklistCache {
	cacheOnce.Do(func() {
		blacklistCache = &BlacklistCache{
			entries:     make(map[string]*models.Blacklist),
			stopRefresh: make(chan struct{}),
		}
	})
	return blacklistCache
}

// InitBlacklistCache initializes the blacklist cache by loading all entries from the database
func InitBlacklistCache() error {
	cache := GetBlacklistCache()
	return cache.Refresh()
}

// StartBlacklistCacheRefresh starts the automatic cache refresh every 5 minutes
func StartBlacklistCacheRefresh() {
	cache := GetBlacklistCache()
	cache.StartAutoRefresh(5 * time.Minute)
}

// StopBlacklistCacheRefresh stops the automatic cache refresh
func StopBlacklistCacheRefresh() {
	cache := GetBlacklistCache()
	cache.StopAutoRefresh()
}

// RefreshBlacklistCache manually refreshes the blacklist cache
func RefreshBlacklistCache() error {
	cache := GetBlacklistCache()
	return cache.Refresh()
}

// Refresh reloads all blacklist entries from the database
func (c *BlacklistCache) Refresh() error {
	if GlobalBlacklistDM == nil {
		logger.Warn("BlacklistCache: DataManager not initialized", "BlacklistCache")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := GlobalBlacklistDM.collection
	if collection == nil {
		logger.Warn("BlacklistCache: Collection not available", "BlacklistCache")
		return nil
	}

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error("BlacklistCache: Error fetching blacklist entries: "+err.Error(), "BlacklistCache")
		return err
	}
	defer func() { _ = cursor.Close(ctx) }()

	newEntries := make(map[string]*models.Blacklist)
	for cursor.Next(ctx) {
		var entry models.Blacklist
		if err := cursor.Decode(&entry); err != nil {
			logger.Warn("BlacklistCache: Error decoding entry: "+err.Error(), "BlacklistCache")
			continue
		}
		newEntries[entry.ID] = &entry
	}

	if err := cursor.Err(); err != nil {
		logger.Error("BlacklistCache: Cursor error: "+err.Error(), "BlacklistCache")
		return err
	}

	c.mu.Lock()
	c.entries = newEntries
	c.mu.Unlock()

	logger.Info(fmt.Sprintf("BlacklistCache: Cache refreshed with %d entries", len(newEntries)), "BlacklistCache")
	return nil
}

// StartAutoRefresh starts automatic cache refresh at the specified interval
// If already refreshing, it will stop the current refresher and start a new one
func (c *BlacklistCache) StartAutoRefresh(interval time.Duration) {
	c.mu.Lock()
	// Stop existing refresher if running
	if c.refreshing {
		close(c.stopRefresh)
		c.refreshing = false
	}
	c.refreshing = true
	c.stopRefresh = make(chan struct{})
	stopChan := c.stopRefresh
	c.mu.Unlock()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		logger.Info("BlacklistCache: Auto-refresh started (interval: "+interval.String()+")", "BlacklistCache")

		for {
			select {
			case <-ticker.C:
				if err := c.Refresh(); err != nil {
					logger.Error("BlacklistCache: Auto-refresh failed: "+err.Error(), "BlacklistCache")
				}
			case <-stopChan:
				logger.Info("BlacklistCache: Auto-refresh stopped", "BlacklistCache")
				return
			}
		}
	}()
}

// StopAutoRefresh stops the automatic cache refresh
func (c *BlacklistCache) StopAutoRefresh() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.refreshing {
		close(c.stopRefresh)
		c.refreshing = false
	}
}

// Get retrieves a blacklist entry from the cache
func (c *BlacklistCache) Get(id string) (*models.Blacklist, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[id]
	return entry, exists
}

// IsBlacklisted checks if an ID is in the blacklist cache
func (c *BlacklistCache) IsBlacklisted(id string) bool {
	_, exists := c.Get(id)
	return exists
}

// IsUserBlacklisted checks if a user ID is in the blacklist and is of type user
func (c *BlacklistCache) IsUserBlacklisted(userID string) (bool, *models.Blacklist) {
	entry, exists := c.Get(userID)
	if !exists || entry.Type != models.BlacklistTypeUser {
		return false, nil
	}
	return true, entry
}

// IsGuildBlacklisted checks if a guild ID is in the blacklist and is of type guild
func (c *BlacklistCache) IsGuildBlacklisted(guildID string) (bool, *models.Blacklist) {
	entry, exists := c.Get(guildID)
	if !exists || entry.Type != models.BlacklistTypeGuild {
		return false, nil
	}
	return true, entry
}

// GetAll returns all blacklist entries from the cache
func (c *BlacklistCache) GetAll() []*models.Blacklist {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]*models.Blacklist, 0, len(c.entries))
	for _, entry := range c.entries {
		result = append(result, entry)
	}
	return result
}

// GetByType returns blacklist entries of a specific type from the cache
func (c *BlacklistCache) GetByType(blacklistType models.BlacklistType) []*models.Blacklist {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []*models.Blacklist
	for _, entry := range c.entries {
		if entry.Type == blacklistType {
			result = append(result, entry)
		}
	}
	return result
}

// Add adds an entry to the cache (called after DB add)
func (c *BlacklistCache) Add(entry *models.Blacklist) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[entry.ID] = entry
}

// Remove removes an entry from the cache (called after DB remove)
func (c *BlacklistCache) Remove(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, id)
}

// Size returns the number of entries in the cache
func (c *BlacklistCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
