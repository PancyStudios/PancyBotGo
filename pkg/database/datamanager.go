// Package database provides the DataManager for cached database operations.
package database

import (
	"container/list"
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DataManagerOptions contains configuration for a DataManager
type DataManagerOptions struct {
	MaxCacheSize int
}

// CacheManager provides shared caching across DataManagers
type CacheManager struct {
	cache     map[string]*list.Element
	cacheList *list.List
	mu        sync.RWMutex
}

// cacheEntry holds a cached value with its key
type cacheEntry struct {
	key   string
	value interface{}
}

// globalCacheManager is shared across all DataManager instances
var globalCacheManager = &CacheManager{
	cache:     make(map[string]*list.Element),
	cacheList: list.New(),
}

// global DataManagers for shared collections
var (
	GlobalWarnDM         *DataManager[models.WarnsDocument]
	GlobalUserPremiumDM  *DataManager[models.UserPremium]
	GlobalGuildPremiumDM *DataManager[models.GuildPremium]
	GlobalPremiumCodeDM  *DataManager[models.PremiumCode]
	GlobalBlacklistDM    *DataManager[models.Blacklist]
)

// InitGlobalDataManagers initializes shared DataManager instances
func InitGlobalDataManagers(db *Database) {
	GlobalWarnDM = NewDataManager[models.WarnsDocument]("warns", db)
	GlobalUserPremiumDM = NewDataManager[models.UserPremium]("premium", db)
	GlobalGuildPremiumDM = NewDataManager[models.GuildPremium]("premium_guilds", db)
	GlobalPremiumCodeDM = NewDataManager[models.PremiumCode]("premium_codes", db)
	GlobalBlacklistDM = NewDataManager[models.Blacklist]("blacklist", db)
}

// DataManager provides cached access to a MongoDB collection
type DataManager[T any] struct {
	collection *mongo.Collection
	dbInstance *Database
	options    DataManagerOptions
}

// DefaultDataManagerOptions returns default options for DataManager
func DefaultDataManagerOptions() DataManagerOptions {
	return DataManagerOptions{
		MaxCacheSize: 1000,
	}
}

// NewDataManager creates a new DataManager for a collection
func NewDataManager[T any](collectionName string, db *Database, opts ...DataManagerOptions) *DataManager[T] {
	dmOptions := DefaultDataManagerOptions()
	if len(opts) > 0 {
		dmOptions = opts[0]
	}

	return &DataManager[T]{
		collection: db.GetCollection(collectionName),
		dbInstance: db,
		options:    dmOptions,
	}
}

// generateCacheKey creates a unique, deterministic key from a query
// It sorts the keys to ensure consistent ordering regardless of map iteration order
func (dm *DataManager[T]) generateCacheKey(query bson.M) string {
	collName := ""
	if dm.collection != nil {
		collName = dm.collection.Name()
	}

	// Sort keys for deterministic serialization
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build a deterministic key string
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, query[k]))
	}

	return fmt.Sprintf("%s:{%s}", collName, strings.Join(parts, ","))
}

