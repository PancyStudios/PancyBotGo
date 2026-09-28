// Package database provides the DataManager for cached database operations.
package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	cachepkg "github.com/PancyStudios/PancyBotGo/pkg/cache"
	"github.com/PancyStudios/PancyBotGo/pkg/logger"
	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrDatabaseUnavailable = errors.New("database unavailable")

// DataManagerOptions contains configuration for a DataManager
type DataManagerOptions struct {
	MaxCacheSize int
}

const cacheOperationTimeout = 900 * time.Millisecond

var (
	cacheMu       sync.RWMutex
	cacheStore    cachepkg.Store = cachepkg.NewMemoryStore(1000)
	cacheEntryTTL                = 5 * time.Minute
)

// ConfigureCache replaces the process-local default with another cache backend.
// MongoDB remains the source of truth, so cache failures never block database access.
func ConfigureCache(store cachepkg.Store, ttl time.Duration) {
	if store == nil {
		return
	}
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	cacheMu.Lock()
	cacheStore = store
	cacheEntryTTL = ttl
	cacheMu.Unlock()
}

func currentCache() (cachepkg.Store, time.Duration) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	return cacheStore, cacheEntryTTL
}

// global DataManagers for shared collections
var (
	GlobalWarnDM         *DataManager[models.WarnsDocument]
	GlobalUserPremiumDM  *DataManager[models.UserPremium]
	GlobalGuildPremiumDM *DataManager[models.GuildPremium]
	GlobalPremiumCodeDM  *DataManager[models.PremiumCode]
	GlobalBlacklistDM    *DataManager[models.Blacklist]
	GlobalEconomyDM      *DataManager[models.GlobalEconomyProfile]
	LocalEconomyDM       *DataManager[models.LocalEconomyProfile]
	LocalLevelsDM        *DataManager[models.UserLevelProfile]
	ItemDM               *DataManager[models.Item]
	GlobalGuildDM        *DataManager[models.GuildDocument]
	GlobalMusicDM        *DataManager[models.MusicSettings]
)

// InitGlobalDataManagers initializes shared DataManager instances
func InitGlobalDataManagers(db *Database) {
	GlobalWarnDM = NewDataManager[models.WarnsDocument]("warns", db)
	GlobalUserPremiumDM = NewDataManager[models.UserPremium]("premium", db)
	GlobalGuildPremiumDM = NewDataManager[models.GuildPremium]("premium_guilds", db)
	GlobalPremiumCodeDM = NewDataManager[models.PremiumCode]("premium_codes", db)
	GlobalBlacklistDM = NewDataManager[models.Blacklist]("blacklist", db)
	GlobalGuildDM = NewDataManager[models.GuildDocument]("guilds", db)
	GlobalMusicDM = NewDataManager[models.MusicSettings]("music", db)
	GlobalEconomyDM = NewDataManager[models.GlobalEconomyProfile]("economy_global", db)
	LocalEconomyDM = NewDataManager[models.LocalEconomyProfile]("economy_local", db)
	LocalLevelsDM = NewDataManager[models.UserLevelProfile]("levels", db)
	ItemDM = NewDataManager[models.Item]("economy_items", db)
}

// DataManager provides cached access to a MongoDB collection
type DataManager[T any] struct {
	collectionName string
	dbInstance     *Database
	options        DataManagerOptions
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
		collectionName: collectionName,
		dbInstance:     db,
		options:        dmOptions,
	}
}

func (dm *DataManager[T]) collection() *mongo.Collection {
	if dm.dbInstance == nil || !dm.dbInstance.Connected() {
		return nil
	}
	return dm.dbInstance.GetCollection(dm.collectionName)
}

