package database

import (
	"sync"
)

const DEFAULT_COLLECTION = "default"

type CollectionManager struct {
	collections map[string]*Database
	mu          sync.RWMutex
}

func NewCollectionManager() *CollectionManager {
	return &CollectionManager{
		collections: make(map[string]*Database),
	}
}

func (cm *CollectionManager) AddCollection(name string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.collections[name] = NewDatabase()
}

func (cm *CollectionManager) GetCollection(name string) (*Database, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	db, exists := cm.collections[name]
	return db, exists
}

func (cm *CollectionManager) GetDefaultCollection() *Database {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	db := cm.collections[DEFAULT_COLLECTION]
	return db
}

func (cm *CollectionManager) GetCollectionNames() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var collectionNames []string
	for key := range cm.collections {
		collectionNames = append(collectionNames, key)
	}
	return collectionNames
}

func (cm *CollectionManager) RemoveCollection(name string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.collections, name)
}