// Get retrieves a document from cache or database
func (dm *DataManager[T]) Get(query bson.M) (*T, error) {
	cacheKey := dm.generateCacheKey(query)

	// Check cache first
	globalCacheManager.mu.RLock()
	if elem, exists := globalCacheManager.cache[cacheKey]; exists {
		// Move to front (LRU)
		globalCacheManager.mu.RUnlock()
		globalCacheManager.mu.Lock()
		globalCacheManager.cacheList.MoveToFront(elem)
		entry := elem.Value.(*cacheEntry)
		globalCacheManager.mu.Unlock()
		return entry.value.(*T), nil
	}
	globalCacheManager.mu.RUnlock()

	// Not in cache, fetch from database
	if !dm.dbInstance.Connected() || dm.collection == nil {
		return nil, fmt.Errorf("database not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result T
	err := dm.collection.FindOne(ctx, query).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.Warn(fmt.Sprintf("Fallo al leer de la DB (%s), intentando desde caché...", dm.collection.Name()), "DataManager")
		return nil, err
	}

	// Add to cache
	globalCacheManager.mu.Lock()
	defer globalCacheManager.mu.Unlock()

	entry := &cacheEntry{key: cacheKey, value: &result}
	elem := globalCacheManager.cacheList.PushFront(entry)
	globalCacheManager.cache[cacheKey] = elem

	// Evict if over capacity
	if dm.options.MaxCacheSize > 0 && globalCacheManager.cacheList.Len() > dm.options.MaxCacheSize {
		oldest := globalCacheManager.cacheList.Back()
		if oldest != nil {
			oldEntry := oldest.Value.(*cacheEntry)
			delete(globalCacheManager.cache, oldEntry.key)
			globalCacheManager.cacheList.Remove(oldest)
		}
	}

	return &result, nil
}

// GetAll retrieves all documents matching a query from the database
func (dm *DataManager[T]) GetAll(query bson.M) ([]*T, error) {
	if !dm.dbInstance.Connected() || dm.collection == nil {
		return nil, fmt.Errorf("database not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := dm.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	var results []*T
	for cursor.Next(ctx) {
		var doc T
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		results = append(results, &doc)
	}

	return results, cursor.Err()
}

// Set updates or inserts a document in the database and cache
func (dm *DataManager[T]) Set(query bson.M, data interface{}) (*T, error) {
	cacheKey := dm.generateCacheKey(query)

	if !dm.dbInstance.Connected() || dm.collection == nil {
		// Queue for later
		logger.Warn(fmt.Sprintf("DB offline. Encolando escritura para '%s'", dm.collection.Name()), "DataManager")
		dm.dbInstance.AddToWriteQueue(QueuedOperation{
			CollectionName: dm.collection.Name(),
			Query:          query,
			Operation:      "set",
			Data:           data,
		})
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result T
	err := dm.collection.FindOneAndUpdate(ctx, query, bson.M{"$set": data}, opts).Decode(&result)
	if err != nil {
		logger.Error("Error en 'set' con DB conectada. Encolando por seguridad.", "DataManager")
		dm.dbInstance.AddToWriteQueue(QueuedOperation{
			CollectionName: dm.collection.Name(),
			Query:          query,
			Operation:      "set",
			Data:           data,
		})
		return nil, err
	}

	// Update cache
	globalCacheManager.mu.Lock()
	defer globalCacheManager.mu.Unlock()

	entry := &cacheEntry{key: cacheKey, value: &result}

	if elem, exists := globalCacheManager.cache[cacheKey]; exists {
		elem.Value = entry
		globalCacheManager.cacheList.MoveToFront(elem)
	} else {
		elem := globalCacheManager.cacheList.PushFront(entry)
		globalCacheManager.cache[cacheKey] = elem

		// Evict if over capacity
		if dm.options.MaxCacheSize > 0 && globalCacheManager.cacheList.Len() > dm.options.MaxCacheSize {
			oldest := globalCacheManager.cacheList.Back()
			if oldest != nil {
				oldEntry := oldest.Value.(*cacheEntry)
				delete(globalCacheManager.cache, oldEntry.key)
				globalCacheManager.cacheList.Remove(oldest)
			}
		}
	}

	return &result, nil
}

// Delete removes a document from the database and cache
func (dm *DataManager[T]) Delete(query bson.M) error {
	cacheKey := dm.generateCacheKey(query)

	// Remove from cache first
	globalCacheManager.mu.Lock()
	if elem, exists := globalCacheManager.cache[cacheKey]; exists {
		globalCacheManager.cacheList.Remove(elem)
		delete(globalCacheManager.cache, cacheKey)
	}
	globalCacheManager.mu.Unlock()

	if !dm.dbInstance.Connected() || dm.collection == nil {
		logger.Warn(fmt.Sprintf("DB offline. Encolando eliminación para '%s'", dm.collection.Name()), "DataManager")
		dm.dbInstance.AddToWriteQueue(QueuedOperation{
			CollectionName: dm.collection.Name(),
			Query:          query,
			Operation:      "delete",
		})
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := dm.collection.DeleteOne(ctx, query)
	if err != nil {
		logger.Error("Error en 'delete' con DB conectada. Encolando por seguridad.", "DataManager")
		dm.dbInstance.AddToWriteQueue(QueuedOperation{
			CollectionName: dm.collection.Name(),
			Query:          query,
			Operation:      "delete",
		})
		return err
	}

	return nil
}

// ClearCache clears the entire cache
func (dm *DataManager[T]) ClearCache() {
	globalCacheManager.mu.Lock()
	defer globalCacheManager.mu.Unlock()

	globalCacheManager.cache = make(map[string]*list.Element)
	globalCacheManager.cacheList = list.New()
}

// CacheSize returns the current cache size
func (dm *DataManager[T]) CacheSize() int {
	globalCacheManager.mu.RLock()
	defer globalCacheManager.mu.RUnlock()
	return globalCacheManager.cacheList.Len()
}

// PrimeCache logs that the cache is ready (caches are filled on demand)
func (dm *DataManager[T]) PrimeCache() {
	collName := ""
	if dm.collection != nil {
		collName = dm.collection.Name()
	}
	logger.System(fmt.Sprintf("Caché para '%s' preparada (tamaño máx: %d). Se llenará bajo demanda.", collName, dm.options.MaxCacheSize), "DataManager")
}