func cloneValue[T any](value *T) (*T, error) {
	if value == nil {
		return nil, nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var cloned T
	if err := json.Unmarshal(data, &cloned); err != nil {
		return nil, err
	}
	return &cloned, nil
}

// generateCacheKey creates a unique, deterministic key from a query
// It sorts the keys to ensure consistent ordering regardless of map iteration order
func (dm *DataManager[T]) generateCacheKey(query bson.M) string {
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

	return fmt.Sprintf("%s:{%s}", dm.collectionName, strings.Join(parts, ","))
}

func (dm *DataManager[T]) getCached(cacheKey string) (*T, bool) {
	store, _ := currentCache()
	ctx, cancel := context.WithTimeout(context.Background(), cacheOperationTimeout)
	defer cancel()

	data, err := store.Get(ctx, cacheKey)
	if errors.Is(err, cachepkg.ErrMiss) {
		return nil, false
	}
	if err != nil {
		logger.Debug(fmt.Sprintf("Caché '%s' no disponible para '%s': %v", store.Name(), dm.collectionName, err), "DataManager")
		return nil, false
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		logger.Warn(fmt.Sprintf("Entrada de caché inválida para '%s': %v", dm.collectionName, err), "DataManager")
		_ = store.Delete(context.Background(), cacheKey)
		return nil, false
	}
	return &value, true
}

func (dm *DataManager[T]) setCached(cacheKey string, value *T) {
	if value == nil {
		return
	}
	data, err := json.Marshal(value)
	if err != nil {
		logger.Debug(fmt.Sprintf("No se pudo serializar caché para '%s': %v", dm.collectionName, err), "DataManager")
		return
	}
	store, ttl := currentCache()
	ctx, cancel := context.WithTimeout(context.Background(), cacheOperationTimeout)
	defer cancel()
	if err := store.Set(ctx, cacheKey, data, ttl); err != nil {
		logger.Debug(fmt.Sprintf("No se pudo actualizar caché '%s' para '%s': %v", store.Name(), dm.collectionName, err), "DataManager")
	}
}

func (dm *DataManager[T]) deleteCached(cacheKey string) {
	store, _ := currentCache()
	ctx, cancel := context.WithTimeout(context.Background(), cacheOperationTimeout)
	defer cancel()
	if err := store.Delete(ctx, cacheKey); err != nil {
		logger.Debug(fmt.Sprintf("No se pudo invalidar caché '%s' para '%s': %v", store.Name(), dm.collectionName, err), "DataManager")
	}
}

// Get retrieves a document from cache or database
func (dm *DataManager[T]) Get(query bson.M) (*T, error) {
	cacheKey := dm.generateCacheKey(query)

	if value, ok := dm.getCached(cacheKey); ok {
		return value, nil
	}

	// Not in cache, fetch from database
	collection := dm.collection()
	if collection == nil {
		return nil, ErrDatabaseUnavailable
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	var result T
	err := collection.FindOne(ctx, query).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.Warn(fmt.Sprintf("Error leyendo de DB (%s): %v", dm.collectionName, err), "DataManager")
		return nil, err
	}

	dm.setCached(cacheKey, &result)
	return &result, nil
}

// GetAll retrieves all documents matching a query from the database
func (dm *DataManager[T]) GetAll(query bson.M) ([]*T, error) {
	collection := dm.collection()
	if collection == nil {
		return nil, ErrDatabaseUnavailable
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	var results []*T
	for cursor.Next(ctx) {
		var doc T
		if err := cursor.Decode(&doc); err != nil {
			logger.Warn(fmt.Sprintf("Error decoding document in GetAll for '%s': %v", dm.collectionName, err), "DataManager")
			continue
		}
		results = append(results, &doc)
	}

	return results, cursor.Err()
}

// Set updates or inserts a document in the database and cache
func (dm *DataManager[T]) Set(query bson.M, data interface{}) (*T, error) {
	cacheKey := dm.generateCacheKey(query)

	// Update cache immediately (Write-Through)
	var cacheValue *T
	if v, ok := data.(*T); ok {
		cacheValue = v
	} else if v, ok := data.(T); ok {
		cacheValue = &v
	} else {
		bytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		var temp T
		if err := json.Unmarshal(bytes, &temp); err != nil {
			return nil, err
		}
		cacheValue = &temp
	}
	cacheValue, err := cloneValue(cacheValue)
	if err != nil {
		return nil, err
	}

	dm.setCached(cacheKey, cacheValue)

	collection := dm.collection()
	if collection == nil {
		logger.Warn(fmt.Sprintf("DB offline. Encolando escritura en '%s' y usando caché.", dm.collectionName), "DataManager")
		dm.dbInstance.AddToWriteQueue(QueuedOperation{
			CollectionName: dm.collectionName,
			Query:          query,
			Operation:      "set",
			Data:           data,
		})
		return cacheValue, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result T
	err = collection.FindOneAndUpdate(ctx, query, bson.M{"$set": data}, opts).Decode(&result)
	if err != nil {
		logger.Debug(fmt.Sprintf("DB offline o timeout. Usando cache local para '%s'", dm.collectionName), "DataManager")
		dm.dbInstance.AddToWriteQueue(QueuedOperation{
			CollectionName: dm.collectionName,
			Query:          query,
			Operation:      "set",
			Data:           data,
		})
		return cacheValue, nil
	}

	dm.setCached(cacheKey, &result)
	return &result, nil
}

// Delete removes a document from the database and cache
func (dm *DataManager[T]) Delete(query bson.M) error {
	cacheKey := dm.generateCacheKey(query)

	dm.deleteCached(cacheKey)

	collection := dm.collection()
	if collection == nil {
		logger.Warn(fmt.Sprintf("DB offline. Encolando eliminación para '%s'", dm.collectionName), "DataManager")
		dm.dbInstance.AddToWriteQueue(QueuedOperation{
			CollectionName: dm.collectionName,
			Query:          query,
			Operation:      "delete",
		})
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	_, err := collection.DeleteOne(ctx, query)
	if err != nil {
		logger.Debug("Eliminación añadida a la cola offline", "DataManager")
		dm.dbInstance.AddToWriteQueue(QueuedOperation{
			CollectionName: dm.collectionName,
			Query:          query,
			Operation:      "delete",
		})
		return err
	}

	return nil
}

// ClearCache clears the entire cache
func (dm *DataManager[T]) ClearCache() {
	store, _ := currentCache()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := store.Clear(ctx); err != nil {
		logger.Warn(fmt.Sprintf("No se pudo limpiar la caché '%s': %v", store.Name(), err), "DataManager")
	}
}

// CacheSize returns the current cache size
func (dm *DataManager[T]) CacheSize() int {
	store, _ := currentCache()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	size, err := store.Size(ctx)
	if err != nil {
		return 0
	}
	return size
}

// PrimeCache logs that the cache is ready (caches are filled on demand)
func (dm *DataManager[T]) PrimeCache() {
	store, ttl := currentCache()
	logger.System(fmt.Sprintf("Caché '%s' para '%s' preparada (TTL: %s). Se llenará bajo demanda.", store.Name(), dm.collectionName, ttl), "DataManager")
}
