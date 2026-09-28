package database

import (
	"math"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

type cacheTestDocument struct {
	Values map[string]int `json:"values"`
}

func TestCheckedAddRejectsOverflow(t *testing.T) {
	if _, ok := checkedAdd(math.MaxInt64, 1); ok {
		t.Fatal("positive overflow was accepted")
	}
	if _, ok := checkedAdd(math.MinInt64, -1); ok {
		t.Fatal("negative overflow was accepted")
	}
	if value, ok := checkedAdd(10, -3); !ok || value != 7 {
		t.Fatalf("valid addition failed: value=%d ok=%v", value, ok)
	}
}

func TestCacheKeyIncludesCollectionWhenDatabaseIsOffline(t *testing.T) {
	db := NewDatabase()
	users := NewDataManager[cacheTestDocument]("users", db)
	guilds := NewDataManager[cacheTestDocument]("guilds", db)
	query := bson.M{"_id": "same-id"}

	if users.generateCacheKey(query) == guilds.generateCacheKey(query) {
		t.Fatal("cache keys from different collections must not collide")
	}
}

func TestCloneValueDoesNotShareMutableMaps(t *testing.T) {
	original := &cacheTestDocument{Values: map[string]int{"stars": 10}}
	cloned, err := cloneValue(original)
	if err != nil {
		t.Fatalf("cloneValue() returned error: %v", err)
	}

	cloned.Values["stars"] = 99
	if original.Values["stars"] != 10 {
		t.Fatal("cloneValue() returned a document sharing mutable map state")
	}
}
